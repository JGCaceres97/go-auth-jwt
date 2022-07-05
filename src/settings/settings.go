package settings

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed settings.yaml
var settingsFile []byte

var JWTSecret, Port *string
var DB *Database

type Database struct {
	Host     *string `yaml:"host"`
	Port     *uint16 `yaml:"port"`
	User     *string `yaml:"user"`
	Password *string `yaml:"password"`
	Name     *string `yaml:"name"`
}

type Settings struct {
	Port      string   `yaml:"port"`
	DB        Database `yaml:"database"`
	JWTSecret string   `yaml:"jwtsecret"`
}

func New() error {
	s := Settings{}

	err := yaml.Unmarshal(settingsFile, &s)
	if err != nil {
		return err
	}

	JWTSecret = &s.JWTSecret
	Port = &s.Port
	DB = &s.DB

	return nil
}
