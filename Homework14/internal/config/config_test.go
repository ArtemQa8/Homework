package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mod.go/internal/config"
)

func TestConfig_Load(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		wantErr bool
	}{
		{
			name: "все переменные заданы",
			env: map[string]string{
				"LOGIN":      "user",
				"PASSWORD":   "pass",
				"JWT_SECRET": "supersecret",
				"JWT_TTL":    "1h",
			},
		},
		{
			name: "нет LOGIN",
			env: map[string]string{
				"LOGIN":      "",
				"PASSWORD":   "pass",
				"JWT_SECRET": "supersecret",
				"JWT_TTL":    "1h",
			},
			wantErr: true,
		},
		{
			name: "нет PASSWORD",
			env: map[string]string{
				"LOGIN":      "user",
				"PASSWORD":   "",
				"JWT_SECRET": "supersecret",
				"JWT_TTL":    "1h",
			},
			wantErr: true,
		},
		{
			name: "нет JWT_SECRET",
			env: map[string]string{
				"LOGIN":      "user",
				"PASSWORD":   "pass",
				"JWT_SECRET": "",
				"JWT_TTL":    "1h",
			},
			wantErr: true,
		},
		{
			name: "нет JWT_TTL",
			env: map[string]string{
				"LOGIN":      "user",
				"PASSWORD":   "pass",
				"JWT_SECRET": "supersecret",
				"JWT_TTL":    "",
			},
			wantErr: true,
		},
		{
			name: "невалидный JWT_TTL",
			env: map[string]string{
				"LOGIN":      "user",
				"PASSWORD":   "pass",
				"JWT_SECRET": "supersecret",
				"JWT_TTL":    "не_длительность",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfg, err := config.Load()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "user", cfg.Login)
			assert.Equal(t, "pass", cfg.Password)
			assert.Equal(t, []byte("supersecret"), cfg.JWTSecret)
			assert.Equal(t, 3600000000000, int(cfg.JWTTTL)) // 1h в наносекундах
		})
	}
}
