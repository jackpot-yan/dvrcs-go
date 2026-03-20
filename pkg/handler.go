package server

import (
	"dcrcs-go/agent"
	logger "dcrcs-go/utils"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UserInput struct {
	Msg     string   `json:"message"`
	Tool    []string `json:"tool"`
	AnyElse any      `json:"anyelse"`
}

func Message(c *gin.Context) {}

func CorsHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func SseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	flusher, ok := w.(http.Flusher)
	if !ok {
		logger.Error(w, "客户端创建失败: %s", http.StatusInternalServerError)
		return
	}
	prompt := r.URL.Query().Get("prompt")
	if prompt == "" {
		prompt = r.FormValue("prompt")
	}
	if prompt == "" {
		prompt = "你好"
	}
	logger.Infof("[SSE] prompt=%q remote=%s", prompt, r.RemoteAddr)
	ctx := r.Context()
	flusher.Flush()
	chunks := agent.SimulateAI(prompt)
	for _, chunk := range chunks {
		select {
		case <-ctx.Done():
			logger.Info("[SSE] client disconnected mid-stream")
			return
		default:
		}
		escaped := strings.ReplaceAll(chunk, "\n", "\\n")
		fmt.Fprintf(w, "event: delta\ndata:%s\n\n", escaped)
		flusher.Flush()
		delay := time.Duration(20+rand.Intn(60)) * time.Millisecond
		select {
		case <-ctx.Done():
			return
		case <-time.After(delay):
		}
	}
	fmt.Fprintf(w, "event:done\ndata: [DONE]\n\n")
	flusher.Flush()
	logger.Infof("[SSE] stream finished, prompt=%q", prompt)
}
