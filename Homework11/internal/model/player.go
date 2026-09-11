package model

import "encoding/json"

type Player struct {
	name  string
	color Color
	id    int
}

func NewPlayer(newName string, newColor Color) *Player {
	return &Player{name: newName, color: newColor}
}

func (p Player) Name() string          { return p.name }
func (p Player) Color() Color          { return p.color }
func (p Player) ID() int               { return p.id }
func (p *Player) SetID(id int)         { p.id = id }
func (p *Player) SetColor(color Color) { p.color = color }
func (p Player) ObjectType() string    { return "игрок" }

// MarshalJSON реализует сериализацию в JSON с сохранением приватных полей.
func (p Player) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID    int    `json:"id"`
		Name  string `json:"имя"`
		Color Color  `json:"цвет"`
	}{
		ID:    p.id,
		Name:  p.name,
		Color: p.color,
	})
}

// UnmarshalJSON реализует десериализацию из JSON с заполнением приватных полей.
func (p *Player) UnmarshalJSON(data []byte) error {
	var payload struct {
		ID    int    `json:"id"`
		Name  string `json:"имя"`
		Color Color  `json:"цвет"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	p.id = payload.ID
	p.name = payload.Name
	p.color = payload.Color
	return nil
}
