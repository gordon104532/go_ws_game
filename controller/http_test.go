package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var httpServer *HttpServer

func TestMain(m *testing.M) {
	cache := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	httpServer = NewHttpServer("8080", cache)
	go httpServer.Run()
	m.Run()
}

func TestGetQuiz(t *testing.T) {
	// Setup
	router := gin.Default()
	router.GET("/api/quiz", httpServer.getQuiz)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/quiz?username=test", nil)

	// Test
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "罐製預設題1")
}

func TestGetNotAnsweredQuiz(t *testing.T) {
	// Setup
	router := gin.Default()
	router.GET("/api/quiz", httpServer.getQuiz)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/quiz?username=罐製", nil)

	// Test
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "罐製預設題2")
}

func TestSetQuiz(t *testing.T) {
	// Setup
	router := gin.Default()
	router.POST("/api/quiz", httpServer.setQuiz)
	w := httptest.NewRecorder()
	quiz := Quiz{Question: "test_question?", Choices: []string{"A", "B"}, Author: "test_displayName"}
	jsonBytes, _ := json.Marshal(quiz)
	req, _ := http.NewRequest("POST", "/api/quiz", bytes.NewReader(jsonBytes))

	// Test
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "setQuiz success")
}

func TestAnswerQuizNotFound(t *testing.T) {
	// Setup
	router := gin.Default()
	router.POST("/api/quiz", httpServer.answerQuiz)
	w := httptest.NewRecorder()
	answer := AnswerPayload{Username: "test", QuizId: "1", Answer: "A"}
	jsonBytes, _ := json.Marshal(answer)
	req, _ := http.NewRequest("POST", "/api/quiz", bytes.NewReader(jsonBytes))

	// Test
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "quiz not found")
}

func TestAnswerQuiz(t *testing.T) {
	// Setup
	router := gin.Default()
	router.POST("/api/quiz", httpServer.answerQuiz)
	w := httptest.NewRecorder()
	answer := AnswerPayload{Username: "罐製", QuizId: "罐製預設題1", Answer: "A"}
	jsonBytes, _ := json.Marshal(answer)
	req, _ := http.NewRequest("POST", "/api/quiz", bytes.NewReader(jsonBytes))

	// Test
	router.ServeHTTP(w, req)

	// Assertions
	fmt.Println(w.Body.String())
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "answerQuiz success")
}
