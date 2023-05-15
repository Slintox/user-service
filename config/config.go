package config

import "github.com/ilyakaznacheev/cleanenv"

// Костыль
const (
	PostgresDev = true
)

type (
	Config struct {
		GRPC     *GRPCServerConfig
		Postgres *PostgresConfig
	}

	GRPCServerConfig struct {
		Port string `yaml:"grpc_port" env:"GRPC_PORT" env-default:":50052"`
	}

	PostgresConfig struct {
		DSN string `yaml:"postgres_dsn" env:"PG_DSN" env-default:"host=localhost port=54322 dbname=user user=user-user password=user-password sslmode=disable"`
	}
)

func InitConfig(configPath string) (*Config, error) {
	var err error
	cfg := Config{
		GRPC:     &GRPCServerConfig{},
		Postgres: &PostgresConfig{},
	}

	// Игнорируем файлы конфигураций, если путь к файлу не указан.
	// Если путь указан, но не валиден, возвращается ошибка
	if configPath != "" {
		err = cleanenv.ReadConfig(configPath, cfg.GRPC)
		if err != nil {
			return nil, err
		}

		if err = cleanenv.ReadConfig(configPath, cfg.Postgres); err != nil {
			return nil, err
		}
	}

	// Чтение конфигов из переменных окружения

	if err = cleanenv.ReadEnv(cfg.GRPC); err != nil {
		return nil, err
	}

	if err = cleanenv.ReadEnv(cfg.Postgres); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) GetPostgresConfig() *PostgresConfig {
	return c.Postgres
}

func (c *Config) GetGRPCConfig() *GRPCServerConfig {
	return c.GRPC
}
