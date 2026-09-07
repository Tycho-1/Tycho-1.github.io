package projects

import (
	"fmt"
	"os"
	"sort"

	"gopkg.in/yaml.v3"
)

// Project describes one portfolio entry.
type Project struct {
	Slug     string   `yaml:"slug"`
	Title    string   `yaml:"title"`
	Summary  string   `yaml:"summary"`
	Tags     []string `yaml:"tags"`
	Repo     string   `yaml:"repo"`
	Featured bool     `yaml:"featured"`
	Draft    bool     `yaml:"draft"`
	Order    int      `yaml:"order"`
}

// List is the full project catalog from YAML.
type List struct {
	Items []Project
}

// Load reads projects from a YAML file.
func Load(path string) (List, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return List{}, fmt.Errorf("read file: %w", err)
	}

	var items []Project
	if err := yaml.Unmarshal(raw, &items); err != nil {
		return List{}, fmt.Errorf("parse yaml: %w", err)
	}

	for i := range items {
		if items[i].Slug == "" {
			return List{}, fmt.Errorf("project at index %d is missing slug", i)
		}
		if items[i].Title == "" {
			return List{}, fmt.Errorf("project %q is missing title", items[i].Slug)
		}
	}

	return List{Items: items}, nil
}

// Published returns non-draft projects sorted by order then title.
func (l List) Published() []Project {
	out := make([]Project, 0, len(l.Items))
	for _, p := range l.Items {
		if !p.Draft {
			out = append(out, p)
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Order != out[j].Order {
			return out[i].Order < out[j].Order
		}
		return out[i].Title < out[j].Title
	})

	return out
}

// Featured returns published projects marked featured.
func (l List) Featured() []Project {
	var out []Project
	for _, p := range l.Published() {
		if p.Featured {
			out = append(out, p)
		}
	}
	return out
}
