package pkg

import (
	"fmt"
	"log"
	"os"
)

func WorkerDirectory(cwd string) {
	// 判断是否有目录，没有创建
	_, err := os.Stat(cwd)
	if os.IsNotExist(err) {
		err := os.Mkdir(cwd, 0755)
		if err != nil {
			fmt.Println("创建目录失败:", err)
		}
	}

	err = os.Chdir(cwd)
	if err != nil {
		log.Fatalf("Failed to change working directory: %v", err)
	}
	fmt.Println("改变当前工作目录到:", cwd)
}
