package model_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mod.go/internal/model"
)

func TestParsePieceType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    model.PieceType
		wantErr bool
	}{
		{name: "пешка", input: "пешка", want: model.Pawn},
		{name: "Пешка", input: "Пешка", want: model.Pawn},
		{name: "ПЕШКА", input: "ПЕШКА", want: model.Pawn},
		{name: "  Пешка  ", input: "  Пешка  ", want: model.Pawn},
		{name: "пеШКа", input: "пеШКа", want: model.Pawn},

		{name: "ладья", input: "ладья", want: model.Rook},
		{name: "Л", input: "Л", want: model.Rook},
		{name: "r", input: "r", want: model.Rook},
		{name: "Rook", input: "Rook", want: model.Rook},

		{name: "КОНЬ", input: "КОНЬ", want: model.Knight},
		{name: "к", input: "к", want: model.Knight},
		{name: "knight", input: "knight", want: model.Knight},
		{name: "n", input: "n", want: model.Knight},

		{name: "слОн", input: "слОн", want: model.Bishop},
		{name: "с", input: "с", want: model.Bishop},
		{name: "BISHOP", input: "BISHOP", want: model.Bishop},
		{name: "B", input: "B", want: model.Bishop},

		{name: "ферзь", input: "ферзь", want: model.Queen},
		{name: "Ф", input: "Ф", want: model.Queen},
		{name: "q", input: "q", want: model.Queen},
		{name: "Queen", input: "Queen", want: model.Queen},

		{name: "Король", input: "Король", want: model.King},
		{name: "король", input: "король", want: model.King},
		{name: "КОРОЛЬ", input: "КОРОЛЬ", want: model.King},
		{name: "КоРоЛь", input: "КоРоЛь", want: model.King},

		// дефолтное значение для PieceType = Pawn.
		// Но мы выйдем раньше и никто его читать не будет
		{name: "Пустая строка", input: "", wantErr: true},
		{name: "Пробел", input: " ", wantErr: true},
		{name: "п - короткое для пешки нет", input: "п", wantErr: true},
		{name: "pawn - англ. для пешки нет", input: "pawn", wantErr: true},
		{name: "кор - короткое для короля нет", input: "кор", wantErr: true},
		{name: "king - англ. для короля нет", input: "king", wantErr: true},
		{name: "abc", input: "abc", wantErr: true},
		{name: "1234", input: "1234", wantErr: true},
		{name: "!@", input: "!@", wantErr: true},
		{name: "Queeen - опечатка, лишняя e", input: "Queeen", wantErr: true},
		{name: "кроль - опечатка, пропущено о", input: "кроль", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := model.ParsePieceType(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)

		})
	}
}

func TestPiece_MarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		piece *model.Piece
		want  string
	}{
		{
			name:  "белая пешка",
			piece: model.NewPiece(model.White, model.Pawn),
			want:  `{"цвет":"Белые","тип":"Пешка"}`,
		},
		{
			name:  "чёрный король",
			piece: model.NewPiece(model.Black, model.King),
			want:  `{"цвет":"Чёрные","тип":"Король"}`,
		},
		{
			name:  "белый ферзь",
			piece: model.NewPiece(model.White, model.Queen),
			want:  `{"цвет":"Белые","тип":"Ферзь"}`,
		},
		{
			name:  "чёрный конь",
			piece: model.NewPiece(model.Black, model.Knight),
			want:  `{"цвет":"Чёрные","тип":"Конь"}`,
		},
		{
			name:  "неизвестные тип и цвет",
			piece: model.NewPiece(model.Color(50), model.PieceType(100)),
			want:  `{"цвет":"Неизвестно","тип":"Неизвестно"}`,
		},
		{
			name:  "неизвестный тип, валидный цвет",
			piece: model.NewPiece(model.White, model.PieceType(77)),
			want:  `{"цвет":"Белые","тип":"Неизвестно"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.piece)
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestPiece_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantColor model.Color
		wantType  model.PieceType
		wantErr   bool
	}{
		{
			name:      "белая пешка",
			input:     `{"цвет":"Белые","тип":"Пешка"}`,
			wantColor: model.White,
			wantType:  model.Pawn,
		},
		{
			name:      "чёрный король",
			input:     `{"цвет":"Чёрные","тип":"Король"}`,
			wantColor: model.Black,
			wantType:  model.King,
		},
		{
			name:      "порядок ключей наоборот",
			input:     `{"тип":"Ладья","цвет":"Белые"}`,
			wantColor: model.White,
			wantType:  model.Rook,
		},
		{
			name:      "маленькие буквы",
			input:     `{"цвет":"белые","тип":"ферзь"}`,
			wantColor: model.White,
			wantType:  model.Queen,
		},
		{
			name:      "пустой",
			input:     `{}`,
			wantColor: model.White, // дефолтное значение
			wantType:  model.Pawn,  // дефолтное значение
		},
		{
			name:      "только цвет, без типа",
			input:     `{"цвет":"Чёрные"}`,
			wantColor: model.Black,
			wantType:  model.Pawn,
		},
		{
			name:      "только тип, без цвета",
			input:     `{"тип":"Ферзь"}`,
			wantColor: model.White,
			wantType:  model.Queen,
		},

		{
			name:    "сломанный JSON",
			input:   `{`,
			wantErr: true,
		},
		{
			name:    "неизвестный цвет",
			input:   `{"цвет":"Красные","тип":"Пешка"}`,
			wantErr: true,
		},
		{
			name:    "неизвестный тип",
			input:   `{"цвет":"Белые","тип":"Телепузик"}`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p model.Piece
			err := json.Unmarshal([]byte(tt.input), &p)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantColor, p.Color())
			assert.Equal(t, tt.wantType, p.Type())
		})
	}
}

func TestColor_MarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		color model.Color
		want  string
	}{
		{
			name:  "белые",
			color: model.White,
			want:  `"Белые"`,
		},
		{
			name:  "чёрные",
			color: model.Black,
			want:  `"Чёрные"`,
		},
		{
			name:  "неизвестный цвет",
			color: model.Color(100),
			want:  `"Неизвестно"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.color)
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestColor_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    model.Color
		wantErr bool
	}{
		{name: "Белые", input: `"Белые"`, want: model.White},
		{name: "белые", input: `"белые"`, want: model.White},

		{name: "Черные", input: `"Черные"`, want: model.Black},
		{name: "черные", input: `"черные"`, want: model.Black},
		{name: "Чёрные", input: `"Чёрные"`, want: model.Black},
		{name: "чёрные", input: `"чёрные"`, want: model.Black},

		{name: "БЕЛЫЕ (капс)", input: `"БЕЛЫЕ"`, want: model.White},
		{name: "ЧЁРНЫЕ (капс, с ё)", input: `"ЧЁРНЫЕ"`, want: model.Black},
		{name: "ЧЕРНЫЕ (капс, без ё)", input: `"ЧЕРНЫЕ"`, want: model.Black},

		{name: "неизвестный цвет", input: `"Красные"`, wantErr: true},
		{name: "пустая строка", input: `""`, wantErr: true},
		{name: "null", input: `null`, wantErr: true},
		{name: "число вместо строки", input: `3421`, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c model.Color
			err := json.Unmarshal([]byte(tt.input), &c)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, c)
		})
	}
}

