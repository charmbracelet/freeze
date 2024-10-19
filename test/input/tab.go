package main

// freeze/issues/50

type Config struct { //nolint: revive
	Telegram struct {
		Token   string `env:"TG_TOKEN"`
		ChatID  string `env:"TG_CHAT"`
		OwnerID string `env:"TG_ADMIN"`
	}

	Database struct {
		DSN string `env:"DB_DSN"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL"`
	}

	Debug bool
}

func Load() (*Config, error) { //nolint: revive
	var c Config
	var err error
	return &c, err
}
