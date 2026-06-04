package input

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/taka1156/bastion/internal/domain/entity"
)

type JsonInput struct{}

func NewJsonInput() *JsonInput {
	return &JsonInput{}
}

func (j *JsonInput) Load(path string) (entity.BastionConfig, error) {
	var config entity.BastionConfig
	file, err := os.Open(path)
	if err != nil {
		return config, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close file: %v\n", err)
		}
	}()

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return config, fmt.Errorf("failed to decode config file: %w", err)
	}

	return config, nil
}
