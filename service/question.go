package service

type quizService struct {
	userMap map[string]*user
}

type user struct {
	Name         string
	AnsweredList []int64
}

type quiz struct {
	Describe string
}

var NewquizService = func() *quizService {
	return &quizService{}
}

func (q *quizService) Login(username string) error {
	if _, ok := q.userMap[username]; !ok {
		q.userMap[username] = &user{
			Name: username,
		}
	}

	return nil
}
