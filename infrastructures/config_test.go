package infrastructures

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Arrange: create a temp directory to place the fake config.env file
	tempDir := t.TempDir()

	configContent := `
APP_PORT=8080
BASE_URL=http://localhost

DB_HOST=localhost
DB_PORT=5432
DB_USER=testuser
DB_PASSWORD=testpass
DB_NAME=testdb

JWT_SECRET_KEY=secret
JWT_EXPIRATION_MINUTES=60
REFRESH_TOKEN_SECRET=refreshsecret
REFRESH_TOKEN_EXPIRATION_MINUTES=1440

SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=smtpuser
SMTP_PASSWORD=smtppass
SMTP_KEY=smtpkey
EMAIL_FROM=noreply@example.com
EMAIL_FROM_NAME=Test App

OPENAI_API_KEY=dummykey

DEFAULT_PAGE_SIZE=10
MAX_PAGE_SIZE=50

ALLOWED_ORIGINS=http://localhost,http://example.com
LOG_LEVEL=debug
TIMEZONE=Africa/Addis_Ababa
`

	// Write the config to a file
	err := os.WriteFile(filepath.Join(tempDir, "config.env"), []byte(configContent), 0644)
	assert.NoError(t, err)

	// Override viper config path to the temp directory
	viper.Reset()
	viper.AddConfigPath(tempDir)
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Act: call LoadConfig
	cfg, err := LoadConfig()

	// Assert: no error and all fields are correctly parsed
	assert.NoError(t, err)
	assert.Equal(t, "8080", cfg.AppPort)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, 60, cfg.JWTExpirationMinutes)
	assert.Equal(t, 1440, cfg.RefreshTokenExpirationMin)
	assert.Equal(t, []string{"http://localhost", "http://example.com"}, cfg.AllowedOrigins)
	assert.Equal(t, "Africa/Addis_Ababa", cfg.Timezone)
}
