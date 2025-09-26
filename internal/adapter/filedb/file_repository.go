package filedb

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// FileRepository is a driven adapter that persists to a newline-delimited text file.
type FileRepository struct {
	path string
	mu   sync.Mutex // protect concurrent CLI calls
}

func New(path string) *FileRepository {
	return &FileRepository{path: path}
}

func (f *FileRepository) checkFile() error {

	if f.path == "" {
		return errors.New("file path is empty")
	}

	err := os.MkdirAll(filepath.Dir(f.path), 0o755)
	if err != nil {
		return err
	}

	if _, err := os.Stat(f.path); os.IsNotExist(err) {
		file, err := os.OpenFile(f.path, os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		_ = file.Close()
	}

	return nil
}

func (f *FileRepository) Add(value string) error {
	// Lock Writing to the file
	f.mu.Lock()
	// Unlock before return
	defer f.mu.Unlock()

	if err := f.checkFile(); err != nil {
		return err
	}

	value = strings.ReplaceAll(value, "\n", " ")

	file, err := os.OpenFile(f.path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(value + "\n"); err != nil {
		return err
	}

	return nil
}

func (f *FileRepository) List() ([]string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if err := f.checkFile(); err != nil {
		return nil, err
	}

	file, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	s := bufio.NewScanner(file)

	for s.Scan() {
		lines = append(lines, s.Text())
	}

	if err := s.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func (f *FileRepository) DeleteLast() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if err := f.checkFile(); err != nil {
		return err
	}

	// Read all lines
	file, err := os.Open(f.path)
	if err != nil {
		return err
	}
	var lines []string
	s := bufio.NewScanner(file)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	_ = file.Close()
	if err := s.Err(); err != nil {
		return err
	}

	if len(lines) == 0 {
		return errors.New("nothing to delete: file is empty")
	}
	// Drop the last non-empty line (treat trailing empties as noise)
	idx := len(lines) - 1
	for idx >= 0 && strings.TrimSpace(lines[idx]) == "" {
		idx--
	}
	if idx < 0 { // all blank
		// Truncate file to empty
		return os.WriteFile(f.path, []byte{}, 0o644)
	}
	newContent := strings.Join(lines[:idx], "\n")
	if newContent != "" {
		newContent += "\n"
	}
	return os.WriteFile(f.path, []byte(newContent), 0o644)
}
