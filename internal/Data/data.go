package data

import error_project "PlayerVaultAPI/internal/Error"

type PlayerStorage interface {
	Get()
	GetAll()
	Save()
}
type player struct {
	ID       int    `json: "id"`
	Nickname string `json: "nickname"`
	Score    int    `json: "score"`
}
type memoryStorage struct {
	Players map[int]player
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		Players: make(map[int]player),
	}
}

func (m *memoryStorage) get(id int) (int, string, int) {
	player, ok := m.Players[id]
	if !ok {
		error_project.ErrorWriteLog(error_project.ErrorGetUser)
	}

	return id, player.Nickname, player.Score
}

func (m *memoryStorage) save(id int, nick string, score int) {
	m.Players[id] = player{
		ID:       id,
		Nickname: nick,
		Score:    score,
	}
}

func (m *memoryStorage) getAll() map[int]player {
	return m.Players
}