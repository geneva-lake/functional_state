package service

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

//   - -------------------------------------------------------------------------------------------------------------------
//     Configuration for service
//   - -------------------------------------------------------------------------------------------------------------------
type Config struct {
	Name               string `yaml:"name"`
	DBConnectionString string `yaml:"db_connection_sring"`
	FiancialsURL       string `yaml:"fiancials_url"`
	Port               string `yaml:"port"`
}

func NewConfig() *Config {
	return &Config{}
}

//   - -------------------------------------------------------------------------------------------------------------------
//     Envelope for proceeding main struct
//   - -------------------------------------------------------------------------------------------------------------------
type config struct {
	object *Config
	reader io.Reader
	Error  error
}

//   - -------------------------------------------------------------------------------------------------------------------
//     Create envelope from reader
//   - -------------------------------------------------------------------------------------------------------------------
func (this *Config) from(reader io.Reader, err error) *config {
	return &config{
		object: this,
		reader: reader,
		Error:  err,
	}
}

//   - -------------------------------------------------------------------------------------------------------------------
//     Set file
//   - -------------------------------------------------------------------------------------------------------------------
func (this *Config) FromFile(name string) *config {
	return this.from(os.Open(name))
}

//   - -------------------------------------------------------------------------------------------------------------------
//     Set reader
//   - -------------------------------------------------------------------------------------------------------------------
func (this *Config) FromReader(reader io.Reader) *config {
	return this.from(reader, nil)
}

//   - -------------------------------------------------------------------------------------------------------------------
//     Parse yaml file
//   - -------------------------------------------------------------------------------------------------------------------
func (this *config) Yaml() (*Config, error) {
	if this.Error != nil {
		return nil, this.Error
	}
	if this.reader == nil {
		return nil, nil
	}
	return this.object, yaml.NewDecoder(this.reader).Decode(this.object)
}
