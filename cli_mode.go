package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 显示菜单
func showMenu() {
	fmt.Println("\n=== 单词记忆程序 ===")
	fmt.Println("1. 添加单词")
	fmt.Println("2. 删除单词")
	fmt.Println("3. 修改单词")
	fmt.Println("4. 列出所有单词")
	fmt.Println("5. 搜索单词")
	fmt.Println("6. 退出")
}

// 获取用户输入
func inputWithPrompt(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func cli_mode() {
	for {
		showMenu()
		choice := inputWithPrompt("请选择操作(1-6): ")

		switch choice {
		case "1":
			addWord()
		case "2":
			removeWord()
		case "3":
			editWord()
		case "4":
			listWords()
		case "5":
			searchWord()
		case "6":
			fmt.Println("感谢使用，再见！")
			return
		default:
			fmt.Println("无效选项，请重新输入")
		}
	}
}
