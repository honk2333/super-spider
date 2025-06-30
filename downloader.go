package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// 下载资源
func downloadResources(resources []Resource) {
	for i, resource := range resources {
		fmt.Printf("正在下载资源 %d/%d: %s\n", i+1, len(resources), resource.URL)

		// 创建保存路径
		outputDir := fmt.Sprintf("downloads/%s", resource.Type)
		os.MkdirAll(outputDir, 0755)

		// 提取文件名
		parts := strings.Split(resource.URL, "/")
		filename := parts[len(parts)-1]
		outputPath := fmt.Sprintf("%s/%s", outputDir, filename)

		// 使用ffmpeg下载资源
		cmd := exec.Command("ffmpeg", "-i", resource.URL, "-c", "copy", outputPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("下载失败: %v\n输出: %s\n", err, string(output))
			continue
		}

		fmt.Printf("下载成功: %s\n", outputPath)
	}
}

