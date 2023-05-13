package model

type Config struct {
	GptHost              string   `yaml:"GptHost"`
	GptProxy             string   `yaml:"GptProxy"`
	GptModel             string   `yaml:"GptModel"`
	GptKeys              []string `yaml:"GptKeys"`
	RedisHost            string   `yaml:"RedisHost"`
	RedisPass            string   `yaml:"RedisPass"`
	RedisDB              int      `yaml:"RedisDB"`
	CqHttpHost           string   `yaml:"CqHttpHost"`
	CqHttpPath           string   `yaml:"CqHttpPath"`
	AhuCalendarStartDate string   `yaml:"AhuCalendarStartDate"`
}
