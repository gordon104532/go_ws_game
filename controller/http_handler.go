package controller

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/secure"
)

type HttpServer struct {
	isHttps     bool
	httpPort    string
	cache       *redis.Client
	tgBotClient *tgBot.BotAPI
	tgToken     string
	tgChannel   int64
}

func NewHttpServer(
	httpPort string,
	cache *redis.Client,
) *HttpServer {
	cannel, err := strconv.Atoi(os.Getenv("TG_CHANNEL"))
	if err != nil {
		log.Fatal(err)
	}
	h := &HttpServer{
		isHttps:   true,
		httpPort:  httpPort,
		cache:     cache,
		tgToken:   os.Getenv("TG_BOT_TOKEN"),
		tgChannel: int64(cannel),
	}

	h.initTgBot()
	h.init()
	return h
}

func (h *HttpServer) Run() {
	log.Info("httpServer start")
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// allow origin
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))

	// 前端靜態檔
	router.Static("/button_with_sound", "../template")
	// 憑證challenge
	router.Static("/.well-known", "./.well-known")
	// SoundPad, AIPopeCat
	router.Static("/media", "../media")

	api := router.Group("/api")
	{
		// 通用
		api.GET("/health_check", healthCheck)
		api.GET("/easter_egg", h.checkEasterEgg)
		api.POST("/easter_egg", h.activateEasterEgg)
		api.POST("/report", healthCheck)

		// 題目
		api.GET("/quiz", h.getQuiz)
		api.POST("/quiz", h.answerQuiz)
		api.PUT("/quiz", h.setQuiz)
		api.GET("/quiz/init", h.initQuiz)    // 從快取中載入題目
		api.GET("/quiz/list", h.getQuizList) // 取得所有題目
		api.DELETE("/quiz", h.reportQuiz)    // 舉報題目

		// 打蘑菇
		api.POST("/mole", h.setMoleScore)
		api.GET("/mole", h.getMoleLeaderBoard) // 排行榜
	}

	// 除了有定義路由的 API 之外，其他都會到前端框架
	router.NoRoute(func(ctx *gin.Context) {
		ctx.File("../template/index.html")
	})

	if h.isHttps && os.Getenv("ENV") == "docker" {
		log.Info("httpServer TLS mode")
		router.Use(h.TlsHandler())

		err := router.RunTLS(":"+h.httpPort, "./.ssl/fullchain.pem", "./.ssl/privkey.pem")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		if err := router.Run(":" + h.httpPort); err != nil {
			log.Fatal(err)
		}
	}
}

func healthCheck(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "text/plain", []byte("health_check ok"))
}

func (h *HttpServer) TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     ":" + h.httpPort,
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}

func (h *HttpServer) initQuiz(ctx *gin.Context) {
	err := h.initQuizFromCache(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (h *HttpServer) setQuiz(ctx *gin.Context) {
	var quiz *Quiz
	err := ctx.ShouldBindJSON(&quiz)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	quiz.Id = strconv.FormatInt(time.Now().Unix(), 10)
	quiz.AnsweredBy = make([][]string, 2)
	quiz.AnsweredBy[0] = []string{}
	quiz.AnsweredBy[1] = []string{}

	err = h.setQuizToCache(quiz)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	quizList[quiz.Id] = quiz

	errStr := fmt.Sprintf("[quiz] setQuiz success, quiz: %+v", quiz.Question)
	go h.SendTgMsg(errStr)

	ctx.JSON(http.StatusOK, gin.H{"message": "setQuiz success"})
}

type AnswerPayload struct {
	Username string `json:"username" validator:"required"`
	QuizId   string `json:"quiz_id" validator:"required"`
	Answer   string `json:"answer"` // A or B
}

func (h *HttpServer) answerQuiz(ctx *gin.Context) {
	var answer AnswerPayload
	err := ctx.ShouldBindJSON(&answer)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	answer.Username = strings.TrimSpace(answer.Username)

	if _, ok := quizList[answer.QuizId]; !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "quiz not found"})
		return
	}

	err = h.setChoiceToCache(ctx, answer.QuizId, answer.Answer, answer.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	err = h.setUserAnswered(ctx, answer.QuizId, answer.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	if answer.Answer == "A" {
		quizList[answer.QuizId].AnsweredBy[0] = append(quizList[answer.QuizId].AnsweredBy[0], answer.Username)
	} else {
		quizList[answer.QuizId].AnsweredBy[1] = append(quizList[answer.QuizId].AnsweredBy[1], answer.Username)
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "answerQuiz success"})
}

