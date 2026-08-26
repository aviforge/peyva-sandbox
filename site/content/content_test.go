package content

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestRoadmapHasTwentyOneChapters(t *testing.T) {
	if len(Roadmap) != 21 {
		t.Errorf("len(Roadmap) = %d, want 21", len(Roadmap))
	}
}

func TestAllHasTwentyOneChapters(t *testing.T) {
	if len(All) != 21 {
		t.Errorf("len(All) = %d, want 21", len(All))
	}
}

func TestEveryChapterHasRequiredFields(t *testing.T) {
	seenNumbers := map[int]bool{}
	seenSlugs := map[string]bool{}

	for _, c := range All {
		if c.Title == "" {
			t.Errorf("chapter %d: Title is empty", c.Number)
		}
		if c.Slug == "" {
			t.Errorf("chapter %d: Slug is empty", c.Number)
		}
		if c.HeroImage == "" {
			t.Errorf("chapter %d: HeroImage is empty", c.Number)
		}
		if len(c.Intuition) == 0 {
			t.Errorf("chapter %d: Intuition is empty", c.Number)
		}
		if len(c.Concepts) == 0 {
			t.Errorf("chapter %d: expected at least one concept", c.Number)
		}
		if len(c.UnderTheHood) == 0 {
			t.Errorf("chapter %d: expected at least one under-the-hood point", c.Number)
		}
		if c.BuildIt.Technique == "" {
			t.Errorf("chapter %d: BuildIt.Technique is empty", c.Number)
		}
		if c.BuildIt.Why == "" {
			t.Errorf("chapter %d: BuildIt.Why is empty", c.Number)
		}
		if len(c.BuildIt.Prompts) == 0 {
			t.Errorf("chapter %d: no prompts", c.Number)
		}
		for _, prompt := range c.BuildIt.Prompts {
			if prompt.Label == "" {
				t.Errorf("chapter %d: a prompt has no label", c.Number)
			}
			// Every turn ends somewhere the reader can check, not only the
			// last one. A middle turn without this is a turn nobody can tell
			// they have finished.
			if !strings.Contains(prompt.Text, "Done when") {
				t.Errorf("chapter %d: prompt %q is missing a 'Done when' line", c.Number, prompt.Label)
			}
			// A prompt is copied out of the page on its own, so it must not
			// depend on anything only visible here.
			for _, dangling := range []string{"steps above", "the steps", "above in"} {
				if strings.Contains(prompt.Text, dangling) {
					t.Errorf("chapter %d: prompt %q refers to off-page content (%q). It must stand alone",
						c.Number, prompt.Label, dangling)
				}
			}
		}
		if len(c.BreakIt.Exercises) == 0 {
			t.Errorf("chapter %d: expected at least one break-it exercise", c.Number)
		}

		if seenNumbers[c.Number] {
			t.Errorf("duplicate chapter Number: %d", c.Number)
		}
		seenNumbers[c.Number] = true

		if seenSlugs[c.Slug] {
			t.Errorf("duplicate chapter Slug: %q", c.Slug)
		}
		seenSlugs[c.Slug] = true
	}
}

// A technique nobody else has named is a technique we made up, and a made-up
// name is worse than no name: the reader carries it into rooms where it means
// nothing. Every technique must cite a corpus the reader can go and read.
func TestEveryTechniqueCitesARecognisedSource(t *testing.T) {
	for _, c := range All {
		if c.BuildIt.Source == "" {
			t.Errorf("chapter %d: BuildIt.Source is empty. Technique %q must cite where it is documented",
				c.Number, c.BuildIt.Technique)
			continue
		}
		recognised := false
		for _, corpus := range RecognisedSources {
			if strings.HasPrefix(c.BuildIt.Source, corpus) {
				recognised = true
				break
			}
		}
		if !recognised {
			t.Errorf("chapter %d: BuildIt.Source %q cites no recognised corpus (want one of %v)",
				c.Number, c.BuildIt.Source, RecognisedSources)
		}
	}
}

// Teaching a range of prompting techniques is a goal of the book, so a
// repeated technique means a chapter is missing its lesson.
func TestEveryChapterTeachesADistinctTechnique(t *testing.T) {
	seen := map[string]int{}
	for _, c := range All {
		if prev, ok := seen[c.BuildIt.Technique]; ok {
			t.Errorf("chapters %d and %d both use technique %q", prev, c.Number, c.BuildIt.Technique)
		}
		seen[c.BuildIt.Technique] = c.Number
	}
}

