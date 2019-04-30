# go-mod-tidy

## 下载
go get github.com/seamounts/go-mod-tidy

## 使用
1. 在使用 go mod 包管理的项目下执行改命令 `go-mod-tidy`，会基于 `go mod tidy` 下载缺失的包，如果因为墙的原因导致某些包无法下载，则会自动替换为 github 上对应的包。格式如下：
    ```shell
    cloud.google.com/go v0.26.0 => github.com/googleapis/google-cloud-go v0.26.0
    google.golang.org/appengine v1.1.0 => github.com/golang/appengine v1.1.0
    golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3 => github.com/golang/lint v0.0.0-20190313153728-d0100b6bd8b3
    ```