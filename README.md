go mod tidy

## 下载
go get github.com/seamounts/go-mod-tidy

## 使用
1. 在使用 go mod 包管理的项目下执行改命令 `go-mod-tidy`，会基于 `go mod tidy` 下载缺失的包，如果因为墙的原因导致某些包无法下载，则会自动转换为 github 上对应的包。格式如下：
```shell
Set Follow Replace To go.mod:
cloud.google.com/go v0.26.0 => github.com/googleapis/google-cloud-go v0.26.0
google.golang.org/appengine v1.1.0 => github.com/golang/appengine v1.1.0
golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3 => github.com/golang/lint v0.0.0-20190313153728-d0100b6bd8b3
google.golang.org/genproto v0.0.0-20180817151627-c66870c02cf8 => github.com/google/go-genproto v0.0.0-20180817151627-c66870c02cf8
golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
```

2. 只需要讲 以上命令输出的 需要替换的包，拷贝到 `go.mod` 的 `replace` 模块。保存后，继续执行 `go-mod-tidy`， 直到所有依赖保下载成功。