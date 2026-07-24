// Package state atomically persists monitor cards outside the repository.
package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rex/apple-submission-monitor/internal/domain"
)

const currentVersion = 1

// Store reads and atomically writes versioned monitor state.
type Store struct {
	path string
}

type fileState struct {
	Version int                 `json:"version"`
	Cards   []domain.Submission `json:"cards"`
}

// NewStore creates a state store for an explicit local file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Load returns an empty state when the file has not been created.
func (s *Store) Load() ([]domain.Submission, error) {
	content, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read monitor state: %w", err)
	}
	var persisted fileState
	if err := json.Unmarshal(content, &persisted); err != nil {
		return nil, fmt.Errorf("decode monitor state: %w", err)
	}
	if persisted.Version != currentVersion {
		return nil, fmt.Errorf("unsupported monitor state version %d", persisted.Version)
	}
	return persisted.Cards, nil
}

// Save writes cards through a private temporary file and atomic rename.
func (s *Store) Save(cards []domain.Submission) error {
	directory := filepath.Dir(s.path)
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return fmt.Errorf("create state directory: %w", err)
	}
	temp, err := os.CreateTemp(directory, ".state-*.json")
	if err != nil {
		return fmt.Errorf("create temporary state: %w", err)
	}
	tempName := temp.Name()
	closed := false
	defer func() {
		if !closed {
			_ = temp.Close()
		}
		_ = os.Remove(tempName)
	}()
	if err := temp.Chmod(0o600); err != nil {
		return fmt.Errorf("protect temporary state: %w", err)
	}
	encoder := json.NewEncoder(temp)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(fileState{Version: currentVersion, Cards: cards}); err != nil {
		return fmt.Errorf("encode monitor state: %w", err)
	}
	if err := temp.Sync(); err != nil {
		return fmt.Errorf("sync monitor state: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close monitor state: %w", err)
	}
	closed = true
	if err := os.Rename(tempName, s.path); err != nil {
		return fmt.Errorf("replace monitor state: %w", err)
	}
	if err := os.Chmod(s.path, 0o600); err != nil {
		return fmt.Errorf("protect monitor state: %w", err)
	}
	return nil
}
