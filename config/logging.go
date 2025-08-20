package config

/* logging parameters */
type Logging struct {
	File       string `yaml:"file,omitempty"`
	MaxSize    int    `yaml:"max_size,omitempty"`
	MaxBackups int    `yaml:"max_backups,omitempty"`
	MaxAge     int    `yaml:"max_age,omitempty"`
	Compress   bool   `yaml:"compress,omitempty"`
}

/* normalization function */
func (l *Logging) Normalize() error {

	/* set default file to log/app.log */
	if l.File == "" {
		l.File = "log/app.log"
	}

	/* set default max size to 100MB */
	if l.MaxSize == 0 {
		l.MaxSize = 100
	}

	/* set default max backups to 3 */
	if l.MaxBackups == 0 {
		l.MaxBackups = 3
	}

	/* let compression remain false by default */

	return nil
}
