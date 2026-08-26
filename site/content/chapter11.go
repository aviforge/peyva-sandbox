package content

var Chapter11 = ChapterContent{
	Number:     11,
	Slug:       "chapter-11",
	Title:      "Sharing Work: Caching",
	Subtitle:   "Caching is like keeping popular items ready in advance, so we don't make (or check) everything again.",
	Category:   "System Design",
	Difficulty: "Intermediate",
	EstTime:    "15 min",
	QuickTip:   "A stale cache is worse than no cache. Invalidation matters more than the cache itself.",

	HeroImage:   "images/chapter-11.webp",
	HeroCaption: "Cache = keep frequently used data nearby, so it's fast to get back.",

	Intuition: []string{
		"Alice checks her balance a dozen times a day.",
		"Each check means another trip to the database: often for data that hasn't changed.",
		"A cache is a small, fast store that keeps recent answers ready, the way a café pre-makes a popular drink instead of brewing it fresh every time.",
	},

	Concepts: []ConceptItem{
		{Term: "Cache", Description: "A fast, small store that holds recently or frequently requested data close at hand."},
		{Term: "Cache Hit", Description: "The data was already in the cache, so the answer comes back without the database being asked at all."},
		{Term: "Cache Miss", Description: "The data wasn't cached. peyva fetches it from the database and stores it in the cache for next time."},
		{Term: "Invalidation", Description: "Removing or updating a cached value when the underlying data changes, so the cache never lies."},
	},

	UnderTheHood: []string{
		"Request Data (Get Alice's balance) -> Cache. Cache Hit -> Return Data (fast). Cache Miss -> Fetch from the Database and store in cache, then return it.",
		"A cache trades a little staleness risk for speed. That tradeoff has to be deliberate, not accidental.",
	},

	BuildIt: BuildIt{
		Intro:     "The Vault learns to answer faster: without ever answering wrong.",
		Technique: "Self-Refine",
		Why:       "Produce, critique your own output against a stated standard, revise, repeat until the critique comes back empty. A cache bug doesn't crash. It returns a confident wrong number, so the critique pass catches what a passing test won't.",
		Source:    "The Prompt Report: Self-Criticism, Self-Refine",
		Prompt: "Every balance enquiry hits the Vault's storage, even when the same balance was read a moment ago and hasn't changed since.\n\n" +
			"Add an in-memory cache in front of the Vault's reads, invalidated whenever a payment changes that account. A map guarded by a mutex is enough: no Redis, no cache library.\n\n" +
			"Once it works, review your own implementation as if you were trying to make it serve a stale balance. Walk every path that changes a balance and check whether it invalidates. Consider a concurrent read and write, a payment that rolled back, and a payment that touches two accounts at once.\n\n" +
			"Report what you found as a list, fix it, then run that same review again on the fixed version. Keep going until a review turns up nothing, and tell me how many rounds it took.\n\n" +
			"Done when a repeated enquiry is served from cache, a payment makes the next enquiry show the new balance, and you've told me which staleness bug your first attempt had.",
	},

	BreakIt: BreakIt{
		Intro: "A cache that's wrong is worse than no cache at all. Prove invalidation actually works.",
		Exercises: []string{
			"Read Alice's balance a few thousand times with the cache on, then again with it off, and compare the totals. A single hit is far too fast to see. The saving only shows up in the aggregate.",
			"Make a transfer that changes Alice's balance, then check it again. It must reflect the new balance, not the stale cached one.",
			"Deliberately skip the invalidation step and repeat the test. Watch the cache confidently return the wrong balance. This is why cache invalidation is famously hard.",
		},
	},
}
