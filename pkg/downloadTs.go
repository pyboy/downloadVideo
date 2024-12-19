package pkg

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func DownloadTs(tsUrl string, tsFile string) {
	// 找到最后一个 '/' 的位置
	lastSlashIndex := strings.LastIndex(tsUrl, "/")

	// 如果找到了 '/', 则进行替换
	if lastSlashIndex != -1 {
		// 替换最后一个 '/' 之后的部分为 "123.ts"
		tsUrl = tsUrl[:lastSlashIndex+1] + tsFile
	}

	response, err := NewHttpGet(tsUrl)
	if err != nil {
		fmt.Println("Error downloading .ts file:", err)
	}
	defer response.Body.Close()

	// 创建临时TS文件
	tsTempFile, err := os.Create(tsFile)
	if err != nil {
		fmt.Println("Error creating temporary TS file:", err)

	}
	defer tsTempFile.Close()

	// 写入下载的数据到临时TS文件
	if _, err := io.Copy(tsTempFile, response.Body); err != nil {
		fmt.Println("Error writing to temporary TS file:", err)

	}

}
