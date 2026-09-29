// Package site builds the static portfolio from YAML, Markdown, and templates.
package site

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Tycho-1/gh-pages-portfolio/internal/config"
	"github.com/Tycho-1/gh-pages-portfolio/internal/markdown"
	"github.com/Tycho-1/gh-pages-portfolio/internal/projects"
	"github.com/Tycho-1/gh-pages-portfolio/internal/render"
	"github.com/Tycho-1/gh-pages-portfolio/internal/resources"
	"github.com/Tycho-1/gh-pages-portfolio/internal/static"
	"github.com/Tycho-1/gh-pages-portfolio/internal/webdata"
)

// Config holds paths and metadata for site generation.
type Config struct {
	RootDir string
	OutDir  string
}

// Generator renders the portfolio into OutDir.
type Generator struct {
	cfg Config
}

// NewGenerator returns a generator rooted at rootDir.
func NewGenerator(rootDir, outDir string) *Generator {
	return &Generator{cfg: Config{RootDir: rootDir, OutDir: outDir}}
}

// Build generates the full static site.
func (g *Generator) Build() error {
	root := g.cfg.RootDir
	out := g.cfg.OutDir

	siteMeta, err := config.LoadSite(filepath.Join(root, "data", "site.yaml"))
	if err != nil {
		return fmt.Errorf("load site config: %w", err)
	}

	list, err := projects.Load(filepath.Join(root, "data", "projects.yaml"))
	if err != nil {
		return fmt.Errorf("load projects: %w", err)
	}

	published := list.Published()
	if len(published) == 0 {
		return fmt.Errorf("no published projects found")
	}

	converter := markdown.NewConverter()
	renderer, err := render.New(filepath.Join(root, "templates"))
	if err != nil {
		return fmt.Errorf("init templates: %w", err)
	}

	if err := os.RemoveAll(out); err != nil {
		return fmt.Errorf("clean output dir: %w", err)
	}
	if err := os.MkdirAll(out, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	pages, err := g.loadProjectPages(root, published, converter)
	if err != nil {
		return fmt.Errorf("load project pages: %w", err)
	}

	if err := renderer.Home(out, render.HomeData{
		Page:     render.Page{Site: siteMeta},
		Projects: published,
	}); err != nil {
		return fmt.Errorf("render home: %w", err)
	}

	if err := renderer.ProjectsIndex(out, render.ProjectsIndexData{
		Page:     render.Page{Site: siteMeta, PageTitle: "Projects"},
		Projects: published,
	}); err != nil {
		return fmt.Errorf("render projects index: %w", err)
	}

	for _, page := range pages {
		if err := renderer.Project(out, render.ProjectData{
			Page: render.Page{
				Site:       siteMeta,
				PageTitle:  page.Meta.Title,
				WideLayout: true,
			},
			Project: page.Meta,
			Body:    template.HTML(page.HTML),
		}); err != nil {
			return fmt.Errorf("render project %q: %w", page.Meta.Slug, err)
		}
	}

	resourceCategories, err := resources.Load(filepath.Join(root, "data", "resources.yaml"))
	if err != nil {
		return fmt.Errorf("load resources: %w", err)
	}

	if err := renderer.Resources(out, render.ResourcesData{
		Page:       render.Page{Site: siteMeta, PageTitle: "Resources"},
		Categories: resourceCategories,
	}); err != nil {
		return fmt.Errorf("render resources: %w", err)
	}

	if err := webdata.Write(out, published, resourceCategories); err != nil {
		return fmt.Errorf("write web data: %w", err)
	}

	if err := static.Copy(filepath.Join(root, "static"), out); err != nil {
		return fmt.Errorf("copy static assets: %w", err)
	}

	if err := os.WriteFile(filepath.Join(out, ".nojekyll"), nil, 0o644); err != nil {
		return fmt.Errorf("write .nojekyll: %w", err)
	}

	return nil
}

type projectPage struct {
	Meta projects.Project
	HTML string
}

func (g *Generator) loadProjectPages(root string, published []projects.Project, converter *markdown.Converter) ([]projectPage, error) {
	contentDir := filepath.Join(root, "content", "projects")
	pages := make([]projectPage, 0, len(published))

	for _, project := range published {
		path := filepath.Join(contentDir, project.Slug+".md")
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", path, err)
		}

		doc, err := markdown.ParseDocument(raw)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}

		html, err := converter.ToHTML(doc.Body)
		if err != nil {
			return nil, fmt.Errorf("convert %s: %w", project.Slug, err)
		}

		pages = append(pages, projectPage{
			Meta: project,
			HTML: html,
		})
	}

	return pages, nil
}

// Root returns the configured source root, or the current directory when unset.
func Root() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return wd, nil
}

// EnsureOutDir creates the output directory if needed.
func EnsureOutDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// WalkDist returns a fs.FS for serving generated files.
func WalkDist(outDir string) (fs.FS, error) {
	return os.DirFS(outDir), nil
}
