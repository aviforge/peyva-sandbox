package main

import (
	"os"
	"path/filepath"
	"regexp"
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
		images := []string{c.HeroImage}
		if c.Aside != nil {
			images = append(images, c.Aside.HeroImage)
		}
		for _, img := range images {
			name := filepath.Base(filepath.FromSlash(img))
			if err := os.WriteFile(filepath.Join(imagesDir, name), []byte("fake-webp"), 0o644); err != nil {
				t.Fatalf("writing fixture hero %s: %v", name, err)
			}
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

	for _, want := range []string{"What Are We Building?", "1. What", "2. Why",
		"3. How", "Mark chapter as complete", "Build the Vault", "Inside One Computer", "Zero-shot prompting", "Documented in The Prompt Report"} {
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
		// The version query is a content digest, so match the path only.
		`href="styles.css?v=`,
		`src="app.js?v=`,
		`href="chapter-0.html"`,
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

// A page that asks for code names the language it wants, so a prompt never
// reaches an assistant that has to pick one itself.
//
// Chapter 9 is the exception and names none: it estimates capacity and writes
// no code, so there is nothing for a language to apply to.
func TestEveryPageBakesInALanguage(t *testing.T) {
	outDir, indexPath := testRun(t)

	buildsCode := map[int]bool{}
	for _, c := range content.All {
		for _, prompt := range c.BuildIt.Prompts {
			if !prompt.Thinking {
				buildsCode[c.Number] = true
			}
		}
	}

	pages, err := filepath.Glob(filepath.Join(outDir, "chapter-*.html"))
	if err != nil {
		t.Fatalf("globbing: %v", err)
	}
	want := content.LanguageByID(content.DefaultLanguage).Name
	_ = indexPath
	for _, page := range pages {
		b, readErr := os.ReadFile(page)
		if readErr != nil {
			t.Fatalf("reading %s: %v", page, readErr)
		}
		name := filepath.Base(page)
		number := 0
		if name != "index.html" {
			number, _ = strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(name, "chapter-"), ".html"))
		}
		if !buildsCode[number] {
			continue
		}
		if !strings.Contains(string(b), want) {
			t.Errorf("%s does not name a language", name)
		}
	}
}

// No prompt sends the assistant off to read what earlier chapters built.
//
// The preamble used to, on every chapter after the first. It cost thousands of
// tokens per chapter to rediscover what one opening sentence says, and on a
// small context window the later chapters did not fit at all. Each prompt now
// states the state it starts from.
func TestNoPromptTellsTheAssistantToReadTheCodebase(t *testing.T) {
	outDir, _ := testRun(t)

	banned := []string{
		"Continue the codebase",
		"read the codebase",
		"read the existing code",
		"from earlier chapters",
		"in previous chapters",
	}
	for _, c := range content.All {
		b, err := os.ReadFile(filepath.Join(outDir, "chapter-"+strconv.Itoa(c.Number)+".html"))
		if err != nil {
			t.Fatalf("reading chapter %d: %v", c.Number, err)
		}
		for _, phrase := range banned {
			if strings.Contains(string(b), phrase) {
				t.Errorf("chapter %d says %q, which sends the assistant through the whole project",
					c.Number, phrase)
			}
		}
	}
}

// Every prompt opens by saying what it starts from, because nothing else does
// any more. A prompt that opens with an instruction and no starting state
// leaves the assistant to guess what already exists, which is how a chapter
// rebuilds a component that was already there.
//
// This only measures length. Whether an opening describes the right state is a
// judgement, and an earlier version of this test tried to make it by looking
// for words like "currently" and "already": it failed prompts that plainly did
// describe a situation and passed any sentence containing "The". A wrong test
// gets worked around rather than obeyed. An opening too short to hold a
// situation definitely does not hold one, and that is worth pinning.
func TestEveryPromptOpensWithEnoughToStandOn(t *testing.T) {
	const minWords = 12
	for _, c := range content.All {
		for _, prompt := range c.BuildIt.Prompts {
			opening := prompt.Text
			if i := strings.Index(opening, "\n\n"); i > 0 {
				opening = opening[:i]
			}
			if n := len(strings.Fields(opening)); n < minWords {
				t.Errorf("chapter %d prompt %q opens with %d words, too few to say what it starts from: %q",
					c.Number, prompt.Label, n, opening)
			}
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
			if strings.Contains(allPromptText(c), b) {
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

	// The entry point is the landing page, which offers no choices at all.
	b, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("reading index: %v", err)
	}
	if strings.Contains(string(b), "data-language-select") {
		t.Error("the landing page offers the picker, which belongs to chapter 0")
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
	_ = indexPath
	for _, page := range pages {
		b, readErr := os.ReadFile(page)
		if readErr != nil {
			t.Fatalf("reading %s: %v", page, readErr)
		}
		body := string(b)
		name := filepath.Base(page)
		// Scoped to the language select. Chapter 10 offers an operating system
		// from a second select in the same section, and counting options across
		// the whole page reads those as languages.
		hasOptions := strings.Contains(languageSelect.FindString(body), "<option value=")
		wantsOptions := name == "chapter-0.html"
		if hasOptions != wantsOptions {
			t.Errorf("%s: carries language options = %v, want %v", name, hasOptions, wantsOptions)
		}
		if strings.Contains(body, "data-languages") {
			t.Errorf("%s still embeds a language list", name)
		}
	}
}

// languageSelect isolates the language control, so a page that also offers an
// operating system does not have its options counted as languages.
var languageSelect = regexp.MustCompile(`(?s)<select[^>]*data-language-select.*?</select>`)

// The standing rules travel with every prompt that asks for code, because a
// prompt is copied out of the page on its own and nothing else tells the
// assistant where code goes or how to hold money.
//
// A turn that produces an answer rather than a change carries none of it. Four
// lines about rounding money are noise on a turn asking for an analogy, and
// chapter 9 is nothing but such turns: it estimates capacity and writes no
// code, so its page states no build rules at all.
func TestEveryPromptCarriesTheStandingRules(t *testing.T) {
	outDir, indexPath := testRun(t)

	rules := map[string]string{
		"folder layout": "peyva/&lt;component&gt;/, one folder per component",
		"exactness":     "never floating point",
		"precision":     "two decimal places",
	}

	buildsCode := map[int]bool{}
	for _, c := range content.All {
		for _, prompt := range c.BuildIt.Prompts {
			if !prompt.Thinking && !prompt.Portal {
				buildsCode[c.Number] = true
			}
		}
	}

	pages, err := filepath.Glob(filepath.Join(outDir, "chapter-*.html"))
	if err != nil {
		t.Fatalf("globbing: %v", err)
	}
	_ = indexPath
	for _, page := range pages {
		b, readErr := os.ReadFile(page)
		if readErr != nil {
			t.Fatalf("reading %s: %v", page, readErr)
		}
		name := filepath.Base(page)

		number := 0
		if name != "index.html" {
			number, _ = strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(name, "chapter-"), ".html"))
		}

		for what, rule := range rules {
			has := strings.Contains(string(b), rule)
			if buildsCode[number] && !has {
				t.Errorf("%s: asks for code but does not state the %s rule", name, what)
			}
			if !buildsCode[number] && has {
				t.Errorf("%s: states the %s rule but asks for no code", name, what)
			}
		}
	}
}

// The preamble is written by the generator and the script only swaps the
// language word inside it. If the script ever rebuilds the sentence again, the
// rules exist in two places and the copy that is not the generator will fall
// behind, which is exactly what happened the first time.
func TestScriptDoesNotRebuildThePreamble(t *testing.T) {
	b, err := os.ReadFile(filepath.Join("..", "..", "assets", "app.js"))
	if err != nil {
		t.Fatalf("reading app.js: %v", err)
	}
	for _, fragment := range []string{"Build this in", "standard library", "floating point", "peyva/"} {
		if strings.Contains(string(b), fragment) {
			t.Errorf("app.js contains %q: the preamble belongs to the generator alone", fragment)
		}
	}
}

// allPromptText is every prompt of a chapter joined.
func allPromptText(c content.ChapterContent) string {
	var b strings.Builder
	for _, p := range c.BuildIt.Prompts {
		b.WriteString(p.Text)
		b.WriteString("\n")
	}
	return b.String()
}

// The find box hides a chapter by setting the hidden property on its li. That
// property only hides anything if nothing in the stylesheet says otherwise, and
// .chapter-list li sets display: block. An author rule beats the browser's own
// [hidden] rule whatever its specificity, so without an explicit override the
// box filtered nothing and every chapter stayed on screen.
func TestStylesheetHidesFilteredChapters(t *testing.T) {
	b, err := os.ReadFile(filepath.Join("..", "..", "assets", "styles.css"))
	if err != nil {
		t.Fatalf("reading styles.css: %v", err)
	}
	css := string(b)
	if !strings.Contains(css, ".chapter-list li[hidden]") {
		t.Error("styles.css does not hide a filtered chapter: .chapter-list li sets display, so hidden alone does nothing")
	}
	// Whichever selector sets display on a list item must come before the
	// override, or the override loses to the later rule of equal weight.
	if i, j := strings.Index(css, ".chapter-list li,"), strings.Index(css, ".chapter-list li[hidden]"); i >= 0 && j >= 0 && j < i {
		t.Error("the [hidden] override is declared before .chapter-list li, so the display rule wins")
	}
}
