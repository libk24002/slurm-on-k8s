/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"log"
	"strings"
	"time"

	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	slurmv1 "github.com/AaronYang0628/slurm-on-k8s/api/v1"

	utils "github.com/AaronYang0628/slurm-on-k8s/internal/utils"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// SlurmDeploymentReconciler reconciles a SlurmDeployment object
type SlurmDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=slurm.ay.dev,resources=slurmdeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=slurm.ay.dev,resources=slurmdeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=slurm.ay.dev,resources=slurmdeployments/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;create;update;patch;delete;watch
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;create;update;patch;delete;watch
// +kubebuilder:rbac:groups=apps,resources=deployments;statefulsets,verbs=get;list;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings,verbs=get;list;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=get;list;create;update;patch;delete
// +kubebuilder:rbac:groups=policy,resources=poddisruptionbudgets,verbs=get;list;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SlurmDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.

// SlurmDeploymentFinalizer is the name of the finalizer added to SlurmDeployment resources
const SlurmDeploymentFinalizer = "slurm.ay.dev/finalizer"

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.0/pkg/reconcile
func (r *SlurmDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// get CR SlurmDeployment instance
	release := &slurmv1.SlurmDeployment{}

	if findReleaseErr := r.Get(ctx, req.NamespacedName, release); findReleaseErr != nil {
		return ctrl.Result{}, client.IgnoreNotFound(findReleaseErr)
	}
	log.Printf("Find SlurmDeployment %s, Going to run Job: %v", release.Name, release.Spec.Job)

	// Create namespace if not exist
	if _, createNamespaceErr := r.CreateNamespaceIfNotExist(ctx, release.Spec.Chart.Namespace); createNamespaceErr != nil {
		log.Printf("Failed to create namespace [%s]: %v", release.Spec.Chart.Namespace, createNamespaceErr)
		return ctrl.Result{}, createNamespaceErr
	}
	log.Printf("Find Namespace [%s]", release.Spec.Chart.Namespace)

	// Initialize Helm settings and configuration
	helmSettings := cli.New()
	actionConfig := new(action.Configuration)
	if helmClientErr := actionConfig.Init(helmSettings.RESTClientGetter(), release.Spec.Chart.Namespace, "secret", log.Printf); helmClientErr != nil {
		log.Fatalf("Failed to initialize Helm configuration: %v", helmClientErr)
		return ctrl.Result{}, helmClientErr
	}
	log.Printf("Helm configuration initialized in namespace [%s]", release.Spec.Chart.Namespace)

	// Check if the SlurmDeployment is being deleted
	if !release.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is being deleted
		if utils.CheckIfExistInArray(release.ObjectMeta.Finalizers, SlurmDeploymentFinalizer) {
			// Our finalizer is present, so we need to clean up resources
			log.Printf("Deleting SlurmDeployment %s in namespace [%s]", release.Name, release.Spec.Chart.Namespace)

			// Uninstall the Helm release
			uninstallClient := action.NewUninstall(actionConfig)
			// uninstallClient.DisableHooks = true
			uninstallClient.Timeout = 60 * time.Second
			uninstallClient.Wait = false

			if _, uninstallErr := uninstallClient.Run(release.Name); uninstallErr != nil {
				if strings.Contains(uninstallErr.Error(), "release: not found") {
					log.Printf("SlurmDeployment %s not found, skipping uninstall", release.Name)
				} else if strings.Contains(uninstallErr.Error(), "timed out") || strings.Contains(uninstallErr.Error(), "BackoffLimitExceeded") {
					log.Printf("SlurmDeployment %s uninstall timed out or job failed, continuing with CR deletion: %v", release.Name, uninstallErr)
				} else {
					log.Printf("Failed to uninstall SlurmDeployment %s: %v", release.Name, uninstallErr)
					return ctrl.Result{}, uninstallErr
				}
			}

			// Remove our finalizer from the list and update it
			release.ObjectMeta.Finalizers = utils.SplitHeadArray(release.ObjectMeta.Finalizers, SlurmDeploymentFinalizer)
			if updateStatusErr := r.Update(ctx, release); updateStatusErr != nil {
				return ctrl.Result{}, updateStatusErr
			}
		}
		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	// Add finalizer if it doesn't exist
	if !utils.CheckIfExistInArray(release.ObjectMeta.Finalizers, SlurmDeploymentFinalizer) {
		log.Printf("Adding finalizer to SlurmDeployment %s", release.Name)
		release.ObjectMeta.Finalizers = append(release.ObjectMeta.Finalizers, SlurmDeploymentFinalizer)
		if updateStatusErr := r.Update(ctx, release); updateStatusErr != nil {
			return ctrl.Result{}, updateStatusErr
		}
		// Requeue to continue with installation after finalizer is added
		return ctrl.Result{Requeue: true}, nil
	}

	// build values yaml content for Slurm Chart
	chartValues := utils.BuildSlurmValues(&release.Spec.Values)

	// Check release if exists
	histClient := action.NewHistory(actionConfig)
	slurmChart := utils.DownloadChart(release.Spec.Chart.Name, release.Spec.Chart.Repository, release.Spec.Chart.Version)
	if _, getHistoryErr := histClient.Run(release.Name); getHistoryErr == nil {
		// upgrade release
		upgradeClient := action.NewUpgrade(actionConfig)
		upgradeClient.Namespace = release.Spec.Chart.Namespace

		if _, upgradeError := upgradeClient.Run(release.Name, slurmChart, chartValues); upgradeError != nil {
			log.Printf("Failed to upgrade release %s in namespace [%s]: %v", release.Name, release.Spec.Chart.Namespace, upgradeError)
			return ctrl.Result{}, upgradeError
		} else {
			if _, updateStatusErr := r.UpdateReleaseStatus(ctx, release); updateStatusErr != nil {
				return ctrl.Result{}, updateStatusErr
			}
		}
		return ctrl.Result{}, getHistoryErr
	} else {
		log.Printf("Cannot find release %s in namespace [%s] : %v", release.Name, release.Spec.Chart.Namespace, getHistoryErr)
		// install a new release
		installClient := action.NewInstall(actionConfig)
		installClient.ReleaseName = release.Name
		installClient.Namespace = release.Spec.Chart.Namespace

		if _, installErr := installClient.Run(slurmChart, chartValues); installErr != nil {
			log.Printf("Failed to install release %s in namespace [%s]: %v", release.Name, release.Spec.Chart.Namespace, installErr)
			return ctrl.Result{}, installErr
		} else {
			if _, updateStatusErr := r.UpdateReleaseStatus(ctx, release); updateStatusErr != nil {
				return ctrl.Result{}, updateStatusErr
			}
		}
	}

	return ctrl.Result{}, nil
}

