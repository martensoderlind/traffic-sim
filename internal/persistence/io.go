package persistence

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/sqweek/dialog"
)

const SaveDir = "saves"

func init() {
	if err := os.MkdirAll(SaveDir, 0755); err != nil {
		log.Printf("Warning: Could not create saves directory: %v", err)
	}
}

func SaveToFile(saveData *SaveFormat) error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	defaultFilename := filepath.Join(SaveDir, fmt.Sprintf("simulation_%s.json", timestamp))

	filename, err := dialog.File().
		Title("Save Simulation").
		Filter("JSON files", "json").
		SetStartFile(defaultFilename).
		Save()

	if err != nil {
		if err == dialog.ErrCancelled {
			log.Println("Save cancelled by user")
			return nil
		}
		return fmt.Errorf("file dialog error: %w", err)
	}

	if filepath.Ext(filename) != ".json" {
		filename += ".json"
	}

	data, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal save data: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	log.Printf("Simulation saved to: %s", filename)
	return nil
}

func LoadFromFile() (*SaveFormat, error) {
	filename, err := dialog.File().
		Title("Load Simulation").
		Filter("JSON files", "json").
		SetStartDir(SaveDir).
		Load()

	if err != nil {
		if err == dialog.ErrCancelled {
			log.Println("Load cancelled by user")
			return nil, nil
		}
		return nil, fmt.Errorf("file dialog error: %w", err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var saveData SaveFormat
	if err := json.Unmarshal(data, &saveData); err != nil {
		return nil, fmt.Errorf("failed to parse save file: %w", err)
	}

	log.Printf("Simulation loaded from: %s", filename)
	return &saveData, nil
}