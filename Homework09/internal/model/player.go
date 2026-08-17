package model

import "encoding/json"

type Игрок struct {
	имя  string
	цвет ЦветФигуры
}

func НовыйИгрок(новоеИмя string, новыйЦвет ЦветФигуры) *Игрок {
	return &Игрок{имя: новоеИмя, цвет: новыйЦвет}
}

func (и Игрок) Имя() string        { return и.имя }
func (и Игрок) Цвет() ЦветФигуры   { return и.цвет }
func (и Игрок) ТипОбъекта() string { return "игрок" }

// MarshalJSON реализует сериализацию в JSON с сохранением приватных полей.
func (и Игрок) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Имя  string     `json:"имя"`
		Цвет ЦветФигуры `json:"цвет"`
	}{
		Имя:  и.имя,
		Цвет: и.цвет,
	})
}

// UnmarshalJSON реализует десериализацию из JSON с заполнением приватных полей.
func (и *Игрок) UnmarshalJSON(data []byte) error {
	var данные struct {
		Имя  string     `json:"имя"`
		Цвет ЦветФигуры `json:"цвет"`
	}
	if err := json.Unmarshal(data, &данные); err != nil {
		return err
	}
	и.имя = данные.Имя
	и.цвет = данные.Цвет
	return nil
}
