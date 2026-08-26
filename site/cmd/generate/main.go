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
	outDir := filepath.Join("site", "dist")

	if err := run(templatesDir, assetsDir, outDir); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if err := writeREADMETOC("README.md", content.All); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println("generated site into", outDir)
}

func run(templatesDir, assetsDir, outDir string) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
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

	for _, chapter := range content.All {
		data := PageData{
			Chapter:       chapter,
			HeroAvailable: fileExists(filepath.Join(assetsDir, filepath.FromSlash(chapter.HeroImage))),
			Roadmap:       roadmap,
			Labs:          content.Labs,
			AssetPrefix:   "",
		}

		outPath := filepath.Join(outDir, fmt.Sprintf("chapter-%d.html", chapter.Number))
		if err := renderPage(tmpl, outPath, data); err != nil {
			return err
		}
	}

	if err := copyAssets(assetsDir, outDir); err != nil {
		return err
	}

	firstChapterPath := filepath.Join(outDir, fmt.Sprintf("chapter-%d.html", content.All[0].Number))
	return copyFile(firstChapterPath, filepath.Join(outDir, "index.html"))
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
	return copyImages(filepath.Join(assetsDir, "images"), filepath.Join(outDir, "images"))
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

// copyImages copies every top-level .webp file directly under srcDir (the
// flat chapter-<N>.webp convention) into dstDir. It does not recurse into
// subdirectories, so scratch/draft folders alongside the real assets are
// never accidentally published.
func copyImages(srcDir, dstDir string) error {
	entries, err := os.ReadDir(srcDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading %s: %w", srcDir, err)
	}

	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", dstDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".webp") {
			continue
		}
		if err := copyFile(filepath.Join(srcDir, entry.Name()), filepath.Join(dstDir, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}
