package config

type Config struct {
	HTTPHost         string   `envconfig:"HTTP_HOST" default:"localhost"`
	HTTPPort         string   `envconfig:"HTTP_PORT" default:"8080"`
	HTTPReadTimeout  int      `envconfig:"HTTP_READ_TIMEOUT" default:"10"`
	HTTPWriteTimeout int      `envconfig:"HTTP_WRITE_TIMEOUT" default:"10"`
	LogLevel         string   `envconfig:"LOG_LEVEL" default:"debug"`
	Postgres         Postgres `envconfig:"POSTGRES"`
	Redis            Redis    `envconfig:"REDIS"`
}

type Postgres struct {
	Host     string `envconfig:"HOST" default:"localhost"`
	Port     string `envconfig:"PORT" default:"5432"`
	Username string `envconfig:"USER" default:"postgres"`
	Password string `envconfig:"PASSWORD" default:"postgres"`
	Database string `envconfig:"DATABASE" default:"marketplace"`
}

type Redis struct {
	Addr     string `envconfig:"ADDR" default:"localhost:6379"`
	Password string `envconfig:"PASSWORD" default:""`
	DB       int    `envconfig:"DB" default:"0"`
}
