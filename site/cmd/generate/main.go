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

	defaultLanguage := content.LanguageByID(content.DefaultLanguage)
	pageData := func(chapter content.ChapterContent, prefix string) PageData {
		return PageData{
			Chapter:       chapter,
			HeroAvailable: fileExists(filepath.Join(outDir, filepath.FromSlash(chapter.HeroImage))),
			Roadmap:       roadmap,
			Labs:          content.Labs,
			AssetPrefix:   prefix,
			LanguageLine:  languageLine(defaultLanguage, chapter.Number),
			UILine:        uiLine(defaultLanguage),
			Setup:         setupFor(chapter.Number),
			Languages:     content.Languages,

			LanguageIsChosenHere: chapter.Number == languagePickerChapter,
			LanguageName:         defaultLanguage.Name,
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

// languagePickerChapter is the one chapter that offers the choice. It is the
// first, because the choice has to be made before any code exists.
const languagePickerChapter = 0

// languageLine is the instruction that travels with every prompt. The prompts
// describe what to build and never name a language, so without this the
// assistant chooses one itself, and can choose differently on each chapter.
//
// Chapter 0 is the only one with nothing before it, so it is the only one that
// does not ask the assistant to continue an existing codebase.
func languageLine(l content.Language, chapterNumber int) template.HTML {
	// The language name is wrapped so the script can swap that word and nothing
	// else. Rebuilding the sentence in JavaScript as well would mean two copies
	// of these rules, and they drift: the first version of this already had the
	// script writing an older wording that silently dropped two of them.
	name := `<span data-language-name>` + template.HTMLEscapeString(l.Name) + `</span>`

	line := "Build this in " + name + ", standard library only."
	if chapterNumber > 0 {
		line += "\nContinue the codebase from earlier chapters."
	}
	// Every word here is paid twenty-one times, once per prompt, so the rules
	// are stated as tightly as they can be without losing what they mean. The
	// minor-units clause stays because five of the twelve languages have no
	// decimal type in their standard library, and without it the assistant
	// reaches for a dependency the first line just ruled out.
	line += "\nCode in peyva/&lt;component&gt;/, one folder per component."
	line += "\nMoney: exact decimal or integer minor units, never floating point, two decimal places."
	// The constraints stay inline because they are cheap and always present.
	// The goal and the invariants do not fit in a preamble, and an assistant
	// that has to infer them reads the whole codebase instead, which costs far
	// more than the one small file this points at.
	line += "\nThe goal and invariants are in peyva/goal.md."
	return template.HTML(line)
}

// uiLine is the preamble for the portal prompt. It repeats the language and the
// layout because the portal is a component like any other, and drops the money
// rule because a page renders a balance rather than holding one.
//
// The no-dependency line is the whole reason the portal can exist in a book
// that allows twelve backend languages: a page built from plain HTML and CSS
// needs no toolchain, so it is the same page whichever language serves it.
func uiLine(l content.Language) template.HTML {
	name := `<span data-language-name>` + template.HTMLEscapeString(l.Name) + `</span>`
	line := "The portal is plain HTML and CSS. No framework, no build step, no dependencies."
	line += "\nIt lives in peyva/portal/. Anything it needs from the server is " + name + ", standard library only."
	// Without the first half the assistant builds an operator's console: a table
	// of every account, with everyone's balance on it. That is a different
	// product, and by chapter 18 a privacy bug the sign-in cannot undo.
	//
	// The switcher is what makes the first half survivable. One customer at a
	// time on one laptop means a reader sends alice's money and never sees it
	// land, and the arriving half of a payment is half of what this book is
	// about.
	line += "\nOne customer's wallet at a time, never everyone's at once. A switcher at the top says whose."
	line += "\nNew work joins the menu."
	line += "\nThe goal and invariants are in peyva/goal.md."
	return template.HTML(line)
}

// setupFor returns the files to save only for the chapter that starts the
// project. Every chapter after it points at peyva/goal.md, which the reader
// saved here, and inherits the rules from the assistant's own instruction file.
func setupFor(chapterNumber int) []content.SetupFile {
	if chapterNumber != languagePickerChapter {
		return nil
	}
	return content.SetupFiles
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
