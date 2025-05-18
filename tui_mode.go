package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	menu   []string
	cursor int
}

type Model interface {
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "ESC":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j": // 修正拼写
			if m.cursor < len(m.menu)-1 {
				m.cursor++
			}

		case "enter", " ":
			tuiAction = m.cursor
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	ui := "operation list:\n\n"

	for i, choice := range m.menu {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		ui += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	return ui
}

var initModel = model{
	menu: []string{"添加单词", "删除单词", "修改单词", "列出所有单词", "搜索单词", "退出"},
}

var tuiAction = 0

func tui_mode() {
	for {
		cmd := tea.NewProgram(initModel)
		if err := cmd.Start(); err != nil {
			fmt.Println("start failed:", err)
			os.Exit(1)
		}
		switch tuiAction {
		case 0:
			addWord()
		case 1:
			removeWord()
		case 2:
			editWord()
		case 3:
			listWords()
		case 4:
			searchWord()
		case 5:
			fmt.Println("退出")
			os.Exit(0)
		default:
			fmt.Println("无效的操作")
			os.Exit(1)
		}
	}
}
