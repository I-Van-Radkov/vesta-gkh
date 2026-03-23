package news

import "time"

type ParserConfig struct {
	Threshold float64       `env:"PARSER_THRESHOLD" env-default:"130"`
	Interval  time.Duration `env:"PARSER_INTERVAL" env-default:"120s"`

	TitleWeight        float64 `env:"PARSER_TITLE_WEIGHT" env-default:"3.0"`
	DescWeight         float64 `env:"PARSER_DESC_WEIGHT" env-default:"1.2"`
	ContentWeight      float64 `env:"PARSER_CONTENT_WEIGHT" env-default:"0.6"`
	MaxContribKeyword  float64 `env:"PARSER_MAX_CONTRIB_KEYWORD" env-default:"3.5"`
	MaxContribCategory float64 `env:"PARSER_MAX_CONTRIB_CATEGORY" env-default:"12.0"`
	MaxContribRegex    float64 `env:"PARSER_MAX_CONTRIB_REGEX" env-default:"85.0"`
	DampeningFactor    float64 `env:"PARSER_DAMPENINF_FACTOR" env-default:"0.35"`
}
