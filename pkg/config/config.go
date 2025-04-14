package config

import (
	"financial-backend/internal/entities"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config representa as configurações da aplicação
type Config struct {
	ServerAddress  string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DefaultDueDate int
}

// LoadConfig carrega as configurações do ambiente
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("erro ao carregar arquivo .env: %v", err)
	}

	defaultDueDate, _ := strconv.Atoi(getEnv("DEFAULT_DUE_DATE", "15"))

	config := &Config{
		ServerAddress:  getEnv("SERVER_ADDRESS", ":8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "financial"),
		DefaultDueDate: defaultDueDate,
	}

	return config, nil
}

// SetupDatabase configura a conexão com o banco de dados
func SetupDatabase(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,          // Don't include params in the SQL log
			Colorful:                  true,         // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %v", err)
	}
	db.AutoMigrate(&entities.Budget{}, &entities.Expense{}, &entities.Income{})
	return db, nil
}

// getEnv retorna o valor de uma variável de ambiente ou um valor padrão
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
