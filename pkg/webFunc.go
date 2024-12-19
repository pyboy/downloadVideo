package pkg

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func NewHttpGet(url string) (*http.Response, error) {
	client := &http.Client{}
	trimmedStr := strings.TrimSpace(url)
	request, err := http.NewRequest("GET", trimmedStr, nil)
	if err != nil {
		fmt.Println("Error creating GET request:", err)
		return nil, err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending GET request:", err)
		return nil, err
	}

	return response, nil
}

func GetWebPageData(url string) (string, error) {

	response, err := NewHttpGet(url)
	if err != nil {
		fmt.Println("Error sending GET request:", err)
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	return string(body), err
}

func DownloadTsFile(urls string, ts string) {
	// 找到最后一个 '/' 的位置
	lastSlashIndex := strings.LastIndex(urls, "/")

	// 如果找到了 '/', 则进行替换
	if lastSlashIndex != -1 {
		// 替换最后一个 '/' 之后的部分为 "123.ts"
		urls = urls[:lastSlashIndex+1] + ts
	}

	response, err := NewHttpGet(urls)
	if err != nil {
		fmt.Println("Error downloading .ts file:", err)
		return
	}

	// 创建临时TS文件
	tsFile, err := os.Create(ts)
	if err != nil {
		fmt.Println("Error creating temporary TS file:", err)
		return
	}
	defer tsFile.Close()

	// 写入下载的数据到临时TS文件
	if _, err := io.Copy(tsFile, response.Body); err != nil {
		fmt.Println("Error writing to temporary TS file:", err)
		return
	}
	tsFile.Close()

}

func StartFFmpeg(concatFileName string, output *os.File) {
	cmd := exec.Command(
		"ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", concatFileName,
		"-c", "copy",
		"-bsf:a", "aac_adtstoasc",
		"-f", "mp4",
		"-",
	)
	cmd.Stdout = output

	// 开始执行ffmpeg命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting ffmpeg:", err)
		return
	}

	// 等待ffmpeg命令完成
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for ffmpeg:", err)
		return
	}

	fmt.Println("Merging completed.")
}
