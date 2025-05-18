package main

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

const (
	CONFIG_FILE = "./config.lua"
)

type config struct {
	WordFilePath string
	UI_Mode      string
}

var CONFIG config

func get_config() {
	L := lua.NewState()
	defer L.Close()

	if err := L.DoFile(CONFIG_FILE); err != nil {
		fmt.Println("Error loading config file:", err)
		return
	}

	wordFilePath := L.GetGlobal("WordFilePath")
	if wordFilePath.Type() == lua.LTString {
		CONFIG.WordFilePath = wordFilePath.String()
	} else {
		fmt.Println("WordFilePath应该是字符串")
	}

	UI_Mode := L.GetGlobal("UI_Mode")
	if UI_Mode.Type() == lua.LTString {
		CONFIG.UI_Mode = UI_Mode.String()
	} else {
		fmt.Println("UI_Mode应该是一个字符串")
	}
}
