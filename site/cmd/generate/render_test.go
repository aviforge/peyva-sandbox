package main

import (
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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
			if n := len(optionTag.FindAllString(section, -1)); n != len(content.Languages) {
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

// The operating system is chosen once, in setup, and reaches the prompts that
// name one. A chapter whose prompt says {os} and renders it literally hands the
// reader a script for an operating system called "{os}".
func TestSystemChoiceReachesThePromptsThatNeedIt(t *testing.T) {
	byNumber := map[int]content.ChapterContent{}
	for _, c := range content.All {
		byNumber[c.Number] = c
	}

	usesIt := 0
	for _, c := range content.All {
		if namesASystem(c) {
			usesIt++
		}
	}
	if usesIt == 0 {
		t.Fatal("no chapter uses the placeholder, so the picker changes nothing")
	}

	for _, p := range renderedPages(t) {
		if strings.Contains(p.body, osPlaceholder) {
			t.Errorf("%s: renders %s literally instead of an operating system", p.name, osPlaceholder)
		}

		num := 0
		if p.name != "index.html" {
			num, _ = strconv.Atoi(regexp.MustCompile(`chapter-(\d+)`).FindStringSubmatch(p.name)[1])
		}
		// A chapter whose prompt says {os} renders the name in a span. The
		// chapter that owns the picker shows the select instead, which is the
		// choice itself rather than a report of it.
		wantsName := namesASystem(byNumber[num])
		hasName := strings.Contains(p.body, "data-os-name")
		if wantsName != hasName {
			t.Errorf("%s: names an operating system = %v, want %v", p.name, hasName, wantsName)
		}

		// The picker belongs to the first chapter that needs a system, not to
		// setup: nothing before it cares which one the reader is on. Only that
		// chapter offers it, for the same reason only chapter 0 offers the
		// language. A PowerShell runner in chapter 10 with bash commands to
		// operate it in chapter 19 is worse than either alone.
		wantPicker := num == systemPickerChapter() && p.name != "index.html"
		hasPicker := strings.Contains(p.body, "data-system-select")
		if wantPicker != hasPicker {
			t.Errorf("%s: offers the system picker = %v, want %v", p.name, hasPicker, wantPicker)
		}
		if wantPicker {
			// Every system's script is present, and exactly one is shown.
			for _, r := range content.RunnerScripts {
				if !strings.Contains(p.body, `data-runner-for="`+r.SystemID+`"`) {
					t.Errorf("%s: carries no runner script for %s", p.name, r.SystemID)
				}
			}
			for _, sys := range content.Systems {
				if !strings.Contains(p.body, `>`+sys.Name+`</option>`) {
					t.Errorf("%s: the picker does not offer %s", p.name, sys.Name)
				}
			}
		}
		// Every other chapter that names a system has to say where to change it.
		if wantsName && !wantPicker && !strings.Contains(p.body, "chosen in") {
			t.Errorf("%s: names a system with no way back to the choice", p.name)
		}
	}
}

// One copy button per prompt. A chapter that grows the portal has two prompts,
// because a reader copies one thing at a time and would otherwise get a
// component and a page in a single paste with no place to stop between them.
func TestEachPromptHasItsOwnCopyButton(t *testing.T) {
	byNumber := map[int]content.ChapterContent{}
	for _, c := range content.All {
		byNumber[c.Number] = c
	}

	for _, p := range renderedPages(t) {
		want := 1
		num := 0
		if p.name != "index.html" {
			num, _ = strconv.Atoi(regexp.MustCompile(`chapter-(\d+)`).FindStringSubmatch(p.name)[1])
		}
		if byNumber[num].BuildIt.UIPrompt != "" {
			want++
		}
		// Chapter 0 also carries the files the reader saves before starting.
		if num == 0 {
			want += len(content.SetupFiles)
		}
		// The runner chapter carries one script per operating system. Only one
		// is visible at a time, but all of them are in the page so switching
		// needs no network.
		if num == content.RunnerChapter {
			want += len(content.RunnerScripts)
		}
		if got := strings.Count(p.body, "data-copy-prompt"); got != want {
			t.Errorf("%s: %d copy buttons, want %d", p.name, got, want)
		}
		if got := strings.Count(p.body, `<pre class="prompt">`); got != want {
			t.Errorf("%s: %d prompt blocks, want %d", p.name, got, want)
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

// The portal prompt is copied on its own, so it carries its own rules. The
// money rule is deliberately absent: a page renders a balance rather than
// holding one, and the no-dependency line is what keeps the portal buildable
// whichever of the twelve backend languages the reader picked.
func TestPortalPromptCarriesItsOwnRules(t *testing.T) {
	byNumber := map[int]content.ChapterContent{}
	for _, c := range content.All {
		byNumber[c.Number] = c
	}

	for _, p := range renderedPages(t) {
		num := 0
		if p.name != "index.html" {
			num, _ = strconv.Atoi(regexp.MustCompile(`chapter-(\d+)`).FindStringSubmatch(p.name)[1])
		}
		if byNumber[num].BuildIt.UIPrompt == "" {
			continue
		}
		for _, rule := range []string{
			"plain HTML and CSS",
			"No framework, no build step, no dependencies",
			"peyva/portal/",
			"peyva/goal.md",
			"never everyone's at once",
			"switcher at the top says whose",
		} {
			if !strings.Contains(p.body, rule) {
				t.Errorf("%s: portal prompt is missing %q", p.name, rule)
			}
		}
	}
}

// A portal prompt must not name a language either. The page is the same page
// whichever language serves it, which is the only reason one portal can exist
// across twelve of them.
func TestPortalPromptsNameNoLanguage(t *testing.T) {
	banned := []string{"go run", "goroutine", "fmt.Println", "npm", "React", "Lit", "node_modules"}
	for _, c := range content.All {
		for _, b := range banned {
			if strings.Contains(c.BuildIt.UIPrompt, b) {
				t.Errorf("chapter %d portal prompt names %q", c.Number, b)
			}
		}
	}
}

// A chapter that adds portal work says what it adds. An empty intro beside a
// filled prompt means the section renders a bare instruction with no context.
func TestPortalPromptsComeWithAnIntro(t *testing.T) {
	for _, c := range content.All {
		if (c.BuildIt.UIPrompt == "") != (c.BuildIt.UIIntro == "") {
			t.Errorf("chapter %d: UIPrompt and UIIntro must both be set or both be empty", c.Number)
		}
		if c.BuildIt.UIPrompt != "" && !strings.Contains(c.BuildIt.UIPrompt, "Done when") {
			t.Errorf("chapter %d: portal prompt has no 'Done when' line", c.Number)
		}
	}
}

// The setup files are rendered into chapter 0 rather than kept as files in this
// repository, so a reader following the published site never has to leave it to
// fetch one. The spec carries the invariants, which no prompt states and
// nothing else would tell an assistant.
func TestChapterZeroCarriesTheSetupFiles(t *testing.T) {
	for _, p := range renderedPages(t) {
		// Match the section element, not the word. "Setup" appears in prose on
		// other pages.
		hasSetup := strings.Contains(p.body, `<section class="block setup">`)
		if p.isCh0 != hasSetup {
			t.Errorf("%s: renders the setup section = %v, want %v", p.name, hasSetup, p.isCh0)
		}
		if !p.isCh0 {
			continue
		}
		for _, line := range []string{
			"Money is never created and never lost",
			"No balance goes negative",
			"No payment is applied twice",
			"Only the Vault changes a balance",
			"peyva/goal.md",
		} {
			if !strings.Contains(p.body, line) {
				t.Errorf("chapter 0 spec is missing %q", line)
			}
		}
		for _, line := range []string{
			"Read peyva/goal.md before your first change",
			"nothing it did not ask for",
		} {
			if !strings.Contains(p.body, line) {
				t.Errorf("chapter 0 agent rules are missing %q", line)
			}
		}
	}
}

// The spec says what to build; the agent file says how to work. A rule that
// appears in both is a rule with two places to change, and the copy that goes
// stale is the one nobody notices.
func TestSetupFilesDoNotRestateEachOther(t *testing.T) {
	for _, claim := range []string{
		"Standard library only",
		"one folder per component",
		"Never floating point",
		"plain HTML and CSS",
		"Structure is earned",
	} {
		if !strings.Contains(content.GoalSpec, claim) {
			t.Errorf("the spec no longer states %q, so nothing does", claim)
		}
		if strings.Contains(content.AgentRules, claim) {
			t.Errorf("the agent file restates %q, which belongs to the spec", claim)
		}
	}
}

// Every setup file is copyable and every one that has no agreed filename says
// which name each assistant reads. A reader who saves the rules as CLAUDE.md
// while using Codex has a file nothing will ever open.
func TestSetupFilesAreCopyableAndNamed(t *testing.T) {
	for _, f := range content.SetupFiles {
		if strings.TrimSpace(f.Content) == "" {
			t.Errorf("setup file %q has no content", f.Path)
		}
		if strings.TrimSpace(f.Purpose) == "" {
			t.Errorf("setup file %q does not say what it is for", f.Path)
		}
	}
	for _, p := range renderedPages(t) {
		if !p.isCh0 {
			continue
		}
		for _, a := range content.AgentFiles {
			if !strings.Contains(p.body, a.Tool) {
				t.Errorf("%s: does not name %q", p.name, a.Tool)
			}
			if !strings.Contains(p.body, template.HTMLEscapeString(a.Path)) {
				t.Errorf("%s: does not give the path for %q", p.name, a.Tool)
			}
		}
	}
}
