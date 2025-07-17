package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	chart "helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
)

// 实现 Chart 下载逻辑（从仓库获取）
func DownloadChart(chartName, repository, version string) *chart.Chart {

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
	if err := downloadFileFromURL(chartURL, filePath); err != nil {
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

// download chart File
func downloadFileFromURL(url, filePath string) error {
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

// unzip helm chart.tgz file
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
