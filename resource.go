package main

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