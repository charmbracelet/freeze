package config

// freeze/issues/50

type Config struct {
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

func Load() (*Config, error) {
	var c Config
	err := parseStruct(&c, "")
	if err != nil {
		return nil, err
	}
	return &c, nil
}
