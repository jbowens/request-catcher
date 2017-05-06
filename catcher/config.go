package catcher

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	Port        int
	Host        string
	RootHost    string `json:"root_host"`
	Database    *Database
	LetsEncrypt *LetsEncrypt `json:"lets_encrypt"`
}

func (c *Configuration) ListenAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type LetsEncrypt struct {
	Email          string `json:"email"`
	CertsDirectory string `json:"certs_directory"`
}

type Database struct {
	Port int
	Name string
	Host string
	User string
}

func LoadConfiguration(filename string) (*Configuration, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	config := Configuration{}
	err = decoder.Decode(&config)

	return &config, err
}
