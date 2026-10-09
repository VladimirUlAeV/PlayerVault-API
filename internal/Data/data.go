package data

type player struct{
	ID int `json: "id"`
	Nickname string `json: "nickname"`
	Score int `json: "score"`
}

type PlayerStorage interface{
	Get()
	GetAll()
	Save()
}

type MemoryStorage struct{
	Players map[int]player
}

func NewMemoryStorage()*MemoryStorage{
	return &MemoryStorage{
		Players: make(map[int]player),
	}
}

func (m * MemoryStorage)Get(id int)(error, int, string, int){
	player,ok := m.Players[id]
	if !ok {}

	return nil, id, player.Nickname, player.Score
}

func (m *MemoryStorage)Save(id int, nick string, score int){
	m.Players[id] = player{
		ID: id,
		Nickname: nick,
		Score: score,
	}
}

func (m *MemoryStorage)GetAll()map[int]player{
	return m.Players
}