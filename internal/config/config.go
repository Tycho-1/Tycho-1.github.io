package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Site holds site-wide metadata shown on every page.
type Site struct {
	Name     string `yaml:"name"`
	Brand    string `yaml:"brand"`
	Title    string `yaml:"title"`
	Tagline  string `yaml:"tagline"`
	GitHub   string `yaml:"github"`
	LinkedIn string `yaml:"linkedin"`
	Email    string `yaml:"email"`
}

// BrandName returns the short header brand, falling back to Name.
func (s Site) BrandName() string {
	if s.Brand != "" {
		return s.Brand
	}
	return s.Name
}

// LoadSite reads site metadata from a YAML file.
// Brand can be overridden with the PORTFOLIO_BRAND environment variable.
func LoadSite(path string) (Site, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Site{}, fmt.Errorf("read file: %w", err)
	}

	var site Site
	if err := yaml.Unmarshal(raw, &site); err != nil {
		return Site{}, fmt.Errorf("parse yaml: %w", err)
	}

	if site.Name == "" {
		return Site{}, fmt.Errorf("site name is required")
	}

	if v := os.Getenv("PORTFOLIO_BRAND"); v != "" {
		site.Brand = v
	}

	return site, nil
}
