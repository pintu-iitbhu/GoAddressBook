package springcloud

import (
	"encoding/json"
	"os"
)

func GetFile(path string) (SpringResponse, error) {
	var cloudConfig SpringResponse
	file, err := os.ReadFile(path)
	if err != nil {
		return SpringResponse{}, err
	}
	err = json.Unmarshal(file, &cloudConfig)
	if err != nil {
		return SpringResponse{}, err
	}
	return cloudConfig, nil
}
