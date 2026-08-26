package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// testRun renders into a fresh temp root, returning the docs dir the pages
// land in and the path of the index.html that fronts them.
func testRun(t *testing.T) (outDir, indexPath string) {
	t.Helper()
	outDir = filepath.Join(t.TempDir(), "docs")
	indexPath = filepath.Join(outDir, "index.html")

	if err := run(filepath.Join("..", "..", "templates"), filepath.Join("..", "..", "assets"), outDir, indexPath); err != nil {
		t.Fatalf("run: %v", err)
	}
	return outDir, indexPath
}

func TestRunGeneratesChapterZeroPage(t *testing.T) {
	outDir, _ := testRun(t)

	pageHTML, err := os.ReadFile(filepath.Join(outDir, "chapter-0.html"))
	if err != nil {
		t.Fatalf("reading chapter-0.html: %v", err)
	}
	body := string(pageHTML)

	for _, want := range []string{"What Are We Building?", "1. Intuition", "2. Concepts", "3. Under the Hood", "4. Build It", "5. Break It", "Mark chapter as complete", "Build the Vault", "Inside One Computer", "Zero-shot prompting", "Documented in The Prompt Report"} {
		if !strings.Contains(body, want) {
			t.Errorf("chapter-0.html missing expected text %q", want)
		}
	}

	if _, err := os.Stat(filepath.Join(outDir, "styles.css")); err != nil {
		t.Errorf("styles.css was not copied: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outDir, "app.js")); err != nil {
		t.Errorf("app.js was not copied: %v", err)
	}

	for _, unwanted := range []string{"http://", "https://", "fetch("} {
		if strings.Contains(body, unwanted) {
			t.Errorf("chapter-0.html unexpectedly contains %q — site must be fully self-contained", unwanted)
		}
	}
}

// index.html sits alongside the pages it fronts, so Pages serving docs/ as
// the site root resolves a bare / to it and every URL below is a clean
// /chapter-N.html. Sharing a directory with its assets, it references them
// with no prefix — the same way every chapter page does.
func TestRunWritesIndexBesideTheChaptersItFronts(t *testing.T) {
	outDir, indexPath := testRun(t)

	if filepath.Dir(indexPath) != outDir {
		t.Fatalf("index.html is at %s, want it inside %s", indexPath, outDir)
	}

	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("reading index.html: %v", err)
	}
	body := string(data)

	for _, want := range []string{
		`href="styles.css"`,
		`src="app.js"`,
		`src="images/chapter-0.webp"`,
		`href="chapter-1.html"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("index.html missing %q", want)
		}
	}

	// A docs/ prefix would mean the entry point still thinks it lives a
	// directory above its assets, which 404s once docs/ is the site root.
	if strings.Contains(body, `"docs/`) {
		t.Error(`index.html still references assets via a "docs/" prefix`)
	}
}

func TestCopyImagesOnlyCopiesFlatWebPs(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := filepath.Join(t.TempDir(), "images")

	if err := os.WriteFile(filepath.Join(srcDir, "chapter-0.webp"), []byte("fake-webp"), 0o644); err != nil {
		t.Fatalf("writing fixture webp: %v", err)
	}
	if err := os.Mkdir(filepath.Join(srcDir, "_pending"), 0o755); err != nil {
		t.Fatalf("creating fixture subdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "_pending", "draft.webp"), []byte("fake-webp"), 0o644); err != nil {
		t.Fatalf("writing fixture draft webp: %v", err)
	}

	if err := copyImages(srcDir, dstDir); err != nil {
		t.Fatalf("copyImages: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dstDir, "chapter-0.webp")); err != nil {
		t.Errorf("expected chapter-0.webp to be copied: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dstDir, "_pending")); err == nil {
		t.Error("expected _pending subdirectory to NOT be copied, but it was")
	}
}

// docs/ is disposable: deleting it and regenerating has to reproduce every
// published file. .nojekyll was hand-made once and vanished on the first
// clean rebuild, which would have let Jekyll rewrite the output.
func TestRunWritesNoJekyllSoOutputSurvivesACleanRebuild(t *testing.T) {
	outDir, _ := testRun(t)

	if _, err := os.Stat(filepath.Join(outDir, ".nojekyll")); err != nil {
		t.Errorf(".nojekyll was not generated: %v", err)
	}
}