func TestChapterNumbersMatchRoadmap(t *testing.T) {
	roadmapTitles := map[int]string{}
	for _, r := range Roadmap {
		roadmapTitles[r.Number] = r.Title
	}

	for _, c := range All {
		want, ok := roadmapTitles[c.Number]
		if !ok {
			t.Errorf("chapter %d (%q) has no matching Roadmap entry", c.Number, c.Title)
			continue
		}
		if c.Title != want {
			t.Errorf("chapter %d: Title %q does not match Roadmap title %q", c.Number, c.Title, want)
		}
	}
}

// componentChapters maps each named component of the system to the chapter
// that introduces it. The names were chosen deliberately so the reader builds
// one architecture across the book rather than twenty-one unrelated exercises.
var componentChapters = map[string]int{
	"Vault":      0,
	"Gateway":    2,
	"Teller":     4,
	"Ledger":     7,
	"Courier":    12,
	"Config":     19,
	"Reconciler": 20,
}

// A component named in Build It but absent from Concepts is a word the reader
// is asked to build without ever being taught. Concepts is where this book
// defines its vocabulary, so that is where a component has to be introduced,
// not only inside the prompt, which is written to be copied out of the page.
func TestEveryComponentIsDefinedInTheConceptsOfItsChapter(t *testing.T) {
	byNumber := map[int]ChapterContent{}
	for _, c := range All {
		byNumber[c.Number] = c
	}

	for component, number := range componentChapters {
		chapter, ok := byNumber[number]
		if !ok {
			t.Errorf("component %q claims chapter %d, which does not exist", component, number)
			continue
		}
		defined := false
		for _, concept := range chapter.Concepts {
			if concept.Term == component {
				defined = true
				if concept.Description == "" {
					t.Errorf("chapter %d: Concept %q has an empty description", number, component)
				}
				break
			}
		}
		if !defined {
			t.Errorf("chapter %d introduces %q in Build It but never defines it in Concepts",
				number, component)
		}
	}
}

// The component belongs to the chapter that builds it. Defining it earlier
// spends the name before the reader has anything to attach it to.
func TestNoChapterNamesAComponentBeforeItsOwnChapter(t *testing.T) {
	for _, c := range All {
		for _, concept := range c.Concepts {
			if intro, ok := componentChapters[concept.Term]; ok && intro != c.Number {
				t.Errorf("chapter %d defines %q, but %q is introduced in chapter %d",
					c.Number, concept.Term, concept.Term, intro)
			}
		}
	}
}

// A chapter may point back at something the reader has built. It may not point
// forward at a chapter number they have not reached, which reads as a promise
// the page cannot keep and tells them nothing they can act on now.
func TestNoChapterPointsAtALaterChapter(t *testing.T) {
	ref := regexp.MustCompile(`Chapter\s+(\d+)|Ch\.\s*(\d+)`)
	for _, c := range All {
		for _, m := range ref.FindAllStringSubmatch(everything(c), -1) {
			num := m[1]
			if num == "" {
				num = m[2]
			}
			target, err := strconv.Atoi(num)
			if err != nil {
				continue
			}
			if target > c.Number {
				t.Errorf("chapter %d points forward at chapter %d", c.Number, target)
			}
		}
	}
}

// The spec says a component does not exist until the chapter that builds it, so
// no chapter may name one earlier. Chapter 4 listed the Ledger among what sits
// behind the API, three chapters before it was built.
func TestNoComponentIsNamedBeforeItIsBuilt(t *testing.T) {
	for name, intro := range componentChapters {
		word := regexp.MustCompile(`\b` + name + `\b`)
		for _, c := range All {
			if c.Number >= intro {
				continue
			}
			if word.MatchString(everything(c)) {
				t.Errorf("chapter %d names the %s, which is built in chapter %d",
					c.Number, name, intro)
			}
		}
	}
}

