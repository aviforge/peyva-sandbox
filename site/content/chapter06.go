package content

var Chapter06 = ChapterContent{
	Number:     6,
	Slug:       "chapter-6",
	Title:      "Finding Things Fast (Indexes)",
	Subtitle:   "You don't scan every page. You use the index.",
	Category:   "Databases",
	Difficulty: "Intermediate",
	EstTime:    "15 min",
	QuickTip:   "An index speeds up reads at the cost of slightly slower writes — every index has to stay current too.",

	HeroImage:   "images/chapter-6.webp",
	HeroCaption: "Index = a map for your data. Fast lookups, less scanning.",

	Intuition: []string{
		"With one account, any search is instant.",
		"With thousands of users, scanning every row is like reading a whole library catalog for one card.",
		"An index lets peyva jump straight to Alice's record instead.",
	},

	Concepts: []ConceptItem{
		{Term: "Full Table Scan", Description: "Checking every row one by one until you find the match — correct, but slow as the table grows."},
		{Term: "Index", Description: "A separate structure that maps a value (like a name) directly to where its row lives."},
		{Term: "B+ Tree", Description: "The tree structure most databases use to keep an index searchable in a handful of steps."},
		{Term: "Index Lookup", Description: "Using the index to jump straight to the data page that holds the row you want."},
	},

	UnderTheHood: []string{
		"Without an index: search(name=\"Alice\") checks Bob, Carol, Dave, Eve... one by one until it finds Alice. Slow.",
		"With an index: the same search walks a small tree (M -> A -> Alice) and lands directly on Alice's record. Fast.",
		"The index tells peyva where the data is; peyva goes straight to it instead of scanning everything.",
	},

	BuildIt: BuildIt{
		Intro:     "The Vault learns to be searched — and you measure whether it worked.",
		Technique: "Grounding: investigate before claiming",
		Why:       "Ask for an optimisation and you'll get one, plus a confident claim that it's faster. The documented remedy is to forbid claims about code the assistant hasn't actually opened or run — which here means timings taken before and after, not estimates.",
		Source:    "Anthropic — Prompting best practices, Minimizing hallucinations in agentic coding",
		Prompt: "The Vault stores accounts keyed by owner. Each account also carries a region, which nothing indexes.\n\n" +
			"Do this in order and don't skip ahead:\n\n" +
			"1. Seed the Vault with 50,000 accounts and randomised regions.\n" +
			"2. Time a lookup of every account in one region. Report the number before changing anything.\n" +
			"3. Add an index on region.\n" +
			"4. Time the same lookup again and report both numbers side by side.\n" +
			"5. Time 1,000 account inserts with the index present, then with it dropped, and report that difference too.\n\n" +
			"Never state a performance claim you have not measured on this machine. Show me real numbers, not estimates, and if the improvement is smaller than expected, say so rather than explaining it away.\n\n" +
			"Done when I have four real numbers: read time before and after, write time before and after.",
	},

	BreakIt: BreakIt{
		Intro: "Indexes aren't free — see the tradeoff.",
		Exercises: []string{
			"Time a lookup on an indexed column versus a full scan on an unindexed one, on a table with thousands of rows.",
			"Insert a new account — the index updates too, so every write pays a small cost to keep it current.",
			"Add five indexes to the same table and compare write speed before and after.",
		},
	},
}
