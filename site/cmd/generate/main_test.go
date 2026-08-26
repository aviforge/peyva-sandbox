package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunGeneratesChapterZeroPage(t *testing.T) {
	outDir := t.TempDir()

	if err := run(filepath.Join("..", "..", "templates"), filepath.Join("..", "..", "assets"), outDir); err != nil {
		t.Fatalf("run: %v", err)
	}

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

	if _, err := os.Stat(filepath.Join(outDir, "index.html")); err != nil {
		t.Errorf("index.html was not generated: %v", err)
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
