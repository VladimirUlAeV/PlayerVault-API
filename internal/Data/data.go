package dataa

import (
	error_project "PlayerVaultAPI/internal/Error"
)

type PlayerStorage interface {
	Get(id int) (Player, error)
	GetAll() (map[int]Player, error)
	Save(structur Player) error
}
type Player struct {
	ID       int    `json: "id"`
	Nickname string `json: "nickname"`
	Score    int    `json: "score"`
}

type MemoryStorage struct {
	players map[int]Player
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		players: make(map[int]Player),
	}
}

func (m *MemoryStorage) Get(id int) (Player, error) {
	if player, ok := m.players[id]; ok == true {
		return player, nil
	} else {
		return Player{}, error_project.ErrorGetUser
	}
}

func (m *MemoryStorage) Save(structur Player) error {

	if _, ok := m.players[structur.ID]; ok == true {
		return error_project.ErrorIdValid
	}
	if structur.Nickname == "" {
		return error_project.ErrorNicknameEmpty
	}
	if structur.Score < 0 {
		return error_project.ErrorScoreInput
	}
	m.players[structur.ID] = Player{
		ID:       structur.ID,
		Nickname: structur.Nickname,
		Score:    structur.Score,
	}
	return nil
}

func (m *MemoryStorage) GetAll() (map[int]Player, error) {
	if l := len(m.players); l <= 0 {
		return make(map[int]Player), error_project.ErrorEmptyMap
	}
	return m.players, nil
}
