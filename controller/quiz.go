package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gordon104532/go_ws_game/enums"
	log "github.com/sirupsen/logrus"
)

func (h *HttpServer) init() {
	quizList = make(map[string]*Quiz)
	h.initQuizFromCache(context.Background())
}

type Quiz struct {
	Id         string     `json:"id"`
	Question   string     `json:"question" validator:"required"`
	Choices    []string   `json:"choices" validator:"required"`
	AnsweredBy [][]string `json:"answered_by"`
	Author     string     `json:"author" validator:"required"`
}

// hSet h:quiz_data "罐製_2" "{\"id\":\"罐製_2\",\"question\":\"ㄅ主播的車最近哪裡壞了?\",\"choices\":[\"變速箱\",\"車燈\"],\"author\":\"罐製\"}"
var quizList map[string]*Quiz

func (h *HttpServer) initQuizFromCache(ctx context.Context) error {
	quizList = map[string]*Quiz{
		"預設題0": {
			Id:         "預設題0",
			Question:   "你自認為你是?",
			Choices:    []string{"波寶", "波黑"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題1": {
			Id:         "預設題1",
			Question:   "波普貓愛吃魚的 '貓' 字代表的動物是?",
			Choices:    []string{"貓咪", "貓頭鷹"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題2": {
			Id:         "預設題2",
			Question:   "波普貓愛吃魚的英文是?",
			Choices:    []string{"popecatengblue", "pepecatengblue"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題3": {
			Id:         "預設題3",
			Question:   "波普貓淚痣是在哪邊?",
			Choices:    []string{"左邊", "右邊"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題4": {
			Id:         "預設題4",
			Question:   "波普貓比較早的soundcloud帳號是哪一個?",
			Choices:    []string{"engblue (EngBlue波普貓在英格藍)", "pope-cat (波普貓PopeCat)"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題5": {
			Id:         "預設題5",
			Question:   "波普貓精華YT頻道是?",
			Choices:    []string{"@PopecatLive (波普貓 Live)", "@popecatengblue (PopeCat波普貓)"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題6": {
			Id:         "預設題6",
			Question:   "波普貓DC群組名稱是?",
			Choices:    []string{"波普貓的第一個家", "波普貓的小城堡"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題7": {
			Id:         "預設題7",
			Question:   "大家常說的 'ㄅ一下', 'ㄅ了' 指的是?",
			Choices:    []string{"斷片", "依上下文情境都適用"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題8": {
			Id:         "預設題8",
			Question:   "波普貓所屬的團體是?",
			Choices:    []string{"貓窩娛樂", "狗窩娛樂"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題9": {
			Id:         "預設題9",
			Question:   "波普貓生日是?",
			Choices:    []string{"5/14", "整年都是"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題10": {
			Id:         "預設題10",
			Question:   "狗窩娛樂一月主題中，波普貓抱著的抱枕是?",
			Choices:    []string{"貓咪", "貓頭鷹"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題11": {
			Id:         "預設題11",
			Question:   "波普貓生日三周目，跟一周目差了幾歲?",
			Choices:    []string{"2 歲", "0 歲"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題12": {
			Id:         "預設題12",
			Question:   "波普貓最近一次去的國家是?",
			Choices:    []string{"西班牙", "日本"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題13": {
			Id:         "預設題13",
			Question:   "小奇點徽章要累計幾點才可以獲得吹風機?",
			Choices:    []string{"10,000點", "1,000點"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題14": {
			Id:         "預設題14",
			Question:   "十萬小奇點徽章是?",
			Choices:    []string{"早安美式鬆餅", "一定要吃到花好月圓"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題15": {
			Id:         "預設題15",
			Question:   "波普貓壽司飯友上面的是?",
			Choices:    []string{"鮭魚", "鮪魚"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題16": {
			Id:         "預設題16",
			Question:   "波普貓是狗窩娛樂的?",
			Choices:    []string{"睡覺擔當", "可愛擔當"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題17": {
			Id:         "預設題17",
			Question:   "哪裡比較適合開會?",
			Choices:    []string{"車上", "辦公室"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題18": {
			Id:         "預設題18",
			Question:   "喜歡的顏色是藍色所以iPhone的顏色是?",
			Choices:    []string{"紫色", "藍色"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題19": {
			Id:         "預設題19",
			Question:   "波普貓會在哪隻手指上戴戒指?",
			Choices:    []string{"右手食指、左手食指、左手小拇指", "右手食指、右手中指、左手食指"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題20": {
			Id:         "預設題20",
			Question:   "人稱'日麻夢想摧毀者'是因為槓了甚麼牌而一戰成名?",
			Choices:    []string{"九萬", "九筒"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題21": {
			Id:         "預設題21",
			Question:   "熱門精華'迷之大抖動-愛惹'，當時ㄅ主播的台詞是?",
			Choices:    []string{"22222", "RRRRR"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題22": {
			Id:         "預設題22",
			Question:   "以下兩個食物怎麼選?",
			Choices:    []string{"炸起司豬排", "火鍋"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題23": {
			Id:         "預設題23",
			Question:   "ㄅ主播追不到聊天室的時候，我們可以把這個現象稱為?",
			Choices:    []string{"在看VOD", "時空旅人"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題24": {
			Id:         "預設題24",
			Question:   "辦公室的甚麼可以帶來小雀幸?",
			Choices:    []string{"兩個零食櫃", "午休睡超過時間"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題25": {
			Id:         "預設題25",
			Question:   "二周年的驚喜翻唱歌名是?",
			Choices:    []string{"Unicorn", "油膩控"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題26": {
			Id:         "預設題26",
			Question:   "2023生日發布的Cover曲是?",
			Choices:    []string{"Young and Beautiful", "Summertime Sadness"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題27": {
			Id:         "預設題27",
			Question:   "在波普貓人設中是永遠的17歲，在甚麼特殊的情況下會變18歲?",
			Choices:    []string{"遇上要開車的時候", "遇上要喝酒的時候"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題28": {
			Id:         "預設題28",
			Question:   "在波普貓人設中，指甲油配色是哪隻手指頭不同色?",
			Choices:    []string{"無名指", "食指"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題29": {
			Id:         "預設題29",
			Question:   "在波普貓人設中，代表色紫藍色的色碼是?",
			Choices:    []string{"#8793BD", "#2a5dad"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題30": {
			Id:         "預設題30",
			Question:   "波普貓於二周年掛軸上哼出來的音符是?",
			Choices:    []string{"四分音符", "八分音符"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題31": {
			Id:         "預設題31",
			Question:   "波普貓二周年掛軸上一共有幾隻貓(含波普貓本貓)?",
			Choices:    []string{"11隻", "17隻"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題32": {
			Id:         "預設題32",
			Question:   "2023波普貓聖誕明信片的造型是?",
			Choices:    []string{"唱詩班", "天使"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
		"預設題33": {
			Id:         "預設題33",
			Question:   "波普貓曾喜歡哪位華語男歌手，並在他爆出私生活爭議後開台唱他的歌?",
			Choices:    []string{"潘瑋柏", "王力宏"},
			Author:     "罐製",
			AnsweredBy: make([][]string, 2),
		},
	}

	rawData := h.cache.HGetAll(ctx, enums.QuizDataKey).Val()

	if len(rawData) == 0 {
		log.Errorf("QuizData cache %s is empty", enums.QuizDataKey)
		return fmt.Errorf("QuizData cache %s is empty", enums.QuizDataKey)
	}

	for _, v := range rawData {
		var quiz Quiz
		if err := json.Unmarshal([]byte(v), &quiz); err != nil {
			log.Error(err)
			continue
		}
		quizList[quiz.Id] = &quiz
	}

	for _, q := range quizList {
		q.AnsweredBy = make([][]string, len(q.Choices))

		AnsweredAList := h.getAnsweredUser(ctx, q.Id, "A")
		if len(AnsweredAList) > 0 {
			q.AnsweredBy[0] = append(q.AnsweredBy[0], AnsweredAList...)
		} else {
			q.AnsweredBy[0] = []string{}
		}

		AnsweredBList := h.getAnsweredUser(ctx, q.Id, "B")
		if len(AnsweredBList) > 0 {
			q.AnsweredBy[1] = append(q.AnsweredBy[1], AnsweredBList...)
		} else {
			q.AnsweredBy[1] = []string{}
		}
	}

	log.Infof("[quiz] init from cache success, quiz len: %d", len(quizList))
	return nil
}

func (h *HttpServer) setQuizToCache(q *Quiz) error {
	qStr, err := json.Marshal(q)
	if err != nil {
		log.WithError(err).Error("[quiz] setQuizToCache Marshal error")
		return err
	}

	err = h.cache.HSet(context.Background(), enums.QuizDataKey, q.Id, qStr).Err()
	return err
}

func (h *HttpServer) setChoiceToCache(ctx context.Context, qId string, choice string, username string) error {
	err := h.cache.RPush(ctx, fmt.Sprintf(enums.QuizAnsweredKey, qId, choice), username).Err()
	if err != nil {
		log.WithError(err).Error("[quiz] setAnsweredToCache RPush error")
		return err
	}

	return nil
}

var showBlackList []string = []string{
	"測試",
	"TEST",
}

func (h *HttpServer) getAnsweredUser(ctx context.Context, qId string, choice string) []string {
	answeredList, err := h.cache.LRange(ctx, fmt.Sprintf(enums.QuizAnsweredKey, qId, choice), 0, -1).Result()
	if err != nil {
		log.WithError(err).Error("[quiz] getChoiceToCache RPush error")
		return nil
	}

	if answeredList == nil {
		return []string{}
	}

	var result []string
	for _, v := range answeredList {
		isBlack := false
		for _, black := range showBlackList {
			if strings.Contains(v, black) {
				isBlack = true
			}
		}

		if isBlack {
			continue
		}

		result = append(result, v)
	}

	return result
}
