package pkg

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func DownloadVideo(urls string) error {
	fmt.Println("开始启动下载程序")
	//urls := "https://www.cnjiean.com/pan/46051/1-1.html"

	// 网页名称
	html := strings.LastIndex(urls, "/")
	htmlName := urls[html+1:]

	fmt.Println("网页名称:", htmlName)

	// 工作目录
	pwd := strings.LastIndex(htmlName, ".")
	pwdName := htmlName[:pwd]
	// 设置工作目录
	WorkerDirectory(pwdName)

	// ffmpeg 二进制文件
	// 将ffmpeg二进制文件写入临时文件
	_, err := WriteFFmpegBinaryToTempFile()
	if err != nil {
		fmt.Println("无法写入ffmpeg二进制文件:", err)
		return err
	}

	// 判断html文件 生成m3u8
	htmlFile, _, err := IsDireToBody(htmlName, urls)
	if err != nil {
		fmt.Println("无法获取html文件内容:", err)
		return err
	}

	// 获取html中的m3u8
	re := regexp.MustCompile(`"(https?:[^"]+.m3u8)"`)
	matches := re.FindAllStringSubmatch(string(htmlFile), -1)
	fmt.Println("正在处理m3u8文件...", matches)

	m3u8Urls := ConvertEncoding(matches)
	fmt.Println("m3u8Urls:", m3u8Urls)

	// mp4 文件名称

	m3u8File, murls, err := IsDireToBody(htmlName+".m3u8", m3u8Urls)
	if err != nil {
		fmt.Println("无法获取m3u8文件内容:", err)
		return err
	}
	fmt.Println("获取m3u8文件内容成功...")

	// 创建tsfilelist.txt文件
	fileList := "tsfilelist.txt"
	file, err := os.Create(fileList)
	if err != nil {
		fmt.Println("Error creating filelist.txt:", err)
		return err
	}
	defer file.Close()

	// 清空文件
	file.Truncate(0)

	filelistContent := ""

	count := 0

	for _, line := range strings.Split(m3u8File, "\n") {
		if strings.HasSuffix(line, ".ts") {
			count++
		}
	}

	// 已存在的英文单词列表
	fileNum := 0

	fmt.Println("开始处理ts文件...")
	for _, line := range strings.Split(m3u8File, "\n") {
		if strings.Contains(line, "murls:") {
			// 获取 murls: 后面内容
			murls = strings.TrimPrefix(line, "murls:")
		}
		if strings.Contains(line, "/") {
			continue
		}

		if strings.HasSuffix(line, ".ts") {
			filelistContent += fmt.Sprintf("file '%s'\n", line)

			_, err := os.Stat("./" + line)
			if err == nil {
				// 文件存在，打印消息并跳过

				fileNum += 1
				continue
			}

			fmt.Println("开始下载ts文件: " + line + "\n")

			AppendLogText("下载文件数：" + fmt.Sprintf("%d", fileNum) + "/" + fmt.Sprintf("%d", count))

			DownloadTsFile(murls, line) // 下载ts文件
			fileNum += 1
		}

	}
	// 写入ts列表信息
	_, err = file.WriteString(filelistContent)
	if err != nil {
		fmt.Println("Error writing to filelist.txt:", err)

		return err
	}
	file.Close()
	mp4Name := GetName(m3u8Urls)
	// 合并ts视频文件
	MergeMP4Files(fileList, mp4Name)

	// 尝试删除一个非空目录及其所有内容
	err = os.RemoveAll(pwdName)
	if err != nil {
		fmt.Println("Failed to remove directory: %v", err)
	} else {
		fmt.Println("Directory and all its contents removed successfully")
	}

	return nil
}
