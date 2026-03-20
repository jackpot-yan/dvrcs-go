package agent

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"
)

// 先模拟AI返回
func SimulateAI(prompt string) []string {
	response := map[string]string{
		"default": "我收到了你的消息。这是一个模拟的流式响应，文字会像真实 AI 一样逐字出现。你可以将这里替换为真正的大模型 API 调用，例如 OpenAI、Anthropic 或本地的 Ollama。",
	}
	prompt = strings.ToLower(strings.TrimSpace(prompt))
	reply := response["default"]
	for key, val := range map[string]string{
		"你好":    "你好！我是一个支持流式输出的 AI 助手。有什么我可以帮你的吗？",
		"hello": "Hello! I'm a streaming AI assistant. How can I help you today?",
		"时间":    fmt.Sprintf("当前服务器时间是：%s", time.Now().Format("2006年01月02日 15:04:05")),
		"go":    "Go 语言非常适合构建高性能的 SSE 服务！它的 goroutine 和 channel 让并发处理变得优雅，`net/http` 标准库就能轻松实现流式响应，不需要任何第三方框架。",
		"sse":   "Server-Sent Events（SSE）是一种服务器向客户端单向推送数据的技术。相比 WebSocket，它更轻量，基于普通 HTTP，天然支持断线重连，非常适合 AI 流式输出场景。",
	} {
		if strings.Contains(prompt, key) {
			reply = val
			break
		}
	}
	var chunks []string
	runes := []rune(reply)
	for i := 0; i < len(runes); {
		// 随机 1-3 个字符一组，模拟真实 token 节奏
		size := rand.Intn(3) + 1
		end := i + size
		if end > len(runes) {
			end = len(runes)
		}
		chunk := string(runes[i:end])
		_ = utf8.RuneCountInString(chunk) // 确保合法
		chunks = append(chunks, chunk)
		i = end
	}
	return chunks
}
