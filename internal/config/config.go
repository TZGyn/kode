package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

type Config struct {
	GEMINI_API_KEY   string `json:"GEMINI_API_KEY"`
	OPENAI_API_KEY   string `json:"OPENAI_API_KEY"`
	DEFAULT_PROVIDER string `json:"default_provider"`
	DEFAULT_MODEL    string `json:"default_model"`
}

func New() (*Config, error) {
	c := &Config{}

	configFilePath, err := xdg.ConfigFile(filepath.Join("kode", "kode.json"))
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(configFilePath)
	if err = os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}

	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(configFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		defaultConfig, _ := json.MarshalIndent(&c, "", "\t")
		_, err = f.WriteString(string(defaultConfig))
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(content, &c); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) SaveConfig() error {
	configFilePath, err := xdg.ConfigFile(filepath.Join("kode", "kode.json"))
	if err != nil {
		return err
	}

	dir := filepath.Dir(configFilePath)
	if err = os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	_, err = os.Stat(configFilePath)

	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(configFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		defaultConfig, _ := json.MarshalIndent(&c, "", "\t")
		_, err = f.WriteString(string(defaultConfig))
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	defaultConfig, _ := json.MarshalIndent(&c, "", "\t")
	err = os.WriteFile(configFilePath, defaultConfig, 0644)

	return err
}
