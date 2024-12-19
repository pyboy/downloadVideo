package main

import (
	"downloadvideo/pkg"
	"fmt"
	"strings"
)

var downloadFuncs = map[string]func(string) error{
	"nxzsxh.com":      pkg.DownloadVideo, // 下载可以通过网站地址直接获取m3u8信息
	"www.cnjiean.com": pkg.DownloadVideo,
	"www.58ys01.com":  pkg.DownloadVideo,

	// 可以继续添加其他域名和下载函数的对应关系  https://www.58ys01.com/vod/play/id/11519/sid/1/nid/1.html
}

func main() {
	// 启动窗口并传递回调函数
	pkg.StartWindow(func(url string) {
		fmt.Println("您输入的网址是:", url)

		for domain, downloadFunc := range downloadFuncs {
			if strings.Contains(url, domain) {
				err := downloadFunc(url)
				if err != nil {
					pkg.AppendLogText(fmt.Sprintf("Error: %v", err))
					return
				}
				pkg.AppendLogText("下载完成")
			}
		}
	})
	select {}
}