// The chapter that introduces Config has to say both halves of the rule. Half
// of it is worse than none: a reader told to externalise settings and not told
// where that stops moves the money rules into a file someone can edit, and the
// invariants become optional without anyone deciding they should be.
func TestConfigChapterStatesBothHalvesOfTheRule(t *testing.T) {
	var chapter ChapterContent
	for _, c := range All {
		if c.Number == componentChapters["Config"] {
			chapter = c
		}
	}
	if chapter.Title == "" {
		t.Fatalf("no chapter %d, which is where Config is built", componentChapters["Config"])
	}

	text := everything(chapter)
	for _, half := range []struct{ name, phrase string }{
		{"what to externalise", "differs between one run"},
		{"what to keep in code", "only one correct value"},
		{"the money rules stay in code", "two decimal places"},
		{"how to decide", "3am"},
	} {
		if !strings.Contains(text, half.phrase) {
			t.Errorf("chapter %d does not say %s (looked for %q)",
				chapter.Number, half.name, half.phrase)
		}
	}

	// Money's shape is stated as a constraint in the spec, so a chapter that
	// later offered it as a setting would contradict the file every prompt
	// points at.
	if !strings.Contains(GoalSpec, "Two decimal places") {
		t.Error("the spec no longer fixes money at two decimal places")
	}
}

// everything is every reader-facing string on a chapter, so a check cannot pass
// only because the text moved to a field it does not look at.
func everything(c ChapterContent) string {
	parts := []string{
		c.Title, c.Subtitle, c.QuickTip, c.HeroCaption,
		c.BuildIt.Intro, c.BuildIt.Why, c.BreakIt.Intro,
	}
	for _, p := range c.BuildIt.Prompts {
		parts = append(parts, p.Label, p.Intro, p.Text)
	}
	parts = append(parts, c.Intuition...)
	parts = append(parts, c.UnderTheHood...)
	parts = append(parts, c.BreakIt.Exercises...)
	for _, x := range c.Concepts {
		parts = append(parts, x.Term, x.Description)
	}
	return strings.Join(parts, "\n")
}

// Chapter 10 replaces the single process with several copies behind a proxy.
// Anything after it that still scopes itself to one process is describing a
// system the reader stopped running two chapters ago. This has been wrong twice
// already: the spec said it, and chapter 12 scoped the Courier by it.
func TestNothingAfterScaleOutClaimsOneProcess(t *testing.T) {
	const scaleOut = 10
	for _, c := range All {
		if c.Number < scaleOut {
			continue
		}
		fields := map[string]string{
			"Prompts":  promptText(c),
			"QuickTip": c.QuickTip,
		}
		for _, s := range append(append([]string{}, c.Intuition...), c.UnderTheHood...) {
			fields[s] = s
		}
		for name, text := range fields {
			// "one process can't do that alone" and "what may live inside one
			// process" are about the limit, not a claim peyva still is one.
			for _, claim := range []string{"one process on a laptop", "one process, one user"} {
				if strings.Contains(strings.ToLower(text), claim) {
					t.Errorf("chapter %d %s scopes peyva to %q, but chapter %d runs several copies",
						c.Number, name, claim, scaleOut)
				}
			}
		}
	}
	if strings.Contains(GoalSpec, "One process on a laptop") {
		t.Error("the spec scopes peyva to one process, which stops being true at chapter 10")
	}
}

// The long prose fields are Go raw strings, which cannot contain a backtick.
// A prompt that wants one has to be reworded or the field has to go back to a
// quoted string with escapes, and the compiler will say so. This says why
// before someone spends time on it.
func TestProseFieldsHoldNoBackticks(t *testing.T) {
	for _, c := range All {
		fields := map[string]string{
			"Prompts": promptText(c),
			"Why":     c.BuildIt.Why,
		}
		for name, text := range fields {
			if strings.Contains(text, "`") {
				t.Errorf("chapter %d %s contains a backtick, which a raw string cannot hold",
					c.Number, name)
			}
		}
	}
	for name, text := range map[string]string{"GoalSpec": GoalSpec, "AgentRules": AgentRules} {
		if strings.Contains(text, "`") {
			t.Errorf("%s contains a backtick, which a raw string cannot hold", name)
		}
	}
}

// promptText is every prompt of a chapter joined, for checks that care whether
// a phrase appears anywhere in what the reader is asked to paste.
func promptText(c ChapterContent) string {
	var b strings.Builder
	for _, p := range c.BuildIt.Prompts {
		b.WriteString(p.Text)
		b.WriteString("\n")
	}
	return b.String()
}
