package config

import "github.com/joho/godotenv"

func GetConfig() *Config {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	return &Config{
		App: App{
			Name:  "go-fiber-api",
			Port:  GetEnvStr("APP_PORT", "3000"),
			Debug: GetEnvBool("APP_DEBUG", true),
		},
		Database: Database{
			Host:     GetEnvStr("DB_HOST", "localhost"),
			Port:     GetEnvInt("DB_PORT", 5432),
			User:     GetEnvStr("DB_USER", "postgres"),
			Password: GetEnvStr("DB_PASSWORD", ""),
			Name:     GetEnvStr("DB_NAME", "postgres"),
		},
		Jwt: Jwt{
			Secret:           GetEnvStr("JWT_SECRET", "secret"),
			AccessExpireMin:  GetEnvInt("JWT_ACCESS_EXPIRE_MIN", 15),
			RefreshExpireMin: GetEnvInt("JWT_REFRESH_EXPIRE_MIN", 60*24*3),
		},
	}
}

type App struct {
	Name  string
	Port  string
	Debug bool
}

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type Jwt struct {
	Secret           string
	AccessExpireMin  int
	RefreshExpireMin int
}

type Config struct {
	App
	Database
	Jwt
}