func (h *HttpServer) getQuiz(ctx *gin.Context) {
	username := strings.TrimSpace(ctx.Query("username"))

	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "username empty"})
		return
	}

	answeredList, err := h.getUserAnswered(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	for _, q := range quizList {
		if _, ok := answeredList[q.Id]; !ok {
			ctx.JSON(http.StatusOK, q)
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "quiz done"})
}

func (h *HttpServer) getQuizList(ctx *gin.Context) {
	if len(quizList) < 1 {
		log.Error("getQuizList err, quizList: ", quizList)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "getQuizList len < 1"})
		return
	}

	ctx.JSON(http.StatusOK, quizList)
}

func (h *HttpServer) reportQuiz(ctx *gin.Context) {
	qid := strings.TrimSpace(ctx.Query("qid"))
	username := strings.TrimSpace(ctx.Query("username"))

	q, ok := quizList[qid]
	if !ok {
		log.Errorf("[quiz] reportQuiz not found qid: %s", qid)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "quiz not found"})
		return
	}

	errStr := fmt.Sprintf("[quiz] reportQuiz by: %s, quiz: %+v", username, q)
	go h.SendTgMsg(errStr)
	ctx.JSON(http.StatusOK, gin.H{"message": "setQuiz success"})
}

func (h *HttpServer) checkEasterEgg(ctx *gin.Context) {
	// 正式用 5/14 中午12點
	activeAt, err := time.Parse("2006-01-02 15:04:05", "2024-05-14 04:00:00")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	log.Info("checkEasterEgg", time.Now(), activeAt)
	if time.Now().Before(activeAt) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "EasterEgg Not Activated Yet", "start_at": activeAt.Unix()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "EasterEgg Activated"})
}

func (h *HttpServer) activateEasterEgg(ctx *gin.Context) {
	username := strings.TrimSpace(ctx.Query("username"))
	h.incrEasterEggCount(ctx, username)

	ctx.JSON(http.StatusOK, gin.H{"message": "activateEasterEgg incr success"})
}

func (h *HttpServer) initTgBot() {
	var err error
	h.tgBotClient, err = tgBot.NewBotAPI(h.tgToken)
	if err != nil {
		log.WithError(err).Error("initTgBot failed")
	}
	h.tgBotClient.Debug = false
}

// 發訊息給我
func (h *HttpServer) SendTgMsg(msg string) {
	NewMsg := tgBot.NewMessage(h.tgChannel, msg)
	NewMsg.ParseMode = tgBot.ModeHTML //傳送html格式的訊息
	_, err := h.tgBotClient.Send(NewMsg)
	if err != nil {
		log.WithError(err).Error("send telegram message error")
	}
}

func (h *HttpServer) setMoleScore(ctx *gin.Context) {
	username := strings.TrimSpace(ctx.Query("username"))
	score, err := strconv.Atoi(strings.TrimSpace(ctx.Query("score")))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "setMoleScore score not int"})
		return
	}

	if username == "" || score == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "param invalid"})
		return
	}

	currentScore, err := h.getMoleScoreFromCache(ctx, username)
	if err != nil && err != redis.Nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 現有分數較高則不更新
	if currentScore != 0 && currentScore > score {
		ctx.JSON(http.StatusOK, gin.H{"rank": 514})
		return
	}

	rank, err := h.setMoleScoreToCache(ctx, username, score)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"rank": rank})
}

func (h *HttpServer) getMoleLeaderBoard(ctx *gin.Context) {
	LeaderBoard, err := h.getMoleLeaderBoardFromCache(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"leader_board": LeaderBoard})
}
