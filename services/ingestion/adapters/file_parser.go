package adapters

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/testpilot-ai/ingestion/domain/entities"
	"gopkg.in/yaml.v3"
)

// FileParser handles parsing of API configuration files
type FileParser struct{}

// NewFileParser creates a new file parser
func NewFileParser() *FileParser {
	return &FileParser{}
}

// ParseFile parses a YAML or JSON configuration file
func (p *FileParser) ParseFile(filePath string) (*entities.APIConfig, string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file: %w", err)
	}

	// Calculate content hash
	hash := sha256.Sum256(data)
	contentHash := hex.EncodeToString(hash[:])

	var config entities.APIConfig

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, "", fmt.Errorf("failed to parse YAML: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, "", fmt.Errorf("failed to parse JSON: %w", err)
		}
	default:
		return nil, "", fmt.Errorf("unsupported file type: %s", ext)
	}

	return &config, contentHash, nil
}

// ScanFolder scans a folder for API configuration files
func (p *FileParser) ScanFolder(folderPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".yaml" || ext == ".yml" || ext == ".json" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan folder: %w", err)
	}

	return files, nil
}

// CalculateHash calculates SHA256 hash of content
func (p *FileParser) CalculateHash(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}

