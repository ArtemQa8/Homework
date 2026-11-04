package auth_test

import (
	"testing"
	"time"

	"mod.go/internal/auth"
)

func TestGenerateToken(t *testing.T) {
	secret := []byte("test-secret-key-at-least-32-chars-long-xxxxx")

	tests := []struct {
		name  string
		login string
		ttl   time.Duration
	}{
		{
			name:  "обычный логин",
			login: "testuser",
			ttl:   time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := auth.GenerateToken(tt.login, secret, tt.ttl)
			// Маловероятные ошибки, страховка на всякий случай
			// (Кто-то осознано должен сломать метод внутри jwt)
			if err != nil {
				t.Fatalf("GenerateToken() вернул ошибку: %v", err)
			}
			// Буквально один сценарий - если кто-то добавит в код return "", nil
			// В остальных случаях должна исключаться верхней ошибкой
			if token == "" {
				t.Fatal("GenerateToken() вернул пустой токен")
			}
		})
	}

}

func TestValidateToken(t *testing.T) {

	secret := []byte("test-secret-key-at-least-32-chars-long-xxxxx")
	wrongSecret := []byte("another-secret-key-at-least-32-chars-long-yyyy")

	validToken, err := auth.GenerateToken("testuser", secret, time.Hour)
	if err != nil {
		t.Fatalf("не удалось подготовить validToken: %v", err)
	}
	expiredToken, err := auth.GenerateToken("testuser", secret, -time.Hour)
	if err != nil {
		t.Fatalf("не удалось подготовить expiredToken: %v", err)
	}

	tests := []struct {
		name        string
		tokenString string
		secret      []byte
		want        bool
		wantErr     bool
	}{
		{
			name:        "Валидный токен",
			tokenString: validToken,
			secret:      secret,
			want:        true,
			wantErr:     false,
		},
		{
			name:        "Неверный секрет",
			tokenString: validToken,
			secret:      wrongSecret,
			want:        false,
			wantErr:     true,
		},
		{
			name:        "Истёкший токен",
			tokenString: expiredToken,
			secret:      secret,
			want:        false,
			wantErr:     true,
		},
		{
			name:        "Мусор",
			tokenString: "это-не-токен",
			secret:      secret,
			want:        false,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.ValidateToken(tt.tokenString, tt.secret)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ValidateToken() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("ValidateToken() succeeded unexpectedly")
			}

			if got != tt.want {
				t.Errorf("ValidateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
