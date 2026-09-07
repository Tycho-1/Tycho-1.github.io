package resources

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Item is one external resource link.
type Item struct {
	Title string `yaml:"title"`
	URL   string `yaml:"url"`
	Note  string `yaml:"note"`
}

// Category groups related resources.
type Category struct {
	Name  string `yaml:"category"`
	Items []Item `yaml:"items"`
}

// Load reads grouped resources from a YAML file.
func Load(path string) ([]Category, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var categories []Category
	if err := yaml.Unmarshal(raw, &categories); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	for i, cat := range categories {
		if cat.Name == "" {
			return nil, fmt.Errorf("category at index %d is missing name", i)
		}
		for j, item := range cat.Items {
			if item.Title == "" {
				return nil, fmt.Errorf("category %q item at index %d is missing title", cat.Name, j)
			}
			if item.URL == "" {
				return nil, fmt.Errorf("category %q item %q is missing url", cat.Name, item.Title)
			}
		}
	}

	return categories, nil
}
