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

	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/cli"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"helm.sh/helm/v3/pkg/chart/loader"

	slurmv1 "github.com/AaronYang0628/slurm-on-k8s/api/v1"
)

// SlurmDeploymentReconciler reconciles a SlurmDeployment object
type SlurmDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=slurm.ay.dev,resources=slurmdeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=slurm.ay.dev,resources=slurmdeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=slurm.ay.dev,resources=slurmdeployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SlurmDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.0/pkg/reconcile
func (r *SlurmDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// 获取自定义资源实例
	release := &slurmv1.SlurmDeployment{}
	if err := r.Get(ctx, req.NamespacedName, release); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Printf("Find SlurmDeployment %s", release.Name)

	// 初始化 Helm 配置
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), release.Spec.Chart.Namespace, "secret", log.Printf); err != nil {
		log.Printf("Failed to initialize Helm configuration: %v", err)
		return ctrl.Result{}, err
	}

	// 构造 Chart 的 values
	values := map[string]interface{}{
		"service": map[string]interface{}{
			"type": release.Spec.Values.Service.Type,
		},
		"replicaCount": release.Spec.Values.ReplicaCount,
	}

	// 检查 Release 是否存在
	histClient := action.NewHistory(actionConfig)
	if _, err := histClient.Run(release.Name); err == nil {
		// 执行升级
		upgradeClient := action.NewUpgrade(actionConfig)
		upgradeClient.Namespace = release.Spec.Chart.Namespace

		_, err = upgradeClient.Run(release.Name, getChart(release), values)
		return handleResult(err)
	}

	// 全新安装
	installClient := action.NewInstall(actionConfig)
	installClient.ReleaseName = release.Name
	installClient.Namespace = release.Spec.Chart.Namespace
	_, err := installClient.Run(getChart(release), values)
	return handleResult(err)
}

func handleResult(err error) (ctrl.Result, error) {
	if err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func getChart(r *slurmv1.SlurmDeployment) *chart.Chart {
	// 实现 Chart 下载逻辑（从仓库获取）
	// 可以使用 helm pull 命令或直接访问仓库
	chartName := r.Spec.Chart.Name
	repository := r.Spec.Chart.Repository
	version := r.Spec.Chart.Version

	// 构造 Chart 的下载 URL
	chartURL := fmt.Sprintf("%s/%s-%s.tgz", repository, chartName, version)
	log.Printf("Downloading chart from %s", chartURL)

	// 创建临时目录用于存储下载的 Chart 文件
	tempDir, err := os.MkdirTemp("", "helm-charts")
	if err != nil {
		log.Printf("Failed to create temporary directory: %v", err)
		return nil
	}
	defer os.RemoveAll(tempDir)

	// 下载 Chart 文件
	filePath := filepath.Join(tempDir, fmt.Sprintf("%s-%s.tgz", chartName, version))
	if err := downloadFile(chartURL, filePath); err != nil {
		log.Printf("Failed to download chart: %v", err)
		return nil
	} else {
		log.Printf("Downloaded chart to %s", filePath)
	}

	// 解压 Chart 文件
	if err := extractTarGz(filePath, tempDir); err != nil {
		log.Printf("Failed to extract chart: %v", err)
		return nil
	} else {
		log.Printf("Extracted chart to %s", tempDir)
	}

	// 加载 Chart
	chartPath := filepath.Join(tempDir, chartName)
	chrt, err := loader.Load(chartPath)
	if err != nil {
		log.Printf("Failed to load chart: %v", err)
		return nil
	} else {
		log.Printf("Loaded chart %s", chrt.Metadata.Name)
	}

	return chrt
}

// downloadFile 下载文件到指定路径
func downloadFile(url, filePath string) error {
	log.Printf("Downloading file from %s to %s", url, filePath)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error downloading file: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to download file: %s", resp.Status)
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return err
	}
	defer out.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		log.Printf("Error writing file: %v", err)
		return err
	}
	log.Printf("Downloaded %d bytes to %s", written, filePath)

	return nil
}

// extractTarGz 解压 .tgz 文件
func extractTarGz(filePath, destPath string) error {
	// 打开压缩文件
	gzFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer gzFile.Close()

	// 创建gzip reader
	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return fmt.Errorf("创建gzip reader失败: %v", err)
	}
	defer gzReader.Close()

	// 创建tar reader
	tarReader := tar.NewReader(gzReader)

	// 获取目标路径绝对地址用于安全检查
	destAbs, err := filepath.Abs(destPath)
	if err != nil {
		return fmt.Errorf("获取绝对路径失败: %v", err)
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // 正常结束
		}
		if err != nil {
			return fmt.Errorf("读取tar条目失败: %v", err)
		}

		// 构建目标路径并进行安全检查
		targetPath := filepath.Join(destPath, header.Name)
		targetAbs, err := filepath.Abs(targetPath)
		if err != nil {
			return fmt.Errorf("路径解析失败: %v", err)
		}

		// 防止路径穿越攻击
		if !strings.HasPrefix(targetAbs, destAbs) {
			return fmt.Errorf("危险路径检测: %s 试图访问目标路径外", header.Name)
		}

		// 根据文件类型处理
		switch header.Typeflag {
		case tar.TypeDir:
			// 创建目录并设置权限
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("创建目录失败: %v", err)
			}
			if err := os.Chmod(targetPath, os.FileMode(header.Mode).Perm()); err != nil {
				return fmt.Errorf("设置目录权限失败: %v", err)
			}

		case tar.TypeReg:
			// 确保父目录存在
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("创建父目录失败: %v", err)
			}

			// 创建文件并设置内容
			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode).Perm())
			if err != nil {
				return fmt.Errorf("创建文件失败: %v", err)
			}

			// 复制文件内容
			if _, err := io.CopyN(outFile, tarReader, header.Size); err != nil {
				outFile.Close()
				return fmt.Errorf("写入文件内容失败: %v", err)
			}

			// 设置文件修改时间
			if err := os.Chtimes(targetPath, time.Time{}, header.ModTime); err != nil {
				outFile.Close()
				return fmt.Errorf("设置修改时间失败: %v", err)
			}

			outFile.Close()

		default:
			// 跳过非常规文件类型（如符号链接、设备文件等）
			continue
		}
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SlurmDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&slurmv1.SlurmDeployment{}).
		Complete(r)
}