// UpdateReleaseStatus updates the SlurmDeployment status with node counts and saves to Kubernetes
func (r *SlurmDeploymentReconciler) UpdateReleaseStatus(ctx context.Context, release *slurmv1.SlurmDeployment) (ctrl.Result, error) {
	cpuSTS, _ := r.RetrieveStatefulSetInfo(ctx, release, "-slurmd-cpu")

	gpuSTS, _ := r.RetrieveStatefulSetInfo(ctx, release, "-slurmd-gpu")

	controldSTS, _ := r.RetrieveStatefulSetInfo(ctx, release, "-slurmctld")

	databasedSTS, _ := r.RetrieveStatefulSetInfo(ctx, release, "-slurmdbd")

	mariadbSTS, _ := r.RetrieveStatefulSetInfo(ctx, release, "-mariadb")

	loginNodeDeploy, _ := r.RetrieveDeployInfo(ctx, release, "-login")

	// Update CPU node count
	release.Status.CPUNodeCount = fmt.Sprintf("%d/%d", cpuSTS.Status.ReadyReplicas, cpuSTS.Status.Replicas)
	// Update GPU node count
	release.Status.GPUNodeCount = fmt.Sprintf("%d/%d", gpuSTS.Status.ReadyReplicas, gpuSTS.Status.Replicas)
	// Update controld node count
	release.Status.ControldDeamonCount = fmt.Sprintf("%d/%d", controldSTS.Status.ReadyReplicas, controldSTS.Status.Replicas)
	// Update database node count
	release.Status.DatabaseDeamonCount = fmt.Sprintf("%d/%d", databasedSTS.Status.ReadyReplicas, databasedSTS.Status.Replicas)
	// Update maridb node count
	release.Status.MariadbServiceCount = fmt.Sprintf("%d/%d", mariadbSTS.Status.ReadyReplicas, mariadbSTS.Status.Replicas)
	// Update login node count
	release.Status.LoginNodeCount = fmt.Sprintf("%d/%d", loginNodeDeploy.Status.AvailableReplicas, loginNodeDeploy.Status.Replicas)
	// Show the command
	release.Status.JobCommand = strings.Join(append(release.Spec.Job.Command, release.Spec.Job.Args...), " ")

	// Update the status in Kubernetes
	if updateStatusErr := r.Status().Update(ctx, release); updateStatusErr != nil {
		log.Printf("Failed to update status: %v", updateStatusErr)
		return ctrl.Result{}, updateStatusErr
	}
	return ctrl.Result{}, nil
}

