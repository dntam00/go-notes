package logging

type LogConfig struct {
	Path  string `mapstructure:"path"`  // application path
	Level string `mapstructure:"level"` // debug, info, warn, error, fatal, panic

	Formatter  string `mapstructure:"formatter"`   // console, json
	MaxSize    int    `mapstructure:"max_size"`    // file max size in MB
	MaxBackups int    `mapstructure:"max_backups"` // max files is kept
}
