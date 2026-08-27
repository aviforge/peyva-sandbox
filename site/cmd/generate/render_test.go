package main

import (
	"html"
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
	"1. What",
	"2. Why",
	"3. How",
}

// pairedTags are counted open against closed. The character class matters:
// "<li" also matches "<link", and "<figure" also matches "<figcaption".
var pairedTags = []string{"ul", "li", "section", "figure", "select", "pre"}

var (
	unrendered   = regexp.MustCompile(`\{\{[^}]{0,60}\}\}`)
	heroImg      = regexp.MustCompile(`<img src="[^"]*images/chapter-[\w-]+\.webp"`)
	sidebarLink  = regexp.MustCompile(`data-slug="chapter-\d+"`)
	optionTag    = regexp.MustCompile(`<option value="`)
	sidebarNav   = regexp.MustCompile(`(?s)<nav class="sidebar">.*?</nav>`)
	buildIt      = regexp.MustCompile(`(?s)<section class="block build-it">.*?</section>`)
	setupSection = regexp.MustCompile(`(?s)<section class="block setup".*?</section>`)
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
	_ = indexPath // the landing page is checked by TestLandingPage, not as a chapter

	pages := make([]renderedPage, 0, len(paths))
	for _, p := range paths {
		b, readErr := os.ReadFile(p)
		if readErr != nil {
			t.Fatalf("reading %s: %v", p, readErr)
		}
		name := filepath.Base(p)
		pages = append(pages, renderedPage{
			name:  name,
			body:  string(b),
			isCh0: name == "chapter-0.html",
		})
	}
	if len(pages) != len(content.All) {
		t.Fatalf("rendered %d pages, want %d", len(pages), len(content.All))
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
		// Every page says the two choices have to be made before its prompts
		// can be copied, because a reader can land on any chapter from a link.
		if !strings.Contains(p.body, "data-choice-gate") {
			t.Errorf("%s: does not tell the reader to choose a language and a system", p.name)
		}
		for _, s := range requiredSections {
			if !strings.Contains(p.body, s) {
				t.Errorf("%s: missing section %q", p.name, s)
			}
		}
		if n := len(sidebarLink.FindAllString(p.body, -1)); n != len(content.All) {
			t.Errorf("%s: sidebar lists %d chapters, want %d", p.name, n, len(content.All))
		}
		if !strings.Contains(p.body, `href="index.html" data-slug="index">Home</a>`) {
			t.Errorf("%s: sidebar has no Home link", p.name)
		}
		if !strings.Contains(p.body, "data-chapter-search") {
			t.Errorf("%s: sidebar has no find box", p.name)
		}
		if !strings.Contains(p.body, `data-terms="`) {
			t.Errorf("%s: sidebar entries carry no terms to find by", p.name)
		}
		if !strings.Contains(p.body, "| Peyva Sandbox</title>") {
			t.Errorf("%s: title is missing the site name", p.name)
		}
		if strings.Contains(p.body, "—") {
			t.Errorf("%s: contains an em dash", p.name)
		}
	}
}

