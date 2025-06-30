package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

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

	// fmt.Println("获取网页内容成功", string(body));
	return string(body), nil
}

func extractVideos(content, baseURL string) []Resource {
	// 提取视频链接
	var resources []Resource
	videoRegex := regexp.MustCompile(`<video[^>]+src=["']([^"']+)["']`)
	videoMatches := videoRegex.FindAllStringSubmatch(content, -1)
	for _, match := range videoMatches {
		resources = append(resources, Resource{
			URL:  resolveURL(match[1], baseURL),
			Type: ResourceTypeVideo,
		})
	}
	return resources
}

// https://img.36krcdn.com/hsossms/20250630/v2_b3c7900a541249fc8103a3c8b52344b7@1743780481@ai_oswg912072oswg1053oswg495_img_png~tplv-1marlgjv7f-ai-v3:600:400:600:400:q70.jpg?x-oss-process=image/format,webp
func extractImages(content, baseURL string) []Resource {
	var resources []Resource

	// 提取图片链接
	imageRegex := regexp.MustCompile(`<img\s+[^>]*?src\s*=\s*["']([^"']+)["']`)
	imageMatches := imageRegex.FindAllStringSubmatch(content, -1)
	for _, match := range imageMatches {
		fmt.Println("提取图片链接", match[1])
		resources = append(resources, Resource{
			URL:  resolveURL(match[1], baseURL),
			Type: ResourceTypeImage,
		})
	}

	return resources
}

func extractDocs(content, baseURL string) []Resource {
	var resources []Resource

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

// 提取资源链接
func extractResources(content, baseURL string) []Resource {
	var resources []Resource
	resources = append(resources, extractVideos(content, baseURL)...)
	resources = append(resources, extractImages(content, baseURL)...)
	resources = append(resources, extractDocs(content, baseURL)...)
	return resources
}

// 解析相对URL为绝对URL
func resolveURL(relativeURL, baseURL string) string {
	if strings.HasPrefix(relativeURL, "http://") || strings.HasPrefix(relativeURL, "https://") {
		return relativeURL
	}

	if strings.HasPrefix(relativeURL, "/") {
		// 处理绝对路径
		parsedBase, _ := url.Parse(baseURL)
		return parsedBase.Scheme + "://" + parsedBase.Host + relativeURL
	}

	// 处理相对路径
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return baseURL + relativeURL
}
