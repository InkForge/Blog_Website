package infrastructures

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv  string
	AppPort string
	BaseURL string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBDriver   string
	DBUri      string

	JWTSecretKey              string
	JWTExpirationMinutes      int
	RefreshTokenSecret        string
	AccessTokenSecret         string
	RefreshTokenExpirationMin int

	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string
	SMTPKey       string
	EmailFrom     string
	EmailFromName string

	OpenAIApiKey string

	DefaultPageSize int
	MaxPageSize     int

	AllowedOrigins []string
	LogLevel       string
	Timezone       string

	GoogleClientID         string
	GoogleClientSecret     string
	GoogleRedirectURL      string
	GithubClientID         string
	GithubClientSecret     string
	GithubRedirectURL      string
	FacebookClientID       string
	FacebookClientSecret   string
	FacebookRedirectURL    string

}

// LoadConfig loads config.env file using absolute project path which looks for
// root/config.env. It returns config reference and nil if loading is a success.
// Else it returns nil and error message.
func LoadConfig() (*Config, error) {
	// Use runtime.Caller to get to root directory
	// it uses absolute path to the project
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "..")

	viper.AddConfigPath(projectRoot)
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{
		AppPort: viper.GetString("APP_PORT"),
		BaseURL: viper.GetString("BASE_URL"),

		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		DBUri:      viper.GetString("DB_URI"),

		JWTSecretKey:              viper.GetString("JWT_SECRET_KEY"),
		JWTExpirationMinutes:      viper.GetInt("JWT_EXPIRATION_MINUTES"),
		RefreshTokenSecret:        viper.GetString("REFRESH_TOKEN_SECRET"),
		RefreshTokenExpirationMin: viper.GetInt("REFRESH_TOKEN_EXPIRATION_MINUTES"),

		SMTPHost:      viper.GetString("SMTP_HOST"),
		SMTPPort:      viper.GetInt("SMTP_PORT"),
		SMTPUsername:  viper.GetString("SMTP_USERNAME"),
		SMTPPassword:  viper.GetString("SMTP_PASSWORD"),
		SMTPKey:       viper.GetString("SMTP_KEY"),
		EmailFrom:     viper.GetString("EMAIL_FROM"),
		EmailFromName: viper.GetString("EMAIL_FROM_NAME"),

		OpenAIApiKey: viper.GetString("OPENAI_API_KEY"),

		DefaultPageSize: viper.GetInt("DEFAULT_PAGE_SIZE"),
		MaxPageSize:     viper.GetInt("MAX_PAGE_SIZE"),

		AllowedOrigins: strings.Split(viper.GetString("ALLOWED_ORIGINS"), ","),
		LogLevel:       viper.GetString("LOG_LEVEL"),
		Timezone:       viper.GetString("TIMEZONE"),

		GoogleClientID:       viper.GetString("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:   viper.GetString("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:    viper.GetString("GOOGLE_REDIRECT_URL"),
		GithubClientID:       viper.GetString("GITHUB_CLIENT_ID"),
		GithubClientSecret:   viper.GetString("GITHUB_CLIENT_SECRET"),
		GithubRedirectURL:    viper.GetString("GITHUB_REDIRECT_URL"),
		FacebookClientID:     viper.GetString("FACEBOOK_CLIENT_ID"),
		FacebookClientSecret: viper.GetString("FACEBOOK_CLIENT_SECRET"),
		FacebookRedirectURL:  viper.GetString("FACEBOOK_REDIRECT_URL"),
	}

	return cfg, nil
}
