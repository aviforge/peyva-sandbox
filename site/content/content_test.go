package content

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestRoadmapHasTwentyTwoChapters(t *testing.T) {
	if len(Roadmap) != 22 {
		t.Errorf("len(Roadmap) = %d, want 22", len(Roadmap))
	}
}

func TestAllHasTwentyTwoChapters(t *testing.T) {
	if len(All) != 22 {
		t.Errorf("len(All) = %d, want 22", len(All))
	}
}

// Why is the only place the book itself teaches; everything else on the page
// is vocabulary or a prompt whose answer comes from someone else's model. A
// chapter without it delegates its whole lesson, and a chapter with one bullet
// has a slogan, not a lesson.
//
// The bounds are tight on purpose. A reader skims this section; a bullet past
// two short sentences, or a list past six, is a paragraph wearing bullets and
// nobody reads it.
func TestEveryChapterSaysWhy(t *testing.T) {
	const minBullets, maxBullets, maxWordsPerBullet = 4, 6, 26
	for _, c := range All {
		if len(c.Why) < minBullets || len(c.Why) > maxBullets {
			t.Errorf("chapter %d: Why has %d bullets, want %d to %d", c.Number, len(c.Why), minBullets, maxBullets)
		}
		for i, b := range c.Why {
			if n := len(strings.Fields(b)); n > maxWordsPerBullet {
				t.Errorf("chapter %d: Why bullet %d is %d words, over %d. One claim per bullet", c.Number, i+1, n, maxWordsPerBullet)
			}
			if strings.Contains(b, "—") {
				t.Errorf("chapter %d: Why bullet %d contains an em dash", c.Number, i+1)
			}
		}
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
		if len(c.Concepts) == 0 {
			t.Errorf("chapter %d: expected at least one concept", c.Number)
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
			// they have finished. A Try turn ends with 'You should see:'
			// instead, and its own test holds it to that.
			if !prompt.Reader && !strings.Contains(prompt.Text, "Done when") {
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
	"Warden":     16,
	"Config":     19,
	"Reconciler": 21,
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

	// Lowercased: the check is whether the rule is stated, not how a sentence
	// starting with it happens to be capitalised.
	text := strings.ToLower(everything(chapter))
	for _, half := range []struct{ name, phrase string }{
		{"what to externalise", "differs between one run"},
		{"what to keep in code", "one correct value"},
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

// The Why paragraph says why a technique suits this chapter. It is not the
// place to define the technique, which is named directly above it and cited in
// the Source line below it.
//
// It had grown to thirty-nine words on average, twelve of the twenty-one built
// on the same antithesis: ask for X and you get Y, ask for Z and you get W. Any
// one reads well. Twenty-one of them down a book is a rhythm a reader starts
// hearing instead of reading.
func TestWhyStaysShort(t *testing.T) {
	const maxWords = 28
	for _, c := range All {
		if n := len(strings.Fields(c.BuildIt.Why)); n > maxWords {
			t.Errorf("chapter %d: Why is %d words, over %d. Say why it suits this chapter and stop",
				c.Number, n, maxWords)
		}
	}
}

// No sentence should appear twice in a chapter's prose. A QuickTip restating
// the Break It intro, or an exercise restating a concept, is the reader being
// told the same thing twice on one page.
//
// Prompts are excluded: one is copied out of the page on its own, so it has to
// restate what it needs even when a paragraph above already said it.
func TestNoChapterRepeatsASentence(t *testing.T) {
	for _, c := range All {
		seen := map[string]bool{}
		for _, sentence := range strings.Split(prose(c), "\n") {
			for _, part := range strings.Split(sentence, ". ") {
				part = strings.ToLower(strings.TrimSpace(strings.TrimSuffix(part, ".")))
				if len(strings.Fields(part)) < 6 {
					continue
				}
				if seen[part] {
					t.Errorf("chapter %d says twice: %q", c.Number, part)
				}
				seen[part] = true
			}
		}
	}
}

// The components list says each one appears in the chapter that builds it, so
// it should read in that order. Config and Runner were both appended when they
// were added rather than inserted, which left chapter 10 sitting after chapter
// 20. An appended list always ends up this way; this notices the next one.
//
// Entries wrap, and three of them carry their chapter number on the second line
// of the bullet. The first version of this read line by line and skipped those,
// which meant moving Runner back to the end did not fail it.
func TestSpecListsComponentsInChapterOrder(t *testing.T) {
	section := GoalSpec[strings.Index(GoalSpec, "## Components"):]

	// Rejoin wrapped bullets, so each entry is one string.
	var bullets []string
	for _, line := range strings.Split(section, "\n") {
		switch {
		case strings.HasPrefix(line, "- "):
			bullets = append(bullets, strings.TrimPrefix(line, "- "))
		case strings.HasPrefix(line, "  ") && len(bullets) > 0:
			bullets[len(bullets)-1] += " " + strings.TrimSpace(line)
		}
	}
	if len(bullets) < 5 {
		t.Fatalf("found %d component entries, which is too few to be the list", len(bullets))
	}

	chapter := regexp.MustCompile(`[Cc]hapter (\d+)`)
	last, lastName := -1, ""
	for _, bullet := range bullets {
		name := strings.SplitN(bullet, ":", 2)[0]
		m := chapter.FindStringSubmatch(bullet)
		if m == nil {
			t.Errorf("%s names no chapter", name)
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err != nil {
			t.Errorf("%s has an unreadable chapter number %q", name, m[1])
			continue
		}
		if n < last {
			t.Errorf("%s is chapter %d but comes after %s at chapter %d", name, n, lastName, last)
		}
		last, lastName = n, name
	}
}

// The book is written in British English. Mixed spelling in one document reads
// as carelessness, and a reader notices it without being able to name it.
//
// Source is exempt: it holds the titles of other people's documents, and a
// citation keeps its source's spelling. Chapter 6 cites an Anthropic page
// called Minimizing hallucinations in agentic coding.
func TestSpellingStaysBritish(t *testing.T) {
	american := regexp.MustCompile(`\b\w+iz(e|ed|ing|ation)\b|\bcolor\w*\b|\bbehavior\w*\b|\bcenter\w*\b|\banalyz\w+\b`)
	// Words that end in -ise in both, so -ize is not evidence either way, and
	// the one HTTP header name a prompt has to spell the way the wire does.
	fine := map[string]bool{
		"size": true, "sized": true, "sizing": true, "prize": true, "capsize": true,
		"authorization": true,
	}
	for _, c := range All {
		for _, text := range []string{prose(c), promptText(c)} {
			for _, w := range american.FindAllString(text, -1) {
				if fine[strings.ToLower(w)] {
					continue
				}
				t.Errorf("chapter %d uses the American %q", c.Number, w)
			}
		}
	}
	for _, name := range []string{"GoalSpec", "DesignBrief", "AgentRules"} {
		text := map[string]string{"GoalSpec": GoalSpec, "DesignBrief": DesignBrief, "AgentRules": AgentRules}[name]
		for _, w := range american.FindAllString(text, -1) {
			if fine[strings.ToLower(w)] {
				continue
			}
			t.Errorf("%s uses the American %q", name, w)
		}
	}
}

// prose is everything on the page except the prompt bodies, sidebar included.
func prose(c ChapterContent) string {
	parts := []string{
		c.Title, c.Subtitle, c.QuickTip, c.HeroCaption,
		c.BuildIt.Why,
	}
	parts = append(parts, c.Why...)
	for _, x := range c.Concepts {
		parts = append(parts, x.Term, x.Description)
	}
	for _, p := range c.BuildIt.Prompts {
		parts = append(parts, p.Label)
	}
	if a := c.Aside; a != nil {
		parts = append(parts, a.Title, a.HeroCaption, a.BuildIt.Why)
		parts = append(parts, a.Why...)
		for _, p := range a.BuildIt.Prompts {
			parts = append(parts, p.Label)
		}
	}
	return strings.Join(parts, "\n")
}

// everything is every reader-facing string on a chapter, so a check cannot pass
// only because the text moved to a field it does not look at.
func everything(c ChapterContent) string {
	return prose(c) + "\n" + promptText(c)
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
			"Prompts":  promptText(c),
			"Commands": commandText(c),
			"Why":      c.BuildIt.Why + strings.Join(c.Why, "\n"),
		}
		for name, text := range fields {
			if strings.Contains(text, "`") {
				t.Errorf("chapter %d %s contains a backtick, which a raw string cannot hold",
					c.Number, name)
			}
		}
	}
	for name, text := range map[string]string{"GoalSpec": GoalSpec, "AgentRules": AgentRules, "DesignBrief": DesignBrief} {
		if strings.Contains(text, "`") {
			t.Errorf("%s contains a backtick, which a raw string cannot hold", name)
		}
	}
}

// promptText is every prompt of a chapter joined, sidebar included, for checks
// that care whether a phrase appears anywhere in what the reader is asked to
// paste.
func promptText(c ChapterContent) string {
	var b strings.Builder
	for _, p := range allPrompts(c) {
		b.WriteString(p.Text)
		b.WriteString("\n")
	}
	return b.String()
}

// commandText is every Try command of a chapter joined, sidebar included.
func commandText(c ChapterContent) string {
	var b strings.Builder
	for _, p := range allPrompts(c) {
		for _, cmd := range p.Commands {
			b.WriteString(cmd.Command)
			b.WriteString("\n")
		}
	}
	return b.String()
}

// allPrompts is the chapter's turns followed by its sidebar's.
func allPrompts(c ChapterContent) []Prompt {
	out := append([]Prompt{}, c.BuildIt.Prompts...)
	if c.Aside != nil {
		out = append(out, c.Aside.BuildIt.Prompts...)
	}
	return out
}

// A prompt is read in full before it is copied, so it is the most-read text
// on the page. Past a certain length a reader stops reading and starts
// skimming, and a skimmed prompt gets pasted with a constraint unseen. The cap
// is generous for the multi-stage builds and tight for everything else.
func TestPromptsStayShort(t *testing.T) {
	const maxWords = 170
	for _, c := range All {
		for _, p := range allPrompts(c) {
			if n := len(strings.Fields(p.Text)); n > maxWords {
				t.Errorf("chapter %d prompt %q is %d words, over %d", c.Number, p.Label, n, maxWords)
			}
		}
	}
}

// Five labels, so a reader learns them once. A sixth word on one chapter is
// a word they have to work out on that page alone.
func TestPromptLabelsComeFromTheFixedSet(t *testing.T) {
	allowed := map[string]bool{"Think": true, "Build": true, "Check": true, "Portal": true, "Try": true}
	for _, c := range All {
		for _, p := range allPrompts(c) {
			if !allowed[p.Label] {
				t.Errorf("chapter %d: prompt label %q is not one of Think, Build, Check, Portal, Try", c.Number, p.Label)
			}
			if p.Label == "Portal" && !p.Portal {
				t.Errorf("chapter %d: a turn labelled Portal must be a Portal turn", c.Number)
			}
			if p.Label == "Think" && !p.Thinking {
				t.Errorf("chapter %d: a Think turn must be a Thinking turn", c.Number)
			}
		}
	}
}

// A sidebar is held to the same bar as a chapter: claims of its own, a cited
// technique no chapter's own Build It already teaches, and prompts that end
// somewhere checkable.
func TestEveryAsideMeetsTheChapterBar(t *testing.T) {
	techniques := map[string]int{}
	for _, c := range All {
		techniques[c.BuildIt.Technique] = c.Number
	}
	for _, c := range All {
		a := c.Aside
		if a == nil {
			continue
		}
		if a.Title == "" || a.HeroImage == "" || a.HeroCaption == "" {
			t.Errorf("chapter %d: sidebar is missing a title, image or caption", c.Number)
		}
		if len(a.Why) < 3 {
			t.Errorf("chapter %d: sidebar has %d Why bullets, want at least 3", c.Number, len(a.Why))
		}
		if a.BuildIt.Technique == "" || a.BuildIt.Why == "" || len(a.BuildIt.Prompts) == 0 {
			t.Errorf("chapter %d: sidebar Build It is incomplete", c.Number)
		}
		if n, ok := techniques[a.BuildIt.Technique]; ok {
			t.Errorf("chapter %d: sidebar technique %q is already chapter %d's", c.Number, a.BuildIt.Technique, n)
		}
		recognised := false
		for _, corpus := range RecognisedSources {
			if strings.HasPrefix(a.BuildIt.Source, corpus) {
				recognised = true
			}
		}
		if !recognised {
			t.Errorf("chapter %d: sidebar Source %q cites no recognised corpus", c.Number, a.BuildIt.Source)
		}
		for _, p := range a.BuildIt.Prompts {
			if !p.Reader && !strings.Contains(p.Text, "Done when") {
				t.Errorf("chapter %d: sidebar prompt %q is missing a 'Done when' line", c.Number, p.Label)
			}
		}
	}
}

// tryTurnsThrough is the last chapter whose Try turn has been written. It
// only ever rises. A chapter at or below it with no Try turn is a chapter
// the reader watches instead of operates, which is the gap Try turns exist
// to close.
const tryTurnsThrough = 21

// buildIts is a chapter's Build It and its sidebar's, each as its own list,
// because "first turn" and "the turn before" only mean anything inside one
// list.
func buildIts(c ChapterContent) []BuildIt {
	out := []BuildIt{c.BuildIt}
	if c.Aside != nil {
		out = append(out, c.Aside.BuildIt)
	}
	return out
}

// A Try turn is the reader's, not the assistant's. It is never sent, so it
// carries none of what a sent prompt needs, and it ends where the reader
// should look rather than where the assistant should stop.
func TestTryTurnsAreTheReadersOwn(t *testing.T) {
	systemIDs := map[string]bool{}
	for _, s := range Systems {
		systemIDs[s.ID] = true
	}
	for _, c := range All {
		hasTry := false
		for _, b := range buildIts(c) {
			for i, p := range b.Prompts {
				if p.Reader != (p.Label == "Try") {
					t.Errorf("chapter %d: a Try turn and a Reader turn are the same thing; label %q has Reader=%v", c.Number, p.Label, p.Reader)
				}
				if !p.Reader {
					continue
				}
				hasTry = true
				if p.Portal || p.Thinking {
					t.Errorf("chapter %d: a Try turn is neither Portal nor Thinking", c.Number)
				}
				if strings.Contains(p.Text, "{os}") {
					t.Errorf("chapter %d: a Try turn has no {os}; its commands are per system already", c.Number)
				}
				if strings.Contains(p.Text, "Done when") {
					t.Errorf("chapter %d: a Try turn ends with 'You should see:', not 'Done when'", c.Number)
				}
				paragraphs := strings.Split(strings.TrimSpace(p.Text), "\n\n")
				if last := paragraphs[len(paragraphs)-1]; !strings.HasPrefix(last, "You should see:") {
					t.Errorf("chapter %d: a Try turn's last paragraph must start 'You should see:', got %q", c.Number, last)
				}
				if i == 0 {
					t.Errorf("chapter %d: a Try turn cannot open a Build It; there is nothing to try yet", c.Number)
				} else if prev := b.Prompts[i-1]; prev.Thinking || prev.Label == "Check" {
					t.Errorf("chapter %d: a Try turn follows a Build, Portal or Try turn, not %q", c.Number, prev.Label)
				}
				if len(p.Commands) == 0 {
					continue
				}
				seen := map[string]bool{}
				for _, cmd := range p.Commands {
					if !systemIDs[cmd.SystemID] {
						t.Errorf("chapter %d: Try command for unknown system %q", c.Number, cmd.SystemID)
					}
					if seen[cmd.SystemID] {
						t.Errorf("chapter %d: Try turn has two commands for %q", c.Number, cmd.SystemID)
					}
					seen[cmd.SystemID] = true
					if strings.TrimSpace(cmd.Command) == "" {
						t.Errorf("chapter %d: Try command for %q is empty", c.Number, cmd.SystemID)
					}
					if strings.Contains(cmd.Command, "`") {
						t.Errorf("chapter %d: Try command for %q contains a backtick, which a raw string cannot hold", c.Number, cmd.SystemID)
					}
					if cmd.SystemID != "macos" && cmd.SystemID != "linux" && strings.Contains(cmd.Command, "curl ") {
						t.Errorf("chapter %d: Try command for %q says 'curl', which PowerShell aliases; say curl.exe", c.Number, cmd.SystemID)
					}
				}
				for id := range systemIDs {
					if !seen[id] {
						t.Errorf("chapter %d: Try turn has commands but none for %q", c.Number, id)
					}
				}
			}
		}
		if c.Number <= tryTurnsThrough && !hasTry {
			t.Errorf("chapter %d has no Try turn, and every chapter up to %d must", c.Number, tryTurnsThrough)
		}
	}
}

// The two helpers build the per-system list, so a chapter file states each
// command once and the test above cannot be failed by a missing entry.
func TestCommandHelpersCoverEverySystem(t *testing.T) {
	for name, cmds := range map[string][]SystemCommand{
		"Commands":      Commands("ps", "bat", "sh"),
		"CommandsSplit": CommandsSplit("ps", "bat", "mac", "linux"),
	} {
		if len(cmds) != len(Systems) {
			t.Errorf("%s returned %d commands, want one per system (%d)", name, len(cmds), len(Systems))
		}
	}
	shared := Commands("ps", "bat", "sh")
	if shared[2].Command != "sh" || shared[3].Command != "sh" {
		t.Error("Commands should give macOS and Linux the same shell command")
	}
	split := CommandsSplit("ps", "bat", "mac", "linux")
	if split[2].SystemID != "macos" || split[2].Command != "mac" || split[3].SystemID != "linux" || split[3].Command != "linux" {
		t.Error("CommandsSplit should give macOS and Linux their own commands, in that order")
	}
}
