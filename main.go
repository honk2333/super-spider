package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// 资源类型定义
const (
	ResourceTypeVideo = "video"
	ResourceTypeImage = "image"
	ResourceTypeDocument = "document"
)

// 资源结构体
type Resource struct {
	URL  string
	Type string
}

// 主函数
func main() {
	if len(os.Args) < 2 {
		fmt.Println("请提供要爬取的网站链接")
		fmt.Println("使用方法: super-spider <url>")
		os.Exit(1)
	}

	url := os.Args[1]
	fmt.Printf("开始爬取网站: %s\n", url)

	// 创建保存目录
	if err := os.MkdirAll("downloads", 0755); err != nil {
		fmt.Printf("创建目录失败: %v\n", err)
		os.Exit(1)
	}

	// 获取网页内容
	content, err := fetchWebPage(url)
	if err != nil {
		fmt.Printf("获取网页内容失败: %v\n", err)
		os.Exit(1)
	}

	// 提取资源链接
	resources := extractResources(content, url)
	fmt.Printf("发现 %d 个资源\n", len(resources))

	// 下载资源
	downloadResources(resources)

	fmt.Println("爬取完成")
}

// 获取网页内容
func fetchWebPage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// 提取资源链接
func extractResources(content, baseURL string) []Resource {
	var resources []Resource

	// 提取视频链接
	videoRegex := regexp.MustCompile(`<video[^>]+src=["']([^"']+)["']`)
	videoMatches := videoRegex.FindAllStringSubmatch(content, -1)
	for _, match := range videoMatches {
		resources = append(resources, Resource{
			URL:  resolveURL(match[1], baseURL),
			Type: ResourceTypeVideo,
		})
	}

	// 提取图片链接
	imageRegex := regexp.MustCompile(`<img[^>]+src=["']([^"']+)["']`)
	imageMatches := imageRegex.FindAllStringSubmatch(content, -1)
	for _, match := range imageMatches {
		resources = append(resources, Resource{
			URL:  resolveURL(match[1], baseURL),
			Type: ResourceTypeImage,
		})
	}

	// 提取文档链接
	docRegex := regexp.MustCompile(`<a[^>]+href=["']([^"']+.(pdf|doc|docx|xls|xlsx|ppt|pptx))["']`)
	docMatches := docRegex.FindAllStringSubmatch(content, -1)
	for _, match := range docMatches {
		resources = append(resources, Resource{
			URL:  resolveURL(match[1], baseURL),
			Type: ResourceTypeDocument,
		})
	}

	return resources
}

// 解析相对URL为绝对URL
func resolveURL(url, baseURL string) string {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}

	if strings.HasPrefix(url, "/") {
		// 处理绝对路径
		parsedBase, _ := url.Parse(baseURL)
		return parsedBase.Scheme + "://" + parsedBase.Host + url
	}

	// 处理相对路径
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return baseURL + url
}

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

SPSC: INIT \
\
Co-authored-by: honk2333 whk19981229@163.com \
Co-authored-by: cameron314 no-reply