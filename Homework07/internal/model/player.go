package model

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
