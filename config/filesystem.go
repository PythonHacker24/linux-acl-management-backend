package config

import (
	"errors"

	"github.com/MakeNowJust/heredoc"
)

/* file system server parameters */
type FileSystemServers struct {
	Path   string  `yaml:"path,omitempty"`
	Method string  `yaml:"method,omitempty"`
	Remote *Remote `yaml:"remote,omitempty"`
}

/* remote parameters for file system server with laclm daemons installed */
type Remote struct {
	Host string `yaml:"host,omitempty"`
	Port int    `yaml:"port,omitempty"`
}

/* normalization function */
func (f *FileSystemServers) Normalize() error {
	if f.Path == "" {
		return errors.New(heredoc.Doc(`
			Remote server file path not specified in the configuration file. 

			Please check the docs for more information: 
		`))
	}

	/* set default method to local */
	if f.Method == "" {
		f.Method = "local"
	}

	/* check if method is remote */
	if f.Method == "remote" {
		/* check if remote is specified */
		if f.Remote == nil {
			return errors.New(heredoc.Doc(`
				Remote file server not specified in the configuration file. 

				Please check the docs for more information: 
			`))
		}

		/* check if host is specified */
		if f.Remote.Host == "" {
			return errors.New(heredoc.Doc(`
				Address not provided for remote file server
				
				Please check the docs for more information: 
			`))
		}

		/* check if port is specified */
		if f.Remote.Port == 0 {
			return errors.New(heredoc.Doc(`
				Port not provided for remote file server 	

				Please check the docs for more information: 
			`))
		}
	}

	return nil
}
