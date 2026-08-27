package content

var Chapter06 = ChapterContent{
	Number:     6,
	Slug:       "chapter-6",
	Title:      "Finding Things Fast (Indexes)",
	Subtitle:   "You don't scan every page. You use the index.",
	Category:   "Databases",
	Difficulty: "Intermediate",
	QuickTip:   "An index speeds up reads at the cost of slightly slower writes. Every index has to stay current too.",

	HeroImage:   "images/chapter-6.webp",
	HeroCaption: "Index = a map for your data. Fast lookups, less scanning.",

	Why: []string{
		"Without an index, finding one row means reading every row. That is fine at a thousand accounts and fatal at ten million, and the code looks identical in both cases. Slowness of this kind arrives silently, with growth.",
		"An index is a second, sorted copy of one column pointing back at the rows. A B-tree keeps it balanced so any lookup takes a handful of page reads however large the table grows, which is why a lookup on a million rows costs about the same as on a thousand.",
		"Every index is a write amplifier. Each insert or update has to maintain every index on the table, so a table with ten indexes does eleven writes per row. Indexes are bought with write speed.",
		"An index only helps queries that use its column the way it was sorted. An index on region does nothing for a search on owner, and a query the planner cannot map onto an index scans the table anyway.",
		"Performance claims are measurements, not reasoning. The same query on the same data is fast or slow depending on the page cache, the disk and the size of the table, and the only honest number is one taken on the machine in question.",
		"The primary key is an index too. Looking an account up by handle is fast for the same reason a region lookup becomes fast here, and for no other.",
	},

	Concepts: []ConceptItem{
		{Term: "Full Table Scan", Description: "Checking every row one by one until you find the match: correct, but slow as the table grows."},
		{Term: "Index", Description: "A separate structure that maps a value (like a name) directly to where its row lives."},
		{Term: "B-Tree", Description: "The balanced tree most databases use to keep an index searchable in a handful of page reads, however large it grows."},
		{Term: "Write Amplification", Description: "Each row write also writes every index on the table. More indexes, slower inserts."},
		{Term: "Query Plan", Description: "The database's account of how it will run a query: which index it will use, or that it will scan. Ask for it before believing a query is fast."},
	},

	BuildIt: BuildIt{
		Technique: "Grounding: investigate before claiming",
		Why:       "Forbid claims about code it has not run, and you get timings instead of confident guesses.",
		Source:    "Anthropic: Prompting best practices, Minimizing hallucinations in agentic coding",
		Prompts: []Prompt{
			{Label: "Build", Text: `The Vault stores accounts keyed by owner. Each account also carries a region, which nothing indexes.

Do this in order and don't skip ahead:

1. Seed the Vault with 50,000 accounts and randomised regions.
2. Time a lookup of every account in one region. Report the number before changing anything, and show me the query plan the database gives for it.
3. Add an index on region.
4. Time the same lookup again and report both numbers side by side, with the new query plan.
5. Time 1,000 account inserts with the index present, then with it dropped, and report that difference too.

Never state a performance claim you have not measured on this machine. Show me real numbers, not estimates, and if the improvement is smaller than expected, say so rather than explaining it away.

Done when I have four real numbers: read time before and after, write time before and after, and two query plans that differ.`},
			{Label: "Portal", Portal: true, Text: `Someone typing a handle from memory should know they have the right person before they send. Have Send look the handle up as it is typed and show whose account it is.

Show the recipient's handle and owner, never their balance. That is not the sender's to see.

Then show me, from a real timing you have run rather than an estimate, how long that lookup takes against a table with thousands of accounts, with the index and without it.

Done when a wrong handle is obvious before sending, and you have shown me both timings.`},
		},
	},
}
