package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/gordon104532/go_ws_game/enums"
)

type User struct {
	Name         string  `json:"name"`
	AnsweredList []int64 `json:"answered_list"`
}

func (h *HttpServer) GetUserDataFromCache(ctx context.Context, name string) (u *User, err error) {
	str, err := h.cache.HGet(ctx, enums.UserDataKey, name).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(str), &u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (h *HttpServer) setUserAnswered(ctx context.Context, quizId string, username string) error {
	err := h.cache.RPush(ctx, fmt.Sprintf(enums.UserAnsweredKey, username), quizId).Err()
	if err != nil {
		return err
	}

	return nil
}

func (h *HttpServer) getUserAnswered(ctx context.Context, username string) (answeredMap map[string]interface{}, err error) {
	answeredList, err := h.cache.LRange(ctx, fmt.Sprintf(enums.UserAnsweredKey, username), 0, -1).Result()
	if err != nil {
		return nil, err
	}

	answeredMap = make(map[string]interface{})
	for _, v := range answeredList {
		answeredMap[v] = v
	}

	return answeredMap, nil
}

func (h *HttpServer) incrEasterEggCount(ctx context.Context, username string) error {
	setMember := &redis.Z{
		Member: username,
		Score:  1, // 加一分
	}

	err := h.cache.ZIncr(ctx, enums.EasterEggCountKey, setMember).Err()
	if err != nil {
		return err
	}

	return nil
}
