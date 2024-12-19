package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// LogText 导出的日志文本框
var LogText *widget.Label

// ScrollContainer 全局变量
var ScrollContainer *container.Scroll

// 历史记录文件路径
const historyFilePath = "history.json"

// 历史记录
var history []string

// 读取历史记录
func loadHistory() ([]string, error) {
	data, err := ioutil.ReadFile(historyFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	var h []string
	err = json.Unmarshal(data, &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// 保存历史记录
func saveHistory(history []string) error {
	data, err := json.Marshal(history)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(historyFilePath, data, 0644)
}

// StartWindow 启动窗口并返回输入的URL
func StartWindow(onStart func(url string)) {
	// 读取历史记录
	history, err := loadHistory()
	if err != nil {
		fmt.Println("加载历史记录失败:", err)
		history = []string{}
	}

	// 创建一个新的应用实例
	myApp := app.New()
	themes := theme.DefaultTheme()
	// 设置应用标题
	myApp.Settings().SetTheme(themes)

	// 创建一个新的窗口
	myWindow := myApp.NewWindow("网页下载器")

	// 设置窗口默认大小
	myWindow.Resize(fyne.NewSize(800, 600))

	// 创建一个文本框用于显示日志 超过最高限制 使用滚动条显示
	LogText = widget.NewLabel("")

	// 创建一个输入框
	input := widget.NewEntry()
	input.SetPlaceHolder("请输入URL")

	// 创建一个按钮
	button := widget.NewButton("开始", func() {
		// 获取输入框中的内容
		url := input.Text
		if url != "" {
			// 添加到历史记录
			history = append(history, url)
			// 去重历史记录
			history = unique(history)
			// 保存历史记录
			err := saveHistory(history)
			if err != nil {
				fmt.Println("保存历史记录失败:", err)
			}
			// 更新下拉菜单
			updateHistoryMenu(input, history, myWindow)
			// 调用回调函数
			onStart(url)
		}
	})

	// 创建一个下拉菜单用于选择历史记录
	updateHistoryMenu(input, history, myWindow)

	// 将输入框和按钮水平排列
	inputButtonContainer := container.NewHSplit(
		input,
		button,
	)

	inputButtonContainer.Offset = 0.9 // 输入框占80%，按钮占20%

	// 使用 ScrollContainer 包裹 LogText
	scrollContainer := container.NewVScroll(LogText)

	// 使用 NewBorder 布局，将输入框和按钮放在顶部，日志文本框占据剩余的空间
	borderContainer := container.NewBorder(inputButtonContainer, nil, nil, nil, scrollContainer)

	// 将布局容器设置为主容器
	myWindow.SetContent(container.NewVBox(
		borderContainer,
	))

	// 显示并运行应用
	myWindow.ShowAndRun()

	// 处理输入框文本变化事件
	input.OnChanged = func(text string) {
		// 这里可以添加文本变化时的处理逻辑
		fmt.Println("输入框文本变化:", text)
	}

	// 处理输入框提交事件
	input.OnSubmitted = func(text string) {
		// 这里可以添加提交时的处理逻辑
		fmt.Println("输入框提交:", text)
	}
}

// AppendLogText 向日志文本框追加内容并滚动到底部
func AppendLogText(content string) {
	LogText.SetText(content + "\n")
	// 强制刷新界面
	LogText.Refresh()
}

// 去重历史记录
func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// 更新下拉菜单
func updateHistoryMenu(input *widget.Entry, history []string, myWindow fyne.Window) {
	// 创建一个新的菜单
	menu := fyne.NewMenu("")
	for _, h := range history {
		item := fyne.NewMenuItem(h, func() {
			input.SetText(h)
			input.Refresh() // 确保输入框更新显示文本
		})
		menu.Items = append(menu.Items, item)
	}

	// 创建一个新的弹出菜单
	historyMenu := widget.NewPopUpMenu(menu, myWindow.Canvas())
	fmt.Println(historyMenu)
}
