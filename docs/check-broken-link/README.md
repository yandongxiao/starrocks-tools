# 检查链接是否会失效的工具
这是一个简单的 Go 程序，可以检查一个目录下的所有 Markdown 文件中的链接是否有效，如果发现无效的链接，会输出到控制台。

## 安装

1. 安装 Golang。要运行这个程序，你需要先安装 Golang 的编译环境。你可以从以下两个网站下载 Golang 的安装包：https://golang.google.cn/dl/
2. 运行程序。要运行这个程序，你需要先将代码文件下载到你的本地目录。然后打开终端或命令行窗口，切换到代码文件所在的目录。输入以下命令：`go run main <directory>`。其中 `<directory>` 是你要检查的目录的路径

## 运行

执行 `go run main /home/disk2/bob/docs.zh-cn` 后，运行结果如下：

```bash
/home/disk2/bob/docs.zh-cn/administration/Query_planning.md link s broken:
        ../using_starrocks/Colocation_join.md
...        
```

它表示，在检查 `docs.zh-cn/administration/Query_planning.md` 文档时，内部有一个超链接，它的地址是 `../using_starrocks/Colocation_join.md`，但是这个地址 `../using_starrocks/Colocation_join.md` 是不存在或者格式不合法的。



