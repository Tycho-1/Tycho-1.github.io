package render

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/Tycho-1/gh-pages-portfolio/internal/config"
	"github.com/Tycho-1/gh-pages-portfolio/internal/projects"
	"github.com/Tycho-1/gh-pages-portfolio/internal/resources"
)

// Renderer executes HTML templates into dist/.
type Renderer struct {
	templates *template.Template
}

// New parses all templates from dir.
func New(dir string) (*Renderer, error) {
	funcs := template.FuncMap{
		"projectURL": func(slug string) string {
			return "/projects/" + slug + "/"
		},
		"brand": func(s config.Site) string {
			return s.BrandName()
		},
	}

	tmpl, err := template.New("").Funcs(funcs).ParseGlob(filepath.Join(dir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}

	return &Renderer{templates: tmpl}, nil
}

// Page is shared metadata passed into the layout template.
type Page struct {
	Site        config.Site
	PageTitle   string
	WideLayout  bool
}

// HomeData is the landing page view model.
type HomeData struct {
	Page
	Projects []projects.Project
}

// ProjectsIndexData is the projects listing view model.
type ProjectsIndexData struct {
	Page
	Projects []projects.Project
}

// ProjectData is a single project page view model.
type ProjectData struct {
	Page
	Project projects.Project
	Body    template.HTML
}

// ResourcesData is the resources page view model.
type ResourcesData struct {
	Page
	Categories []resources.Category
}

// Home writes dist/index.html.
func (r *Renderer) Home(outDir string, data HomeData) error {
	return r.write(filepath.Join(outDir, "index.html"), "home.html", data)
}

// ProjectsIndex writes dist/projects/index.html.
func (r *Renderer) ProjectsIndex(outDir string, data ProjectsIndexData) error {
	dir := filepath.Join(outDir, "projects")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create projects dir: %w", err)
	}
	return r.write(filepath.Join(dir, "index.html"), "projects.html", data)
}

// Project writes dist/projects/<slug>/index.html.
func (r *Renderer) Project(outDir string, data ProjectData) error {
	dir := filepath.Join(outDir, "projects", data.Project.Slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create project dir: %w", err)
	}
	return r.write(filepath.Join(dir, "index.html"), "project.html", data)
}

// Resources writes dist/resources/index.html.
func (r *Renderer) Resources(outDir string, data ResourcesData) error {
	dir := filepath.Join(outDir, "resources")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create resources dir: %w", err)
	}
	return r.write(filepath.Join(dir, "index.html"), "resources.html", data)
}

func (r *Renderer) write(path, name string, data any) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if err := r.templates.ExecuteTemplate(f, name, data); err != nil {
		return fmt.Errorf("execute template %q: %w", name, err)
	}
	return f.Close()
}
