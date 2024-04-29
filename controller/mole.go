package controller

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/gordon104532/go_ws_game/enums"
)

func (h *HttpServer) getMoleScoreFromCache(ctx context.Context, username string) (score int, err error) {
	scoreFloat, err := h.cache.ZScore(ctx, enums.MoleLeaderBoardKey, username).Result()
	if err != nil {
		return 0, err
	}

	return int(scoreFloat), nil
}

func (h *HttpServer) setMoleScoreToCache(ctx context.Context, username string, score int) (rank int64, err error) {
	setMember := &redis.Z{
		Member: username,
		Score:  float64(score),
	}

	err = h.cache.ZAdd(ctx, enums.MoleLeaderBoardKey, setMember).Err()
	if err != nil {
		return 0, err
	}

	rank, err = h.cache.ZRevRank(ctx, enums.MoleLeaderBoardKey, username).Result()
	if err != nil {
		return 0, err
	}

	rank += 1
	return rank, nil
}

type leaderBoardStruct struct {
	Username string  `json:"username"`
	Score    float64 `json:"score"`
}

func (h *HttpServer) getMoleLeaderBoardFromCache(ctx context.Context) (leaderBoard []leaderBoardStruct, err error) {
	leaderBoard = make([]leaderBoardStruct, 0)
	data, err := h.cache.ZRevRange(ctx, enums.MoleLeaderBoardKey, 0, 13).Result()
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return leaderBoard, nil
	}

	// 批量獲取所有成員的分數
	scores, err := h.cache.ZMScore(ctx, enums.MoleLeaderBoardKey, data...).Result()
	if err != nil {
		return nil, err
	}

	if len(scores) != len(data) {
		return nil, fmt.Errorf("failed to get scores for all members")
	}

	for i, member := range data {
		score := scores[i]

		leaderBoard = append(leaderBoard, leaderBoardStruct{
			Username: member,
			Score:    score,
		})
	}

	return leaderBoard, nil
}
