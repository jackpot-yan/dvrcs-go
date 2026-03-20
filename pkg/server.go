package server

import (
	logger "dcrcs-go/utils"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Server() *gin.Engine {
	go func() {
		rand.Seed(time.Now().UnixNano())
		http.HandleFunc("/event", CorsHandler(SseHandler))
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"status": "ok"}`)
		})
		logger.Info("SSE server listening on http://localhost:6000")
		logger.Info("Test: curl -N http://localhost:6000/event?prompt=你好")
		http.ListenAndServe(":6000", nil)
	}()
	router := gin.Default()
	baseHttp := router.Group("/base")
	baseHttp.POST("/input", Message)
	return router
}
