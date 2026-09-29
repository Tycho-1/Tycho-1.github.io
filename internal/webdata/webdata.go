// Package webdata writes JSON files consumed by WebMCP tools in the browser.
package webdata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Tycho-1/gh-pages-portfolio/internal/projects"
	"github.com/Tycho-1/gh-pages-portfolio/internal/resources"
)

// ProjectSummary is the public list view of a project.
type ProjectSummary struct {
	Title   string   `json:"title"`
	Slug    string   `json:"slug"`
	Tags    []string `json:"tags"`
	Summary string   `json:"summary"`
}

// ProjectDetail is the full published metadata for one project.
type ProjectDetail struct {
	Title    string   `json:"title"`
	Slug     string   `json:"slug"`
	Tags     []string `json:"tags"`
	Summary  string   `json:"summary"`
	Repo     string   `json:"repo,omitempty"`
	Featured bool     `json:"featured"`
	URL      string   `json:"url"`
}

// Write emits /data/projects.json, /data/projects/<slug>.json, and /data/resources.json.
func Write(outDir string, published []projects.Project, resourceCategories []resources.Category) error {
	dataDir := filepath.Join(outDir, "data")
	projectsDir := filepath.Join(dataDir, "projects")
	if err := os.MkdirAll(projectsDir, 0o755); err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}

	summaries := make([]ProjectSummary, 0, len(published))
	for _, p := range published {
		summaries = append(summaries, ProjectSummary{
			Title:   p.Title,
			Slug:    p.Slug,
			Tags:    p.Tags,
			Summary: p.Summary,
		})
	}

	if err := writeJSON(filepath.Join(dataDir, "projects.json"), summaries); err != nil {
		return fmt.Errorf("write projects.json: %w", err)
	}

	for _, p := range published {
		detail := ProjectDetail{
			Title:    p.Title,
			Slug:     p.Slug,
			Tags:     p.Tags,
			Summary:  p.Summary,
			Repo:     p.Repo,
			Featured: p.Featured,
			URL:      "/projects/" + p.Slug + "/",
		}
		if err := writeJSON(filepath.Join(projectsDir, p.Slug+".json"), detail); err != nil {
			return fmt.Errorf("write project %q: %w", p.Slug, err)
		}
	}

	if err := writeJSON(filepath.Join(dataDir, "resources.json"), resourceCategories); err != nil {
		return fmt.Errorf("write resources.json: %w", err)
	}

	return nil
}

func writeJSON(path string, v any) error {
	raw, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}
	raw = append(raw, '\n')

	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}
