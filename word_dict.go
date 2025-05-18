package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const dataFile = "words.csv"

// Word 结构体表示一个单词及其各种属性
type Word struct {
	Term          string   // 英文单词
	Translation   string   // 中文翻译
	PartOfSpeech  string   // 词性
	Pronunciation string   // 发音
	Example       string   // 例句
	ExampleTrans  string   // 例句翻译
	Tags          []string // 标签
	Difficulty    int      // 难度等级(1-5)
	Notes         string   // 备注
}

var words []Word

// 加载单词数据
func loadWords() error {
	file, err := os.Open(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在时返回空
		}
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		if len(record) < 9 {
			continue // 跳过不完整的记录
		}

		difficulty, _ := strconv.Atoi(record[7])
		tags := strings.Split(record[6], "|")
		if len(tags) == 1 && tags[0] == "" {
			tags = []string{}
		}

		word := Word{
			Term:          record[0],
			Translation:   record[1],
			PartOfSpeech:  record[2],
			Pronunciation: record[3],
			Example:       record[4],
			ExampleTrans:  record[5],
			Tags:          tags,
			Difficulty:    difficulty,
			Notes:         record[8],
		}
		words = append(words, word)
	}

	return nil
}

// 保存单词数据
func saveWords() error {
	file, err := os.Create(dataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, word := range words {
		record := []string{
			word.Term,
			word.Translation,
			word.PartOfSpeech,
			word.Pronunciation,
			word.Example,
			word.ExampleTrans,
			strings.Join(word.Tags, "|"),
			strconv.Itoa(word.Difficulty),
			word.Notes,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// 显示单词详情
func showWordDetails(word Word) {
	fmt.Println("\n--- 单词详情 ---")
	fmt.Printf("1. 英文单词: %s\n", word.Term)
	fmt.Printf("2. 中文翻译: %s\n", word.Translation)
	fmt.Printf("3. 词性: %s\n", word.PartOfSpeech)
	fmt.Printf("4. 发音: %s\n", word.Pronunciation)
	fmt.Printf("5. 例句: %s\n", word.Example)
	fmt.Printf("6. 例句翻译: %s\n", word.ExampleTrans)
	fmt.Printf("7. 标签: %s\n", strings.Join(word.Tags, ", "))
	fmt.Printf("8. 难度: %d/5\n", word.Difficulty)
	fmt.Printf("9. 备注: %s\n", word.Notes)
}

// 修改单词
func editWord() {
	if len(words) == 0 {
		fmt.Println("\n单词列表为空")
		return
	}

	// 先列出所有单词方便选择
	listWords()

	indexStr := inputWithPrompt("\n请输入要修改的单词编号: ")
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 1 || index > len(words) {
		fmt.Println("无效的编号")
		return
	}

	word := &words[index-1]
	showWordDetails(*word)

	// 选择要修改的字段
	fieldStr := inputWithPrompt("\n请输入要修改的字段编号(1-9, 0取消): ")
	field, err := strconv.Atoi(fieldStr)
	if err != nil || field < 0 || field > 9 {
		fmt.Println("无效的字段编号")
		return
	}

	if field == 0 {
		fmt.Println("取消修改")
		return
	}

	// 根据字段获取新值
	var newValue string
	switch field {
	case 1:
		newValue = inputWithPrompt("请输入新的英文单词: ")
		if newValue == "" {
			fmt.Println("英文单词不能为空")
			return
		}
		word.Term = newValue
	case 2:
		newValue = inputWithPrompt("请输入新的中文翻译: ")
		word.Translation = newValue
	case 3:
		newValue = inputWithPrompt("请输入新的词性: ")
		word.PartOfSpeech = newValue
	case 4:
		newValue = inputWithPrompt("请输入新的发音: ")
		word.Pronunciation = newValue
	case 5:
		newValue = inputWithPrompt("请输入新的例句: ")
		word.Example = newValue
	case 6:
		newValue = inputWithPrompt("请输入新的例句翻译: ")
		word.ExampleTrans = newValue
	case 7:
		newValue = inputWithPrompt("请输入新的标签(多个用逗号分隔): ")
		tags := strings.Split(newValue, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
		word.Tags = tags
	case 8:
		newValue = inputWithPrompt("请输入新的难度等级(1-5): ")
		difficulty, err := strconv.Atoi(newValue)
		if err == nil && difficulty >= 1 && difficulty <= 5 {
			word.Difficulty = difficulty
		} else {
			fmt.Println("无效的难度等级")
			return
		}
	case 9:
		newValue = inputWithPrompt("请输入新的备注: ")
		word.Notes = newValue
	}

	if err := saveWords(); err != nil {
		fmt.Printf("保存失败: %v\n", err)
		return
	}
	fmt.Println("单词修改成功!")
}

// 添加单词
func addWord() {
	fmt.Println("\n--- 添加新单词 ---")

	term := inputWithPrompt("英文单词: ")
	if term == "" {
		fmt.Println("单词不能为空")
		return
	}

	// 检查单词是否已存在
	for _, w := range words {
		if w.Term == term {
			fmt.Println("该单词已存在")
			return
		}
	}

	translation := inputWithPrompt("中文翻译: ")
	partOfSpeech := inputWithPrompt("词性(n/v/adj/adv等): ")
	pronunciation := inputWithPrompt("发音(可选): ")
	example := inputWithPrompt("例句(可选): ")
	exampleTrans := inputWithPrompt("例句翻译(可选): ")
	tagsInput := inputWithPrompt("标签(多个用逗号分隔，可选): ")
	difficultyInput := inputWithPrompt("难度等级(1-5，可选): ")
	notes := inputWithPrompt("备注(可选): ")

	// 处理难度等级
	difficulty := 0
	if difficultyInput != "" {
		fmt.Sscanf(difficultyInput, "%d", &difficulty)
		if difficulty < 1 || difficulty > 5 {
			difficulty = 0
		}
	}

	// 处理标签
	var tags []string
	if tagsInput != "" {
		tags = strings.Split(tagsInput, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	newWord := Word{
		Term:          term,
		Translation:   translation,
		PartOfSpeech:  partOfSpeech,
		Pronunciation: pronunciation,
		Example:       example,
		ExampleTrans:  exampleTrans,
		Tags:          tags,
		Difficulty:    difficulty,
		Notes:         notes,
	}

	words = append(words, newWord)
	if err := saveWords(); err != nil {
		fmt.Printf("保存失败: %v\n", err)
		return
	}
	fmt.Println("单词添加成功!")
}

// 删除单词
func removeWord() {
	if len(words) == 0 {
		fmt.Println("\n单词列表为空")
		return
	}

	listWords()

	indexStr := inputWithPrompt("\n请输入要删除的单词编号: ")
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 1 || index > len(words) {
		fmt.Println("无效的编号")
		return
	}

	term := words[index-1].Term
	confirm := inputWithPrompt(fmt.Sprintf("确定要删除单词 '%s' 吗?(y/n): ", term))
	if strings.ToLower(confirm) != "y" {
		fmt.Println("取消删除")
		return
	}

	words = append(words[:index-1], words[index:]...)
	if err := saveWords(); err != nil {
		fmt.Printf("保存失败: %v\n", err)
		return
	}
	fmt.Println("单词删除成功")
}

// 列出所有单词
func listWords() {
	if len(words) == 0 {
		fmt.Println("\n单词列表为空")
		return
	}

	fmt.Println("\n--- 单词列表 ---")
	for i, word := range words {
		fmt.Printf("%d. %s [%s] /%s/\n", i+1, word.Term, word.PartOfSpeech, word.Pronunciation)
		fmt.Printf("   翻译: %s\n", word.Translation)
		if word.Example != "" {
			fmt.Printf("   例句: %s\n", word.Example)
			if word.ExampleTrans != "" {
				fmt.Printf("   例句翻译: %s\n", word.ExampleTrans)
			}
		}
		if len(word.Tags) > 0 {
			fmt.Printf("   标签: %s\n", strings.Join(word.Tags, ", "))
		}
		if word.Difficulty > 0 {
			fmt.Printf("   难度: %d/5\n", word.Difficulty)
		}
		if word.Notes != "" {
			fmt.Printf("   备注: %s\n", word.Notes)
		}
		fmt.Println()
	}
}

// 搜索单词
func searchWord() {
	query := inputWithPrompt("\n请输入要搜索的单词或部分内容: ")
	if query == "" {
		fmt.Println("搜索内容不能为空")
		return
	}

	var results []Word
	for _, word := range words {
		if strings.Contains(strings.ToLower(word.Term), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(word.Translation), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(word.Notes), strings.ToLower(query)) {
			results = append(results, word)
		}
	}

	if len(results) == 0 {
		fmt.Println("没有找到匹配的单词")
		return
	}

	fmt.Println("\n--- 搜索结果 ---")
	for i, word := range results {
		fmt.Printf("%d. %s [%s] - %s\n", i+1, word.Term, word.PartOfSpeech, word.Translation)
	}
}

func main() {
	get_config()
	fmt.Println("欢迎使用单词记忆程序 (CSV存储版)")

	if err := loadWords(); err != nil {
		fmt.Printf("加载单词失败: %v\n", err)
		return
	}

	if CONFIG.UI_Mode == "cli mode" {
		cli_mode()
	} else if CONFIG.UI_Mode == "tui mode" {
		tui_mode()
	}
}
