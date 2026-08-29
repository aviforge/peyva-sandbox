package content

var Chapter11 = ChapterContent{
	Number:     11,
	Slug:       "chapter-11",
	Title:      "Sharing Work: Caching",
	Subtitle:   "Caching is like keeping popular items ready in advance, so we don't make (or check) everything again.",
	Category:   "System Design",
	Difficulty: "Intermediate",
	QuickTip:   "A stale cache is worse than no cache.",

	HeroImage:   "images/chapter-11.webp",
	HeroCaption: "Cache = keep frequently used data nearby, so it's fast to get back.",

	Why: []string{
		"A cache is a second copy, and any second copy can disagree with the first.",
		"Only the program that sees every write can clear the cache correctly. So the cache lives in the Vault, not in the copies.",
		"Clear it after the write is saved, never before. Clear it early and another read puts the old value straight back.",
		"A payment changes two accounts. Clear both.",
		"A balance on a page may be a moment old. The balance checked before taking money may not be cached at all.",
		"A cache bug never crashes. It hands back a fast, confident, wrong number.",
	},

	Concepts: []ConceptItem{
		{Term: "Cache", Description: "A small fast store holding what was asked for recently, so the next ask is quick."},
		{Term: "Hit and Miss", Description: "A hit answers from the cache. A miss goes to storage and keeps the answer for next time."},
		{Term: "Invalidation", Description: "Clearing or updating a cached value when the real one changes, so the cache never lies."},
		{Term: "Cache Placement", Description: "Which program holds the cache. Only one that sees every write can clear it correctly, so balances are cached in the Vault, not the copies."},
		{Term: "Read-Your-Writes", Description: "Someone who just paid must see the new balance. A cache that can hand them the old one has broken a promise they will notice."},
	},

	BuildIt: BuildIt{
		Technique: "Self-Refine",
		Why:       "Have it review its own work and fix what it finds, until a pass finds nothing. Nothing here crashes, so running the tests will not tell you.",
		Source:    "The Prompt Report: Self-Criticism, Self-Refine",
		Prompts: []Prompt{
			{Label: "Build", Text: `Every balance enquiry reaches the Vault's storage, even for a balance read a moment ago. Three copies sit in front of the Vault, which is its own process.

Add a cache of balances in memory inside the Vault, cleared whenever a payment changes an account. A map safe for several requests at once, no Redis, no library. Enquiries only: the check before a debit never reads the cache.

First, two sentences on why the cache cannot live in the copies.

Done when a repeated enquiry comes from cache, a payment makes the next enquiry from any copy show the new balance, and a payment's own check never reads the cache.`},
			{Label: "Check", Text: `You added a cache of balances inside the Vault, cleared when a payment changes an account.

Try to make it serve an old balance. Check every path that changes a balance clears it. Consider a read and a write at once, a payment rolled back after the clear, a payment touching two accounts, and a read refilling the cache between the clear and the commit.

List what you found, fix it, review again. Repeat until a review finds nothing.

Done when a review finds nothing, and I know how many rounds it took and what the first version got wrong.`},
			{Label: "Portal", Portal: true, Text: `The Portal works, and looks like each piece was added the week it was needed. Make it presentable.

Judge it against every rule in peyva/portal/design.md and the visual idea from chapter 0. Say what passes before changing anything. Fix what fails, judge again, repeat.

Done when a pass finds nothing to fix.`},
		},
	},
}
