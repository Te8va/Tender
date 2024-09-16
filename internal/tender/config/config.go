package config

type Config struct {
	ServerAddress    string `env:"SERVER_ADDRESS" envDefault:"0.0.0.0:8080"`
	ServicePort      int    `env:"SERVICE_PORT"          envDefault:"8080"`
	ServiceHost      string `env:"SERVICE_HOST"          envDefault:"0.0.0.0"`
	PostgresUsername string `env:"POSTGRES_USERNAME"     envDefault:"tender"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"     envDefault:"tender"`
	PostgresDB       string `env:"POSTGRES_DATABASE"           envDefault:"tender"`
	PostgresPort     int    `env:"POSTGRES_PORT"         envDefault:"5432"`
	MigrationsPath   string `env:"MIGRATIONS_PATH"       envDefault:"migrations"`
	LogFilePath      string `env:"LOG_FILE_PATH"         envDefault:"logfile.log"`
	PostgresConn     string `env:"POSTGRES_CONN"         envDefault:"postgres://tender:tender@postgres:5432/tender"`
}
