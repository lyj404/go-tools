package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

// RequestData 定义客户端请求的数据结构
type RequestData struct {
	Input string `json:"input"` // 接收前端传来的原始字符串
}

// ResponseData 定义返回给客户端的响应数据结构
type ResponseData struct {
	Result string `json:"result"` // 转换后的驼峰命名结果
}

func main() {
	// 注册路由处理函数
	http.HandleFunc("/convert", convertHandler)
	http.Handle("/", http.FileServer(http.Dir(filepath.Join(".", "static"))))

	// 启动HTTP服务器
	fmt.Println("Server started at http://localhost:9999")
	log.Fatal(http.ListenAndServe(":9999", nil))
}

// convertHandler 处理驼峰命名转换请求
func convertHandler(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法是否为POST
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体中的JSON数据
	var req RequestData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 按换行符分割输入字符串
	lines := strings.Split(req.Input, "\n")
	// 存储转换后的结果
	var resultLines []string

	// 处理每一行输入
	for _, line := range lines {
		// 去除首尾空白字符
		line = strings.TrimSpace(line)
		if line == "" {
			continue // 跳过空行
		}
		// 转换为驼峰命名并添加到结果集
		resultLines = append(resultLines, toCamelCase(line))
	}

	// 构造响应数据
	response := ResponseData{
		// 用换行符连接转换后的结果
		Result: strings.Join(resultLines, "\n"),
	}

	// 设置响应头并返回JSON格式数据
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response) // 编码并发送响应
}

// toCamelCase 将字符串转换为驼峰命名的核心函数
func toCamelCase(input string) string {
	// 第一步：统一转为小写，处理全大写的特殊情况（如"AMTPOST"）
	input = strings.ToLower(input)

	// 第二步：使用正则表达式分割字符串
	// 匹配所有分隔符：下划线(_)、空格(\s)、连字符(-)
	re := regexp.MustCompile(`[_\s-]`)
	words := re.Split(input, -1) // 分割为单词切片

	// 第三步：转换每个单词（首字母大写）
	for i, word := range words {
		if i > 0 && word != "" { // 第一个单词保持全小写，后续单词首字母大写
			words[i] = strings.ToUpper(word[:1]) + word[1:] // 首字母大写+剩余部分
		}
	}

	// 第四步：合并所有单词，形成最终驼峰命名
	return strings.Join(words, "")
}
