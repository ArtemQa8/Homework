package model

import "encoding/json"

type Move struct {
	FromRow    int
	FromCol    int
	ToRow      int
	ToCol      int
	Captured   *Piece
	Promotion  PieceType
	MovedPiece *Piece
	Check      bool
	Mate       bool
	gameID     int
	id         int
}

func NewMove(fromRow, fromCol, toRow, toCol int) *Move {
	return &Move{FromRow: fromRow, FromCol: fromCol, ToRow: toRow, ToCol: toCol}
}

func (m Move) ObjectType() string { return "ход" }
func (m Move) ID() int            { return m.id }
func (m *Move) SetID(id int)      { m.id = id }
func (m Move) GameID() int        { return m.gameID }
func (m *Move) SetGameID(id int)  { m.gameID = id }

// omitempty - Если наше поле будет "нулевым", то его не будет в джисоне
func (m Move) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		GameID     int       `json:"играID"`
		ID         int       `json:"id"`
		FromRow    int       `json:"отСтрока"`
		FromCol    int       `json:"отСтолбец"`
		ToRow      int       `json:"вСтрока"`
		ToCol      int       `json:"вСтолбец"`
		Captured   *Piece    `json:"съедена,omitempty"`
		Promotion  PieceType `json:"превращение,omitempty"`
		MovedPiece *Piece    `json:"ходившаяФигура,omitempty"`
		Check      bool      `json:"шах,omitempty"`
		Mate       bool      `json:"мат,omitempty"`
	}{
		GameID:     m.gameID,
		ID:         m.id,
		FromRow:    m.FromRow,
		FromCol:    m.FromCol,
		ToRow:      m.ToRow,
		ToCol:      m.ToCol,
		Captured:   m.Captured,
		Promotion:  m.Promotion,
		MovedPiece: m.MovedPiece,
		Check:      m.Check,
		Mate:       m.Mate,
	})
}

func (m *Move) UnmarshalJSON(data []byte) error {
	var payload struct {
		GameID     int       `json:"играID"`
		ID         int       `json:"id"`
		FromRow    int       `json:"отСтрока"`
		FromCol    int       `json:"отСтолбец"`
		ToRow      int       `json:"вСтрока"`
		ToCol      int       `json:"вСтолбец"`
		Captured   *Piece    `json:"съедена,omitempty"`
		Promotion  PieceType `json:"превращение,omitempty"`
		MovedPiece *Piece    `json:"ходившаяФигура,omitempty"`
		Check      bool      `json:"шах,omitempty"`
		Mate       bool      `json:"мат,omitempty"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	m.gameID = payload.GameID
	m.id = payload.ID
	m.FromRow = payload.FromRow
	m.FromCol = payload.FromCol
	m.ToRow = payload.ToRow
	m.ToCol = payload.ToCol
	m.Captured = payload.Captured
	m.Promotion = payload.Promotion
	m.MovedPiece = payload.MovedPiece
	m.Check = payload.Check
	m.Mate = payload.Mate
	return nil
}
