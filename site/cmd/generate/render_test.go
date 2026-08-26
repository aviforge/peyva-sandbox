package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"site/content"
)

// sections every chapter page must render.
var requiredSections = []string{
	"1. Intuition",
	"2. Concepts",
	"3. Under the Hood",
	"4. Build It",
	"5. Break It",
}

// pairedTags are counted open against closed. The character class matters:
// "<li" also matches "<link", and "<figure" also matches "<figcaption".
var pairedTags = []string{"ul", "li", "section", "figure", "select", "pre"}

var (
	unrendered  = regexp.MustCompile(`\{\{[^}]{0,60}\}\}`)
	heroImg     = regexp.MustCompile(`<img src="[^"]*images/chapter-\d+\.webp"`)
	sidebarLink = regexp.MustCompile(`data-slug="chapter-\d+"`)
	optionTag   = regexp.MustCompile(`<option value="`)
	sidebarNav  = regexp.MustCompile(`(?s)<nav class="sidebar">.*?</nav>`)
	buildIt     = regexp.MustCompile(`(?s)<section class="block build-it">.*?</section>`)
)

// renderedPage is one generated file plus what the test needs to know about it.
type renderedPage struct {
	name  string
	body  string
	isCh0 bool
}

func renderedPages(t *testing.T) []renderedPage {
	t.Helper()
	outDir, indexPath := testRun(t)

	paths, err := filepath.Glob(filepath.Join(outDir, "chapter-*.html"))
	if err != nil {
		t.Fatalf("globbing pages: %v", err)
	}
	paths = append(paths, indexPath)

	pages := make([]renderedPage, 0, len(paths))
	for _, p := range paths {
		b, readErr := os.ReadFile(p)
		if readErr != nil {
			t.Fatalf("reading %s: %v", p, readErr)
		}
		name := filepath.Base(p)
		// index.html is chapter 0 rendered a second time, so it is chapter 0
		// for every purpose here.
		pages = append(pages, renderedPage{
			name:  name,
			body:  string(b),
			isCh0: name == "chapter-0.html" || name == "index.html",
		})
	}
	if len(pages) != len(content.All)+1 {
		t.Fatalf("rendered %d pages, want %d", len(pages), len(content.All)+1)
	}
	return pages
}

// A page that renders without error can still be wrong in ways nothing else
// notices: a template action left unexecuted, a section silently dropped, a
// hero replaced by its placeholder. None of those fail a build, and all of
// them reach the reader.
func TestEveryPageRendersCompletely(t *testing.T) {
	for _, p := range renderedPages(t) {
		if m := unrendered.FindString(p.body); m != "" {
			t.Errorf("%s: template action left unrendered: %s", p.name, m)
		}
		if strings.Contains(p.body, "Image pending") {
			t.Errorf("%s: hero fell back to the placeholder", p.name)
		}
		if !heroImg.MatchString(p.body) {
			t.Errorf("%s: no hero image element", p.name)
		}
		for _, s := range requiredSections {
			if !strings.Contains(p.body, s) {
				t.Errorf("%s: missing section %q", p.name, s)
			}
		}
		if n := len(sidebarLink.FindAllString(p.body, -1)); n != len(content.All) {
			t.Errorf("%s: sidebar lists %d chapters, want %d", p.name, n, len(content.All))
		}
		if !strings.Contains(p.body, "| Peyva Sandbox</title>") {
			t.Errorf("%s: title is missing the site name", p.name)
		}
		if strings.Contains(p.body, "—") {
			t.Errorf("%s: contains an em dash", p.name)
		}
	}
}

// The tags are nested by hand in the template, so an edit that opens one and
// forgets to close it produces a page browsers will silently reflow into
// something else.
func TestEveryPageHasBalancedTags(t *testing.T) {
	for _, p := range renderedPages(t) {
		for _, tag := range pairedTags {
			open := regexp.MustCompile(`<`+tag+`[ >]`).FindAllString(p.body, -1)
			closed := strings.Count(p.body, "</"+tag+">")
			if len(open) != closed {
				t.Errorf("%s: <%s> opened %d times, closed %d", p.name, tag, len(open), closed)
			}
		}
	}
}

// The language control belongs beside the prompt it governs, and the choice
// belongs to chapter 0 alone. Both have moved once already; this pins where
// they ended up.
func TestLanguageControlSitsWithThePromptItGoverns(t *testing.T) {
	for _, p := range renderedPages(t) {
		nav := sidebarNav.FindString(p.body)
		if nav != "" && strings.Contains(nav, "data-language-select") {
			t.Errorf("%s: the language picker is back in the sidebar", p.name)
		}

		section := buildIt.FindString(p.body)
		if section == "" {
			t.Errorf("%s: no Build It section", p.name)
			continue
		}

		if p.isCh0 {
			if !strings.Contains(section, "data-language-select") {
				t.Errorf("%s: chapter 0 should offer the picker in Build It", p.name)
			}
			if n := len(optionTag.FindAllString(p.body, -1)); n != len(content.Languages) {
				t.Errorf("%s: offers %d languages, want %d", p.name, n, len(content.Languages))
			}
			if strings.Contains(p.body, "language-locked") {
				t.Errorf("%s: chapter 0 makes the choice, so it should not show it as locked", p.name)
			}
			continue
		}

		if strings.Contains(p.body, "data-language-select") {
			t.Errorf("%s: only chapter 0 may offer the picker", p.name)
		}
		if !strings.Contains(section, "language-locked") {
			t.Errorf("%s: does not show which language is in force", p.name)
		}
		if !strings.Contains(section, "chapter-0.html") {
			t.Errorf("%s: shows the language with no way back to change it", p.name)
		}
	}
}

// The prompts are long enough that selecting them by hand is the friction this
// button exists to remove, and one per page is the whole contract.
func TestEveryPageHasExactlyOneCopyButton(t *testing.T) {
	for _, p := range renderedPages(t) {
		if n := strings.Count(p.body, "data-copy-prompt"); n != 1 {
			t.Errorf("%s: %d copy buttons, want 1", p.name, n)
		}
	}
}

// Nothing may be fetched at read time. The site has to work from a clone with
// no network, which is the promise the README makes.
func TestEveryPageIsSelfContained(t *testing.T) {
	for _, p := range renderedPages(t) {
		for _, unwanted := range []string{"http://", "fetch(", "<iframe"} {
			if strings.Contains(p.body, unwanted) {
				t.Errorf("%s: contains %q, so it is not self-contained", p.name, unwanted)
			}
		}
	}
}
