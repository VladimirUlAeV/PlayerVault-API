package data

type player struct{
	ID int `json: "id"`
	Nickname string `json: "nickname"`
	Score int `json: "score"`
}

type PlayerStorage interface{

}

type MemoryStorage struct{
	Players map[int]player
}

func NewMemoryStorage()*MemoryStorage{
	return &MemoryStorage{
		Players: make(map[int]player),
	}
}

