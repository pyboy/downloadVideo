package pkg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// 判断html文件是否存在
func IsDireToBody(fileName string, urls string) (string, string, error) {
	var body string
	murls := urls
	_, err := os.Stat(fileName)
	if err == nil {
		// 文件不存在，打印消息并跳过
		fmt.Println("从文件获取数据...")
		body, err = ReadMFile(fileName)
		if err != nil {
			fmt.Println("Error:", err)
			return "", "", err
		}
	} else {
		// 文件存在，打印消息并跳过
		fmt.Println("正在获取网页数据...", murls)
		body, err = GetWebPageData(murls)
		if err != nil {
			fmt.Println("Error:", err)
			return "", "", err
		}

		hfile, err := os.Create(fileName)
		if err != nil {
			fmt.Println("Error creating m3u8:", err)
			return "", "", err
		}
		defer hfile.Close()

		if strings.Contains(urls, "m3u8") {
			if strings.Contains(body, "m3u8") {

				re := regexp.MustCompile(`.*.m3u8`)
				matches := re.FindAllStringSubmatch(body, -1)

				html := strings.LastIndex(urls, "/")
				murls = urls[:html+1] + matches[0][0]

				body, err = GetWebPageData(murls)
				if err != nil {
					fmt.Println("Error:", err)
					return "", "", err
				}
				fmt.Println("正在处理m3u8文件...", murls)
				body = "murls: " + string(murls) + "\n" + body
			}
		}

		_, err = hfile.WriteString(body)
		if err != nil {
			fmt.Println("Error writing to html:", err)
			return "", "", err
		}
		hfile.Close()
	}
	return body, murls, nil
}

func ReadMFile(filename string) (string, error) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", err
	}
	defer file.Close() // 确保文件在函数结束时关闭

	// 创建一个bufio.Reader来读取文件
	reader := bufio.NewReader(file)

	// 使用ReadString或ReadBytes方法读取文件内容
	// 这里我们使用ReadString，假设文件内容是以换行符分隔的
	var tsbody string
	for {
		line, err := reader.ReadString('\n') // 读取到换行符为止
		if err != nil && err != io.EOF {
			fmt.Println("Error reading file:", err)
			return "", err
		}
		tsbody += line // 将读取到的行追加到tsbody字符串
		if err == io.EOF {
			break // 如果到达文件末尾，则退出循环
		}
	}
	file.Close()
	return tsbody, err
}

// 转换字符编码
func ConvertEncoding(content [][]string) string {

	decodedURL := ""
	if len(content) > 0 && len(content[0]) > 1 {
		url := content[0][1]

		// 替换URL编码的百分号编码（%5C）为反斜杠（\）
		url = strings.ReplaceAll(url, "%5C", "\\")

		// 转换Unicode转义序列为UTF-8字符

		for len(url) > 0 {
			if url[0] == '\\' && strings.Contains(url, "\\u") {
				// 尝试解析Unicode转义序列
				runeValue, err := strconv.ParseInt(url[2:6], 16, 32)
				if err == nil {
					decodedURL += string(rune(runeValue))
					url = url[6:]
					continue
				}
			}
			decodedURL += string(url[0])
			url = url[1:]
		}
	}
	decodedURL = strings.ReplaceAll(decodedURL, "\\", "")
	return decodedURL
}

func GetName(decodedURL string) string {
	// 找到最后一个'/'和'/'之间的内容
	lastSlashIndex := strings.LastIndex(decodedURL, "/")
	prefix := decodedURL[:lastSlashIndex]

	// 找到倒数第二个'/'和最后一个'/'之间的内容
	secondLastSlashIndex := strings.LastIndex(prefix, "/")
	episode := decodedURL[secondLastSlashIndex+1 : lastSlashIndex]

	fmt.Println("提取的集数是:", episode)
	return episode
}

// func SortTsFiles(tsFiles []string) []string {
// 	// 使用map来存储文件名和索引
// 	files := []string{"0000.ts", "0002.ts", "0004.ts", "0001.ts", "0003.ts"}

// 	// 使用 sort.Slice 对切片进行排序
// 	sort.Slice(files, func(i, j int) bool {
// 		// 将字符串转换为整数，用于比较
// 		// 假设文件名格式为 "XXXX.ts"，其中 XXXX 是数字
// 		numI, _ := strconv.Atoi(files[i][:4])
// 		numJ, _ := strconv.Atoi(files[j][:4])
// 		return numI < numJ
// 	})

// 	return files
// }
