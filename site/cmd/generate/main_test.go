package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"site/content"
)

// testRun renders into a fresh temp root, returning the docs dir the pages
// land in and the path of the index.html that fronts them.
//
// The hero images are seeded first because the generator no longer copies
// them. They are checked in under the output directory, so a run that finds
// none renders every chapter with the "image pending" placeholder instead.
func testRun(t *testing.T) (outDir, indexPath string) {
	t.Helper()
	outDir = filepath.Join(t.TempDir(), "docs")
	indexPath = filepath.Join(outDir, "index.html")

	imagesDir := filepath.Join(outDir, "images")
	if err := os.MkdirAll(imagesDir, 0o755); err != nil {
		t.Fatalf("creating fixture images dir: %v", err)
	}
	for _, c := range content.All {
		name := filepath.Base(filepath.FromSlash(c.HeroImage))
		if err := os.WriteFile(filepath.Join(imagesDir, name), []byte("fake-webp"), 0o644); err != nil {
			t.Fatalf("writing fixture hero %s: %v", name, err)
		}
	}

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
			t.Errorf("chapter-0.html unexpectedly contains %q: site must be fully self-contained", unwanted)
		}
	}
}

// index.html sits alongside the pages it fronts, so Pages serving docs/ as
// the site root resolves a bare / to it and every URL below is a clean
// /chapter-N.html. Sharing a directory with its assets, it references them
// with no prefix, the same way every chapter page does.
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

// docs/ is disposable: deleting it and regenerating has to reproduce every
// published file. .nojekyll was hand-made once and vanished on the first
// clean rebuild, which would have let Jekyll rewrite the output.
func TestRunWritesNoJekyllSoOutputSurvivesACleanRebuild(t *testing.T) {
	outDir, _ := testRun(t)

	if _, err := os.Stat(filepath.Join(outDir, ".nojekyll")); err != nil {
		t.Errorf(".nojekyll was not generated: %v", err)
	}
}

