package filewriter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/taka1156/bastion/internal/domain/entity"
)

type LocalFileWriter struct{}

func NewLocalFileWriter() *LocalFileWriter {
	return &LocalFileWriter{}
}

func (w *LocalFileWriter) Write(config entity.BastionConfig, outputDir string) error {
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0o750); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	outputPath := fmt.Sprintf("%s/bastion.json", outputDir)
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close file: %v\n", err)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to write config to file: %w", err)
	}

	fmt.Printf("Config generated at: %s\n", outputPath)
	return nil
}
