package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"go-newpost/scraper"
)

type Storage struct {
	filename string
}

type Data struct {
	LastChecked string          `json:"last_checked"`
	SeenPosts   map[string]bool `json:"seen_posts"`
}

func New(filename string) *Storage {
	return &Storage{
		filename: filename,
	}
}

func (s *Storage) LoadData() (*Data, error) {
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		return &Data{
			SeenPosts: make(map[string]bool),
		}, nil
	}

	file, err := os.Open(s.filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var data Data
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	if data.SeenPosts == nil {
		data.SeenPosts = make(map[string]bool)
	}

	return &data, nil
}

func (s *Storage) SaveData(data *Data) error {
	file, err := os.Create(s.filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

func (s *Storage) GetNewPosts(posts []scraper.Post) ([]scraper.Post, error) {
	data, err := s.LoadData()
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}

	var newPosts []scraper.Post
	for _, post := range posts {
		if !data.SeenPosts[post.ID] {
			newPosts = append(newPosts, post)
			data.SeenPosts[post.ID] = true
		}
	}

	if len(newPosts) > 0 {
		if err := s.SaveData(data); err != nil {
			return nil, fmt.Errorf("failed to save data: %w", err)
		}
	}

	return newPosts, nil
}