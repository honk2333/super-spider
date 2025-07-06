package downloader

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"super-spider/consts"
)

// 下载资源
func DownloadResources(resources []consts.Resource) {
	for i, res := range resources {
		fmt.Printf("正在下载资源 %d/%d: %s\n", i+1, len(resources), res.URL)

		// 创建保存路径
		outputDir := fmt.Sprintf("downloads/%s", res.Type)
		os.MkdirAll(outputDir, 0755)

		// 提取文件名
		parts := strings.Split(res.URL, "/")
		filename := parts[len(parts)-1]
		outputPath := fmt.Sprintf("%s/%s", outputDir, filename)

		// 使用ffmpeg下载资源
		cmd := exec.Command("ffmpeg", "-i", res.URL, "-c", "copy", outputPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("下载失败: %v\n输出: %s\n", err, string(output))
			continue
		}

		fmt.Printf("下载成功: %s\n", outputPath)
	}
}
