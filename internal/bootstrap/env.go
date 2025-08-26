package bootstrap

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Env holds all configuration; values are populated from real environment vars.
type Env struct {
	Port         string
	AppEnv       string
	DB_Uri       string
	DB_Name      string
	RTS          string
	ATS          string
	RefTEHours   int
	AccTEMinutes int
	CtxTSeconds  int

	Page               int
	PageSize           int
	Recency            string

	// user / auth collections
	UserCollection          string
	RefreshTokenCollection  string
	PasswordResetCollection string
	PasswordResetExpiry     int // minutes

	// email configuration
	SMTPHost     string
	SMTPPort     int
	SMTPFrom     string
	SMTPUsername string
	SMTPPassword string
	ResetURL     string

	// Gemini AI
	GeminiAPIKey    string
	GeminiModelName string

	// ImageKit
	ImageKitPrivateKey string
	ImageKitPublicKey  string
	ImageKitEndpoint   string

	// OTP
	SecretSalt         string
	OtpCollection      string
	OtpExpireMinutes   int
	OtpMaximumAttempts int

	// Redis
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int

	// Cache
	CacheExpirationSeconds int

	// Google OAuth2
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func NewEnv() (*Env, error) {
	_ = godotenv.Load() 

	env := &Env{}
	var err error

	env.Port = getEnv("PORT", ":8080")
	env.AppEnv = getEnv("APP_ENV", "development")
	env.DB_Uri = getEnv("DB_URI", "")
	env.DB_Name = getEnv("DB_NAME", "")
	env.RTS = getEnv("REFRESH_TOKEN_SECRET", "")
	env.ATS = getEnv("ACCESS_TOKEN_SECRET", "")
	env.RefTEHours, err = getEnvAsInt("REFRESH_TOKEN_EXPIRE_HOURS", 24*7)
	if err != nil { return nil, err }
	env.AccTEMinutes, err = getEnvAsInt("ACCESS_TOKEN_EXPIRE_MINUTES", 15)
	if err != nil { return nil, err }
	env.CtxTSeconds, err = getEnvAsInt("CONTEXT_TIMEOUT_SECONDS", 10)
	if err != nil { return nil, err }

	env.Page, err = getEnvAsInt("PAGE", 1); if err != nil { return nil, err }
	env.PageSize, err = getEnvAsInt("PAGE_SIZE", 10); if err != nil { return nil, err }
	env.Recency = getEnv("RECENCY", "new")
	
	env.UserCollection = getEnv("USER_COLLECTION", "users")
	env.RefreshTokenCollection = getEnv("REFRESH_TOKEN_COLLECTION", "refresh_tokens")
	env.PasswordResetCollection = getEnv("PASSWORD_RESET_TOKEN_COLLECTION", "password_resets")
	env.PasswordResetExpiry, err = getEnvAsInt("PASSWORD_RESET_TOKEN_EXPIRE_MINUTES", 15); if err != nil { return nil, err }

	env.SMTPHost = getEnv("SMTP_HOST", "smtp.gmail.com")
	env.SMTPPort, err = getEnvAsInt("SMTP_PORT", 587); if err != nil { return nil, err }
	env.SMTPFrom = getEnv("SMTP_FROM", "")
	env.SMTPUsername = getEnv("SMTP_USERNAME", "")
	env.SMTPPassword = getEnv("SMTP_PASSWORD", "")
	env.ResetURL = getEnv("RESET_URL", "http://localhost:3000/reset-password")

	env.GeminiAPIKey = getEnv("GEMINI_API_KEY", "")
	env.GeminiModelName = getEnv("GEMINI_MODEL_NAME", "gemini-2.0-flash")

	env.ImageKitPrivateKey = getEnv("IMAGEKIT_PRIVATE_KEY", "")
	env.ImageKitPublicKey = getEnv("IMAGEKIT_PUBLIC_KEY", "")
	env.ImageKitEndpoint = getEnv("IMAGEKIT_URL_ENDPOINT", "")

	env.SecretSalt = getEnv("MY_SUPER_SECRET_SALT", "")
	env.OtpCollection = getEnv("OTP_COLLECTION", "otps")
	env.OtpExpireMinutes, err = getEnvAsInt("OTP_EXPIRE_MINUTES", 5); if err != nil { return nil, err }
	env.OtpMaximumAttempts, err = getEnvAsInt("OTP_MAXIMUM_ATTEMPTS", 3); if err != nil { return nil, err }

	env.RedisHost = getEnv("REDIS_HOST", "localhost")
	env.RedisPort, err = getEnvAsInt("REDIS_PORT", 6379); if err != nil { return nil, err }
	env.RedisPassword = getEnv("REDIS_PASSWORD", "")
	env.RedisDB, err = getEnvAsInt("REDIS_DB", 0); if err != nil { return nil, err }

	env.CacheExpirationSeconds, err = getEnvAsInt("CACHE_EXPIRATION_SECONDS", 3600); if err != nil { return nil, err }

	env.GoogleClientID = getEnv("GOOGLE_CLIENT_ID", "")
	env.GoogleClientSecret = getEnv("GOOGLE_CLIENT_SECRET", "")
	env.GoogleRedirectURL = getEnv("GOOGLE_REDIRECT_URL", "")

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	if env.DB_Uri == "" {
		return nil, errors.New("DB_URI is required but empty; set it in .env or environment")
	}

	return env, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) (int, error) {
	if v, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(v)
		if err != nil { return 0, err }
		return i, nil
	}
	return fallback, nil
}