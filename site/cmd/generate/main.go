package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"site/content"
)

func main() {
	templatesDir := filepath.Join("site", "templates")
	assetsDir := filepath.Join("site", "assets")
	outDir := "docs"
	indexPath := filepath.Join(outDir, "index.html")

	if err := run(templatesDir, assetsDir, outDir, indexPath); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if err := writeREADMETOC("README.md", content.All); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("generated site into %s, entry point %s\n", outDir, indexPath)
}

// run renders every chapter into outDir and writes the site's entry point to
// indexPath. The entry point is rendered rather than copied from a chapter
// page, because its asset prefix depends on where it sits relative to outDir,
// which assetPrefix works out.
//
// Images are not copied here. They live in outDir already, checked in beside
// the pages that serve them, because GitHub Pages publishes outDir alone and
// nothing outside it is reachable.
func run(templatesDir, assetsDir, outDir, indexPath string) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	if err := requireImages(outDir); err != nil {
		return err
	}

	tmpl, err := template.ParseFiles(
		filepath.Join(templatesDir, "sidebar.html.tmpl"),
		filepath.Join(templatesDir, "page.html.tmpl"),
	)
	if err != nil {
		return fmt.Errorf("parsing templates: %w", err)
	}

	built := map[int]bool{}
	for _, c := range content.All {
		built[c.Number] = true
	}

	roadmap := make([]roadmapView, 0, len(content.Roadmap))
	for _, r := range content.Roadmap {
		roadmap = append(roadmap, roadmapView{Number: r.Number, Title: r.Title, Built: built[r.Number]})
	}

	if len(content.All) == 0 {
		return fmt.Errorf("no chapters defined in content.All")
	}

	pageData := func(chapter content.ChapterContent, prefix string) PageData {
		return PageData{
			Chapter:       chapter,
			HeroAvailable: fileExists(filepath.Join(outDir, filepath.FromSlash(chapter.HeroImage))),
			Roadmap:       roadmap,
			Labs:          content.Labs,
			AssetPrefix:   prefix,
		}
	}

	for _, chapter := range content.All {
		outPath := filepath.Join(outDir, fmt.Sprintf("chapter-%d.html", chapter.Number))
		if err := renderPage(tmpl, outPath, pageData(chapter, "")); err != nil {
			return err
		}
	}

	if err := copyAssets(assetsDir, outDir); err != nil {
		return err
	}

	// The entry point is the first chapter rendered a second time rather than
	// copied from its page, so that it can carry whatever asset prefix its own
	// location calls for.
	prefix, err := assetPrefix(indexPath, outDir)
	if err != nil {
		return err
	}
	return renderPage(tmpl, indexPath, pageData(content.All[0], prefix))
}

// assetPrefix is the URL prefix a page at pagePath needs in order to reach
// assets published in outDir: "docs/" for the root index.html, "" for a page
// already sitting inside outDir.
func assetPrefix(pagePath, outDir string) (string, error) {
	rel, err := filepath.Rel(filepath.Dir(pagePath), outDir)
	if err != nil {
		return "", fmt.Errorf("relating %s to %s: %w", pagePath, outDir, err)
	}
	if rel == "." {
		return "", nil
	}
	return filepath.ToSlash(rel) + "/", nil
}

// requireImages refuses to build when the illustrations are not where they
// should be. They are the one part of the published site nothing can rebuild,
// so a run that cannot find them is looking at a wiped output directory, not
// a fresh one.
//
// Without this the build would succeed: every chapter would simply render its
// "image pending" placeholder, and a working site would quietly be replaced by
// twenty-one empty panels. Nothing reads these files except the browser, so
// the failure would surface only once it was published.
func requireImages(outDir string) error {
	dir := filepath.Join(outDir, "images")

	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist. It holds the chapter illustrations, "+
			"which the generator does not create and cannot replace. Restore it with "+
			"'git checkout -- %s' before building", dir, filepath.ToSlash(dir))
	}
	if err != nil {
		return fmt.Errorf("reading %s: %w", dir, err)
	}

	for _, e := range entries {
		if !e.IsDir() && strings.EqualFold(filepath.Ext(e.Name()), ".webp") {
			return nil
		}
	}
	return fmt.Errorf("%s holds no .webp files. The chapter illustrations are missing "+
		"and the generator cannot replace them. Restore them with 'git checkout -- %s' "+
		"before building", dir, filepath.ToSlash(dir))
}

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func renderPage(tmpl *template.Template, outPath string, data PageData) error {
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("creating %s: %w", outPath, err)
	}
	defer f.Close()

	if err := tmpl.ExecuteTemplate(f, "page", data); err != nil {
		return fmt.Errorf("rendering %s: %w", outPath, err)
	}
	return nil
}

func copyAssets(assetsDir, outDir string) error {
	files := []string{"styles.css", "app.js"}
	for _, name := range files {
		if err := copyFile(filepath.Join(assetsDir, name), filepath.Join(outDir, name)); err != nil {
			return err
		}
	}

	// GitHub Pages runs the published folder through Jekyll unless this file
	// is present. It has to be written here rather than kept by hand, or
	// deleting outDir and regenerating silently drops it and Jekyll starts
	// rewriting the output.
	if err := os.WriteFile(filepath.Join(outDir, ".nojekyll"), nil, 0o644); err != nil {
		return fmt.Errorf("writing .nojekyll: %w", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading %s: %w", src, err)
	}
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", dst, err)
	}
	return nil
}