func TestPieceType_MarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		piece model.PieceType
		want  string
	}{
		{name: "Пешка", piece: model.Pawn, want: `"Пешка"`},
		{name: "Ладья", piece: model.Rook, want: `"Ладья"`},
		{name: "Конь", piece: model.Knight, want: `"Конь"`},
		{name: "Слон", piece: model.Bishop, want: `"Слон"`},
		{name: "Ферзь", piece: model.Queen, want: `"Ферзь"`},
		{name: "Король", piece: model.King, want: `"Король"`},
		{name: "неизвестный тип", piece: model.PieceType(100), want: `"Неизвестно"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.piece)
			require.NoError(t, err)
			assert.JSONEq(t, tt.want, string(got))
		})
	}
}

func TestPieceType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    model.PieceType
		wantErr bool
	}{
		{name: "Пешка", input: `"Пешка"`, want: model.Pawn},
		{name: "ферзь", input: `"ферзь"`, want: model.Queen},
		{name: "КоролЬ", input: `"КоролЬ"`, want: model.King},
		{name: "ФЕРЗЬ капс", input: `"ФЕРЗЬ"`, want: model.Queen},
		{name: "q", input: `"q"`, want: model.Queen},
		// Т.к. пешка - вариант по умолчанию
		{name: "null", input: `null`, want: model.Pawn},
		{name: "Пустая строка", input: `""`, want: model.Pawn},

		{name: "Число 0", input: `0`, wantErr: true},
		{name: "Число 12", input: `12`, wantErr: true},
		{name: "Неизвестная фигура", input: `"Телепузик"`, wantErr: true},
		{name: "Строка null", input: `"null"`, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pt model.PieceType
			err := json.Unmarshal([]byte(tt.input), &pt)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, pt)
		})
	}
}

func TestPiece_Symbol(t *testing.T) {
	tests := []struct {
		name  string
		piece *model.Piece
		want  string
	}{
		{name: "Белая пешка", piece: model.NewPiece(model.White, model.Pawn), want: "♙"},
		{name: "Белая ладья", piece: model.NewPiece(model.White, model.Rook), want: "♖"},
		{name: "Белый конь", piece: model.NewPiece(model.White, model.Knight), want: "♘"},
		{name: "Белый слон", piece: model.NewPiece(model.White, model.Bishop), want: "♗"},
		{name: "Белый ферзь", piece: model.NewPiece(model.White, model.Queen), want: "♕"},
		{name: "Белый король", piece: model.NewPiece(model.White, model.King), want: "♔"},

		{name: "Чёрная пешка", piece: model.NewPiece(model.Black, model.Pawn), want: "♟"},
		{name: "Чёрная ладья", piece: model.NewPiece(model.Black, model.Rook), want: "♜"},
		{name: "Чёрный конь", piece: model.NewPiece(model.Black, model.Knight), want: "♞"},
		{name: "Чёрный слон", piece: model.NewPiece(model.Black, model.Bishop), want: "♝"},
		{name: "Чёрный ферзь", piece: model.NewPiece(model.Black, model.Queen), want: "♛"},
		{name: "Чёрный король", piece: model.NewPiece(model.Black, model.King), want: "♚"},

		{name: "неизвестный тип", piece: model.NewPiece(model.Black, model.PieceType(321)), want: " "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.piece.Symbol())
		})
	}
}
