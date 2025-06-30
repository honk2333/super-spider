package main

import (
	"fmt"
	"os"
)

// 主函数
func main() {
	if len(os.Args) < 2 {
		fmt.Println("请提供要爬取的网站链接")
		fmt.Println("使用方法: super-spider <url>")
		os.Exit(1)
	}

	targetURL := os.Args[1]
	fmt.Printf("开始爬取网站: %s\n", targetURL)

	// 创建保存目录
	if err := os.MkdirAll("downloads", 0755); err != nil {
		fmt.Printf("创建目录失败: %v\n", err)
		os.Exit(1)
	}

	// 获取网页内容
	content, err := fetchWebPage(targetURL)
	if err != nil {
		fmt.Printf("获取网页内容失败: %v\n", err)
		os.Exit(1)
	}

	// 提取资源链接
	resources := extractResources(content, targetURL)
	fmt.Printf("发现 %d 个资源\n", len(resources))

	// 下载资源
	downloadResources(resources)

	fmt.Println("爬取完成")
}