// The illustrations are the only files under docs/ that nothing can rebuild.
// A build that finds them gone has to stop, because carrying on would publish
// twenty-one placeholder panels over a working site and nothing would notice
// until someone opened it.
func TestRunRefusesToBuildWithoutTheImages(t *testing.T) {
	outDir := filepath.Join(t.TempDir(), "docs")
	indexPath := filepath.Join(outDir, "index.html")

	err := run(filepath.Join("..", "..", "templates"), filepath.Join("..", "..", "assets"), outDir, indexPath)
	if err == nil {
		t.Fatal("expected an error when docs/images is missing, got nil")
	}
	if !strings.Contains(err.Error(), "images") {
		t.Errorf("error should name the missing directory, got: %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(outDir, "chapter-0.html")); statErr == nil {
		t.Error("a page was written despite the images being missing")
	}
}

// Regenerating must never disturb the illustrations, however many times it runs.
func TestRunLeavesTheImagesUntouched(t *testing.T) {
	outDir, indexPath := testRun(t)

	before := map[string][]byte{}
	entries, err := os.ReadDir(filepath.Join(outDir, "images"))
	if err != nil {
		t.Fatalf("reading images: %v", err)
	}
	for _, e := range entries {
		b, readErr := os.ReadFile(filepath.Join(outDir, "images", e.Name()))
		if readErr != nil {
			t.Fatalf("reading %s: %v", e.Name(), readErr)
		}
		before[e.Name()] = b
	}

	for i := 0; i < 3; i++ {
		if err := run(filepath.Join("..", "..", "templates"), filepath.Join("..", "..", "assets"), outDir, indexPath); err != nil {
			t.Fatalf("rebuild %d: %v", i, err)
		}
	}

	after, err := os.ReadDir(filepath.Join(outDir, "images"))
	if err != nil {
		t.Fatalf("reading images after rebuilds: %v", err)
	}
	if len(after) != len(before) {
		t.Fatalf("image count changed: %d before, %d after", len(before), len(after))
	}
	for name, want := range before {
		got, readErr := os.ReadFile(filepath.Join(outDir, "images", name))
		if readErr != nil {
			t.Errorf("%s went missing across rebuilds: %v", name, readErr)
			continue
		}
		if string(got) != string(want) {
			t.Errorf("%s was modified by a rebuild", name)
		}
	}
}

// A prompt that reaches an assistant without naming a language is worse than
// useless: the assistant picks one itself, and can pick differently on
// different chapters, so a later chapter's code will not build on an earlier
// one's. Every page must carry a language before any script runs.
func TestEveryPageBakesInALanguage(t *testing.T) {
	outDir, indexPath := testRun(t)

	pages, err := filepath.Glob(filepath.Join(outDir, "chapter-*.html"))
	if err != nil {
		t.Fatalf("globbing pages: %v", err)
	}
	pages = append(pages, indexPath)

	want := "Build this in " + content.LanguageByID(content.DefaultLanguage).Name + ","
	for _, page := range pages {
		b, readErr := os.ReadFile(page)
		if readErr != nil {
			t.Fatalf("reading %s: %v", page, readErr)
		}
		if !strings.Contains(string(b), want) {
			t.Errorf("%s does not name a language", filepath.Base(page))
		}
	}
}

// Chapter 0 is the only chapter with nothing before it, so it is the only one
// that must not tell the assistant to continue an existing codebase.
func TestOnlyLaterChaptersAskToContinueTheCodebase(t *testing.T) {
	outDir, _ := testRun(t)

	const carryOn = "Continue the same codebase"
	for _, c := range content.All {
		b, err := os.ReadFile(filepath.Join(outDir, "chapter-"+strconv.Itoa(c.Number)+".html"))
		if err != nil {
			t.Fatalf("reading chapter %d: %v", c.Number, err)
		}
		has := strings.Contains(string(b), carryOn)
		if c.Number == 0 && has {
			t.Error("chapter 0 asks the assistant to continue a codebase that does not exist yet")
		}
		if c.Number > 0 && !has {
			t.Errorf("chapter %d does not ask the assistant to continue the existing codebase", c.Number)
		}
	}
}

// The prompts describe what to build and must never name a language, or the
// dropdown would contradict the text below it.
func TestPromptsNameNoLanguage(t *testing.T) {
	banned := []string{"go run", "net.Listen", "goroutine", "buffered channel",
		"fmt.Println", "Go standard library", "Go variable", "a mutex"}
	for _, c := range content.All {
		for _, b := range banned {
			if strings.Contains(c.BuildIt.Prompt, b) {
				t.Errorf("chapter %d prompt names %q, which ties it to one language", c.Number, b)
			}
		}
	}
}

// The language applies to every chapter, so exactly one chapter offers the
// choice. A picker on chapter 12 would invite a reader to switch halfway
// through and leave eleven chapters of code in another language.
func TestOnlyOneChapterOffersTheLanguageChoice(t *testing.T) {
	outDir, indexPath := testRun(t)

	withPicker := []string{}
	for _, c := range content.All {
		page := filepath.Join(outDir, "chapter-"+strconv.Itoa(c.Number)+".html")
		b, err := os.ReadFile(page)
		if err != nil {
			t.Fatalf("reading chapter %d: %v", c.Number, err)
		}
		if strings.Contains(string(b), "data-language-select") {
			withPicker = append(withPicker, "chapter-"+strconv.Itoa(c.Number))
		}
	}
	if len(withPicker) != 1 {
		t.Errorf("expected exactly one chapter to offer the picker, got %v", withPicker)
	}
	if len(withPicker) == 1 && withPicker[0] != "chapter-0" {
		t.Errorf("the picker belongs on chapter 0, found it on %s", withPicker[0])
	}

	// The entry point is chapter 0 rendered again, so it carries the picker too.
	b, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("reading index: %v", err)
	}
	if !strings.Contains(string(b), "data-language-select") {
		t.Error("the entry point is chapter 0 and should offer the picker")
	}
}

// The choice is made on chapter 0, which owns the picker and can resolve an id
// to a name itself. It stores the name, so no other page needs a copy of the
// language list. Carrying one on all twenty-two would be dead weight that only
// two of them could ever use.
func TestOnlyTheChoosingPageCarriesTheLanguageList(t *testing.T) {
	outDir, indexPath := testRun(t)

	pages, err := filepath.Glob(filepath.Join(outDir, "chapter-*.html"))
	if err != nil {
		t.Fatalf("globbing: %v", err)
	}
	for _, page := range append(pages, indexPath) {
		b, readErr := os.ReadFile(page)
		if readErr != nil {
			t.Fatalf("reading %s: %v", page, readErr)
		}
		body := string(b)
		name := filepath.Base(page)
		hasOptions := strings.Contains(body, "<option value=")
		wantsOptions := name == "chapter-0.html" || name == "index.html"
		if hasOptions != wantsOptions {
			t.Errorf("%s: carries language options = %v, want %v", name, hasOptions, wantsOptions)
		}
		if strings.Contains(body, "data-languages") {
			t.Errorf("%s still embeds a language list", name)
		}
	}
}
