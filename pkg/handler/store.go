package handler

import (
	"encoding/json"
	"fmt"
	"github.com/neucn/ipgw/pkg/model"
	"github.com/neucn/ipgw/pkg/utils"
	"os"
)

type StoreHandler struct {
	Path   string
	Config *model.Config
}

func NewStoreHandler(path string) (*StoreHandler, error) {
	path, err := getConfigPath(path)
	if err != nil {
		return nil, fmt.Errorf("fail to open config:\n\t%s", err)
	}
	return &StoreHandler{
		Path: path,
	}, nil
}

func getConfigPath(path string) (string, error) {
	if path == "" {
		homeDir, err := utils.GetHomeDir()
		if err != nil {
			return "", fmt.Errorf("fail to get home dir %v", err)
		}
		path = homeDir + string(os.PathSeparator) + ".ipgw"
	}
	if err := utils.FileMustExist(path); err != nil {
		return "", err
	}
	return path, nil
}

func (h *StoreHandler) Persist() error {
	content, err := json.Marshal(h.Config)
	if err != nil {
		return fmt.Errorf("fail to save config:\n\t%v", err)
	}
	return os.WriteFile(h.Path, content, 0666)
}

func (h *StoreHandler) Load() error {
	content, err := os.ReadFile(h.Path)
	if err != nil {
		return fmt.Errorf("fail to load config:\n\t%v", err)
	}
	h.Config = &model.Config{}
	err = json.Unmarshal(content, &h.Config)
	if err == nil || err.Error() == "unexpected end of JSON input" {
		return nil
	}
	return fmt.Errorf("fail to save config:\n\t%v", err)
}
