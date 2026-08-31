package content

var Chapter06 = ChapterContent{
	Number:     6,
	Slug:       "chapter-6",
	Title:      "Finding Things Fast (Indexes)",
	Subtitle:   "You don't scan every page. You use the index.",
	Category:   "Databases",
	Difficulty: "Intermediate",
	QuickTip:   "An index makes reads fast and writes a little slower. Every one has to be kept up to date.",

	HeroImage:   "images/chapter-6.webp",
	HeroCaption: "Index = a map for your data. Fast lookups, less scanning.",

	Why: []string{
		"With no index, finding one row means reading every row. Fine at a thousand rows, hopeless at ten million.",
		"An index is a sorted list pointing at the rows. Looking something up takes a few steps whatever the size.",
		"Indexes are paid for in writing. Every insert also updates every index on the table.",
		"An index only helps a search on its own column. Anything else still reads the whole table.",
		"Never trust a speed claim you have not timed. Ask for the query plan, then time it on this machine.",
	},

	Concepts: []ConceptItem{
		{Term: "Full Table Scan", Description: "Checking every row one by one until you find the match. Correct, and slower as the table grows."},
		{Term: "Index", Description: "A separate sorted list that points straight at the rows holding a value."},
		{Term: "B-Tree", Description: "The shape databases keep an index in, so a lookup takes a handful of steps however large it grows."},
		{Term: "Write Amplification", Description: "Every row you write also updates every index on the table. More indexes, slower writes."},
		{Term: "Query Plan", Description: "The database's own account of how it will run a query, and whether it will use an index. Ask for it before believing a query is fast."},
	},

	BuildIt: BuildIt{
		Technique: "Grounding: investigate before claiming",
		Why:       "Rule out claims about code it has not run. What comes back is measured timings and query plans, not a confident guess about what an index does.",
		Source:    "Anthropic: Prompting best practices, Minimizing hallucinations in agentic coding",
		Prompts: []Prompt{
			{Label: "Build", Text: `The Vault stores accounts by owner, and each carries a region that nothing indexes.

In order:

1. Seed 50,000 accounts, each in one of four regions: north, south, east, west.
2. Time a lookup of every account in one region. Show the number and the query plan.
3. Add an index on region.
4. Time it again. Show both numbers and the new plan.
5. Time 1,000 inserts with the index, then without.

Expose GET /accounts/by-region/north on the Gateway: the count and the milliseconds the query took. Started with PEYVA_NO_INDEX=1, the Vault drops the index, so the difference can be seen from outside.

Only report what you measured here. If the gain is smaller than expected, say so.

Done when I have four measured numbers and two query plans that differ.`},
			{Label: "Try", Reader: true, Text: `Time the lookup yourself. With the program running, run this: it asks for one region five times and prints the milliseconds each took.

You should see: five small numbers, most of them alike.`,
				Commands: Commands(
					`1..5 | ForEach-Object { curl.exe -s http://127.0.0.1:9310/accounts/by-region/north -w '\n' }`,
					`for /l %i in (1,1,5) do @curl.exe -s http://127.0.0.1:9310/accounts/by-region/north -w "\n"`,
					`for i in 1 2 3 4 5; do curl -s http://127.0.0.1:9310/accounts/by-region/north -w '\n'; done`,
				)},
			{Label: "Try", Reader: true, Text: `Stop the program. In the terminal you start it from, set this, then start it the way you did before and run the five lookups again. Close that terminal when you are done and the setting goes with it.

You should see: the same count with a far larger number of milliseconds every time. That gap is the index: without it the database reads all 50,000 rows to find the ones in north.`,
				Commands: Commands(
					`$env:PEYVA_NO_INDEX = '1'`,
					`set PEYVA_NO_INDEX=1`,
					`export PEYVA_NO_INDEX=1`,
				)},
			{Label: "Portal", Portal: true, Text: `Someone typing a handle should know they have the right person before sending. Have Send look the handle up as it is typed and show whose account it is.

Show the handle and owner, never the balance.

Done when a wrong handle is obvious before sending, and you have timed the lookup with the index and without.`},
		},
	},
}