// Nothing on a page may name a section the site no longer has.
//
// Three sections were removed, and every check written at the time searched for
// their exact headings. The sidebar tagline said "Build it. Break it. Understand
// why it works." on all twenty-two pages and no sweep saw it, because it wrote
// Break it and the searches looked for Break It. This one is case insensitive.
func TestNoPageNamesARemovedSection(t *testing.T) {
	removed := []string{"break it", "under the hood", "intuition"}
	for _, p := range renderedPages(t) {
		body := strings.ToLower(html.UnescapeString(p.body))
		for _, name := range removed {
			if strings.Contains(body, name) {
				t.Errorf("%s names %q, which is not a section any more", p.name, name)
			}
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

		// Every chapter reports the language beside the prompt it governs.
		if !strings.Contains(section, "language-locked") {
			t.Errorf("%s: does not show which language is in force", p.name)
		}

		if p.isCh0 {
			setup := setupSection.FindString(p.body)
			if !strings.Contains(setup, "data-language-select") {
				t.Errorf("%s: the picker belongs in Setup, where the project is prepared", p.name)
			}
			if strings.Contains(section, "data-language-select") {
				t.Errorf("%s: the picker is still in Build It", p.name)
			}
			// Setup carries both choices, and each opens on a blank rather
			// than on a default nobody picked.
			want := len(content.Languages) + len(content.Systems) + 2
			if n := len(optionTag.FindAllString(setup, -1)); n != want {
				t.Errorf("%s: offers %d options, want %d", p.name, n, want)
			}
			// The choice is on this page, so pointing at the page is no help.
			if !strings.Contains(section, `href="#setup"`) {
				t.Errorf("%s: does not point back at the picker above it", p.name)
			}
			continue
		}

		if strings.Contains(p.body, "data-language-select") {
			t.Errorf("%s: only chapter 0 may offer the picker", p.name)
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
		num, _ = strconv.Atoi(regexp.MustCompile(`chapter-(\d+)`).FindStringSubmatch(p.name)[1])
		// A chapter renders the system name when a prompt says {os}, and on
		// the chapter that hands over the runner script, which is written for
		// one system and has to say which.
		wantsName := namesASystem(byNumber[num]) || num == content.RunnerChapter
		hasName := strings.Contains(p.body, "data-os-name")
		if wantsName != hasName {
			t.Errorf("%s: names an operating system = %v, want %v", p.name, hasName, wantsName)
		}

		// The picker sits in Setup beside the language, because chapter 3
		// already asks for commands that run on the reader's machine. Sending
		// them to chapter 10 to correct it means three chapters of commands
		// for a system they are not on.
		wantPicker := num == systemPickerChapter()
		hasPicker := strings.Contains(p.body, "data-system-select")
		if wantPicker != hasPicker {
			t.Errorf("%s: offers the system picker = %v, want %v", p.name, hasPicker, wantPicker)
		}
		if wantPicker {
			for _, sys := range content.Systems {
				if !strings.Contains(p.body, `>`+sys.Name+`</option>`) {
					t.Errorf("%s: the picker does not offer %s", p.name, sys.Name)
				}
			}
		}
		// Every system's script is on the chapter that hands it over, and
		// exactly one is shown.
		if num == content.RunnerChapter {
			for _, r := range content.RunnerScripts {
				if !strings.Contains(p.body, `data-runner-for="`+r.SystemID+`"`) {
					t.Errorf("%s: carries no runner script for %s", p.name, r.SystemID)
				}
			}
		}
		// A chapter that names a system has to say where to change it.
		if wantsName && !strings.Contains(p.body, "change in Chapter") {
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
		want := 0
		num := 0
		num, _ = strconv.Atoi(regexp.MustCompile(`chapter-(\d+)`).FindStringSubmatch(p.name)[1])
		want = len(byNumber[num].BuildIt.Prompts)
		// A sidebar's turns are copied like any other.
		if aside := byNumber[num].Aside; aside != nil {
			want += len(aside.BuildIt.Prompts)
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

// The Portal is a component, named like the Vault and the Teller. This checks
// the rendered page rather than the content, because the two places it was
// still lowercase were the template's divider and the preamble the generator
// writes into every portal prompt, neither of which lives in content.
func TestPortalIsCapitalisedEverywhereItIsNamed(t *testing.T) {
	// The folder is a path, not the component's name.
	path := regexp.MustCompile(`peyva/portal/`)
	for _, p := range renderedPages(t) {
		body := path.ReplaceAllString(p.body, "")
		for _, wrong := range []string{"the portal", "The portal", "a portal"} {
			if strings.Contains(body, wrong) {
				t.Errorf("%s: writes %q, but Portal is a component name", p.name, wrong)
			}
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
// whichever of the backend languages the reader picked.
func TestPortalPromptCarriesItsOwnRules(t *testing.T) {
	byNumber := map[int]content.ChapterContent{}
	for _, c := range content.All {
		byNumber[c.Number] = c
	}

	for _, p := range renderedPages(t) {
		num := 0
		num, _ = strconv.Atoi(regexp.MustCompile(`chapter-(\d+)`).FindStringSubmatch(p.name)[1])
		if len(byNumber[num].BuildIt.PortalPrompts()) == 0 {
			continue
		}
		body := html.UnescapeString(p.body)
		for _, rule := range []string{
			"plain HTML and CSS",
			"No framework, no build step, no dependencies",
			"peyva/portal/",
			"peyva/goal.md",
			"never everyone's at once",
			"switcher at the top says whose",
		} {
			if !strings.Contains(body, rule) {
				t.Errorf("%s: portal prompt is missing %q", p.name, rule)
			}
		}
	}
}

// A portal prompt must not name a language either. The page is the same page
// whichever language serves it, which is the only reason one portal can exist
// across all of them.
func TestPortalPromptsNameNoLanguage(t *testing.T) {
	banned := []string{"go run", "goroutine", "fmt.Println", "npm", "React", "Lit", "node_modules"}
	for _, c := range content.All {
		for _, b := range banned {
			for _, prompt := range c.BuildIt.PortalPrompts() {
				if strings.Contains(prompt.Text, b) {
					t.Errorf("chapter %d portal prompt names %q", c.Number, b)
				}
			}
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
		hasSetup := strings.Contains(p.body, `<section class="block setup"`)
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

// The entry point is the landing page. It has to say what the site is, list
// every chapter, and hand the reader to chapter 0, and it carries none of a
// chapter's machinery: no prompts, no pickers, no setup files.
func TestLandingPage(t *testing.T) {
	_, indexPath := testRun(t)
	b, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("reading landing page: %v", err)
	}
	body := string(b)
	if m := unrendered.FindString(body); m != "" {
		t.Errorf("landing page: template action left unrendered: %s", m)
	}
	for _, want := range []string{
		"Learn distributed systems by building one",
		"Start with Chapter 0",
		`href="chapter-0.html"`,
		"How a chapter works",
		"No solutions. On purpose.",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("landing page missing %q", want)
		}
	}
	// The chapter list is the sidebar, the same one every chapter shows.
	if n := len(sidebarLink.FindAllString(body, -1)); n != len(content.All) {
		t.Errorf("landing page sidebar lists %d chapters, want %d", n, len(content.All))
	}
	for _, c := range content.All {
		if !strings.Contains(body, html.EscapeString(c.Title)) {
			t.Errorf("landing page does not list chapter %d, %q", c.Number, c.Title)
		}
	}
	// Home is the active sidebar entry here, and a link everywhere else.
	if !strings.Contains(body, `<li class="active" data-slug="index">Home</li>`) {
		t.Error("landing page: Home is not the active sidebar entry")
	}
	// It wears the chapter shell, so the same sections in the same order.
	for _, s := range requiredSections {
		if !strings.Contains(body, s) {
			t.Errorf("landing page missing section %q", s)
		}
	}
	for _, absent := range []string{"data-copy-prompt", "data-language-select", "data-system-select", `<section class="block setup"`} {
		if strings.Contains(body, absent) {
			t.Errorf("landing page carries %s, which belongs to a chapter", absent)
		}
	}
}
