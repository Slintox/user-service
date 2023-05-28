package config

import "github.com/ilyakaznacheev/cleanenv"

// Костыль
const (
	PostgresDev = true
)

type (
	Config struct {
		GRPC     *GRPCServerConfig
		HTTP     *HTTPServerConfig
		Postgres *PostgresConfig
		Swagger  *SwaggerConfig
	}

	GRPCServerConfig struct {
		Port string `yaml:"grpc_port" env:"GRPC_PORT" env-default:":50052"`
	}

	HTTPServerConfig struct {
		Port string `yaml:"http_port" env:"HTTP_PORT" env-default:":8080"`
	}

	PostgresConfig struct {
		DSN string `yaml:"postgres_dsn" env:"PG_DSN" env-default:"host=localhost port=54322 dbname=user user=user-user password=user-password sslmode=disable"`
	}

	SwaggerConfig struct {
		Port string `yaml:"swagger_port" env:"SWAGGER_PORT" env-default:":8081"`
	}
)

func InitConfig(configPath string) (*Config, error) {
	var err error
	cfg := Config{
		GRPC:     &GRPCServerConfig{},
		HTTP:     &HTTPServerConfig{},
		Postgres: &PostgresConfig{},
		Swagger:  &SwaggerConfig{},
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

	if err = cleanenv.ReadEnv(cfg.HTTP); err != nil {
		return nil, err
	}

	if err = cleanenv.ReadEnv(cfg.Postgres); err != nil {
		return nil, err
	}

	if err = cleanenv.ReadEnv(cfg.Swagger); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) GetPostgresConfig() *PostgresConfig {
	return c.Postgres
}

func (c *Config) GetHTTPConfig() *HTTPServerConfig {
	return c.HTTP
}

func (c *Config) GetGRPCConfig() *GRPCServerConfig {
	return c.GRPC
}

func (c *Config) GetSwaggerConfig() *SwaggerConfig {
	return c.Swagger
}
