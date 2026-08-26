package content

var Chapter06 = ChapterContent{
	Number:     6,
	Slug:       "chapter-6",
	Title:      "Finding Things Fast (Indexes)",
	Subtitle:   "You don't scan every page. You use the index.",
	Category:   "Databases",
	Difficulty: "Intermediate",
	EstTime:    "15 min",
	QuickTip:   "An index speeds up reads at the cost of slightly slower writes. Every index has to stay current too.",

	HeroImage:   "images/chapter-6.webp",
	HeroCaption: "Index = a map for your data. Fast lookups, less scanning.",

	Concepts: []ConceptItem{
		{Term: "Full Table Scan", Description: "Checking every row one by one until you find the match: correct, but slow as the table grows."},
		{Term: "Index", Description: "A separate structure that maps a value (like a name) directly to where its row lives."},
		{Term: "B+ Tree", Description: "The tree structure most databases use to keep an index searchable in a handful of steps."},
		{Term: "Index Lookup", Description: "Using the index to jump straight to the data page that holds the row you want."},
	},

	BuildIt: BuildIt{
		Technique: "Grounding: investigate before claiming",
		Why:       "Forbid claims about code it has not run, and you get timings instead of confident guesses.",
		Source:    "Anthropic: Prompting best practices, Minimizing hallucinations in agentic coding",
		Prompts: []Prompt{
			{Label: "Build", Text: `The Vault stores accounts keyed by owner. Each account also carries a region, which nothing indexes.

Do this in order and don't skip ahead:

1. Seed the Vault with 50,000 accounts and randomised regions.
2. Time a lookup of every account in one region. Report the number before changing anything.
3. Add an index on region.
4. Time the same lookup again and report both numbers side by side.
5. Time 1,000 account inserts with the index present, then with it dropped, and report that difference too.

Never state a performance claim you have not measured on this machine. Show me real numbers, not estimates, and if the improvement is smaller than expected, say so rather than explaining it away.

Done when I have four real numbers: read time before and after, write time before and after.`},
			{Label: "Portal", Portal: true, Text: `Someone typing a handle from memory should know they have the right person before they send. Have Send look the handle up as it is typed and show whose account it is.

Show the recipient's handle and owner, never their balance. That is not the sender's to see.

Then show me, from a real timing you have run rather than an estimate, how long that lookup takes against a table with thousands of accounts, with the index and without it.

Done when a wrong handle is obvious before sending, and you have shown me both timings.`},
		},
	},
}
