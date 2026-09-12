package model

import "encoding/json"

type Player struct {
	name string
	id   int
}

func NewPlayer(newName string) *Player {
	return &Player{name: newName}
}

func (p Player) Name() string       { return p.name }
func (p Player) ID() int            { return p.id }
func (p *Player) SetID(id int)      { p.id = id }
func (p Player) ObjectType() string { return "игрок" }

func (p Player) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID   int    `json:"id"`
		Name string `json:"имя"`
	}{
		ID:   p.id,
		Name: p.name,
	})
}

func (p *Player) UnmarshalJSON(data []byte) error {
	var payload struct {
		ID   int    `json:"id"`
		Name string `json:"имя"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	p.id = payload.ID
	p.name = payload.Name
	return nil
}
