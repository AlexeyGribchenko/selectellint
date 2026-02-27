package analyzer

type Config struct {
	CheckCapital   bool `mapstructure:"check-capital"`
	CheckInvalid   bool `mapstructure:"check-invalid"`
	CheckSensitive bool `mapstructure:"check-sensitive"`

	AllowedSymbols   string `mapstructure:"allowed-symbols"`
	ReplaceWithSpace string `mapstructure:"replace-with-space"`

	SensitiveWords []string `mapstructure:"sensitive-words"`
}

func DefaultConfig() Config {
	return Config{
		CheckCapital:     true,
		CheckInvalid:     true,
		CheckSensitive:   true,
		AllowedSymbols:   "",
		ReplaceWithSpace: "_-",
		SensitiveWords:   []string{"password", "token", "login", "email", "id", "api", "credential"},
	}
}
