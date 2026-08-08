package config  


import (
	"os"
	"encoding/json"
	"path/filepath"
)


const configFileName = ".gatorconfig.json"


type Config struct {
	DbUrl string `json:"db_url"`
	UserName string `json:"current_user_name"`
}


func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home,configFileName), nil
}


func write(cfg Config) error {
	jsonConfig, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	os.WriteFile(filePath, jsonConfig, 0600)
	return nil
}


func (cfg *Config) SetUser(name string) error {
	cfg.UserName = name
	err := write(*cfg)
	if err != nil {
		return err
	}
	return nil
}


func Read() (Config, error) {
	fileName, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	rawjsonToRead, err:= os.ReadFile(fileName)
	if err != nil {
		return Config{}, err
	}
	cfg := Config{}
	err = json.Unmarshal([]byte(rawjsonToRead), &cfg)
	if err != nil {
		return Config{}, err
	}
 	return cfg, nil
}