func (r *SlurmDeploymentReconciler) CreateNamespaceIfNotExist(ctx context.Context, namespace string) (ctrl.Result, error) {
	if namespace != "" {
		// 尝试获取命名空间
		if err := r.Get(ctx, client.ObjectKey{Name: namespace}, &corev1.Namespace{}); err != nil {
			// 仅当命名空间不存在时才创建
			if apierrors.IsNotFound(err) {
				if createErr := r.Create(ctx, &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{Name: namespace},
				}); createErr != nil {
					log.Printf("Failed to create namespace [%s]: %v", namespace, createErr)
					return ctrl.Result{}, createErr
				}
				log.Printf("Namespace [%s] created successfully", namespace)
			} else {
				// 其他错误（如权限问题）才返回错误
				log.Printf("Failed to get namespace [%s]: %v", namespace, err)
				return ctrl.Result{}, err
			}
		} else {
			// 命名空间已存在，无需操作
			log.Printf("Namespace [%s] already exists", namespace)
		}
	}
	return ctrl.Result{}, nil
}

func (r *SlurmDeploymentReconciler) RetrieveStatefulSetInfo(ctx context.Context, release *slurmv1.SlurmDeployment, suffix string) (appsv1.StatefulSet, error) {
	sts := &appsv1.StatefulSet{}
	if k8sGetInfoErr := r.Client.Get(ctx, types.NamespacedName{
		Name:      release.Name + suffix,
		Namespace: release.Spec.Chart.Namespace,
	}, sts); k8sGetInfoErr != nil {
		log.Printf("Failed to get StatefulSet: %v", k8sGetInfoErr)
		return *sts, k8sGetInfoErr
	}
	return *sts, nil
}

func (r *SlurmDeploymentReconciler) RetrieveDeployInfo(ctx context.Context, release *slurmv1.SlurmDeployment, suffix string) (appsv1.Deployment, error) {
	deploy := &appsv1.Deployment{}
	if k8sGetInfoErr := r.Client.Get(ctx, types.NamespacedName{
		Name:      release.Name + suffix,
		Namespace: release.Spec.Chart.Namespace,
	}, deploy); k8sGetInfoErr != nil {
		log.Printf("Failed to get Deployment : %v", k8sGetInfoErr)
		return *deploy, k8sGetInfoErr
	}
	return *deploy, nil
}

// SetupWithManager sets up the controller with the Manager.
// Helper functions to check and remove string from a slice of strings.
func (r *SlurmDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&slurmv1.SlurmDeployment{}).
		Complete(r)
}
