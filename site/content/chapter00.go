package content

var Chapter00 = ChapterContent{
	Number:     0,
	Slug:       "chapter-0",
	Title:      "What Are We Building?",
	Subtitle:   "A peer-to-peer wallet, built one honest layer at a time.",
	Category:   "Foundations",
	Difficulty: "Beginner",
	EstTime:    "15 min",
	QuickTip:   "Every system starts as one honest process. Resist the urge to design for scale you don't have yet.",

	HeroImage:   "images/chapter-0.webp",
	HeroCaption: "Alice sends Bob $20: the story, and what's really happening underneath.",

	Intuition: []string{
		"Alice opens her wallet, picks Bob, types $20, and hits send. Seconds later, Bob has it.",
		"That's the whole system, from the outside.",
		"Every chapter answers one question about that moment: how the request arrives, what the program does with it, what could go wrong.",
	},

	Concepts: []ConceptItem{
		{Term: "User", Description: "A person with an identity in the system. Owns exactly one account."},
		{Term: "Account", Description: "One owner and an amount, the thing this book protects."},
		{Term: "Balance", Description: "How much an account holds. Must never go negative, be lost, or be double-spent."},
		{Term: "Vault", Description: "The component holding every account and what's in it. Nothing else is allowed to change a balance."},
		{Term: "Transfer", Description: "A request to move money between accounts, the core action of this book."},
		{Term: "History", Description: "The record of every transfer. If the balance is the answer, history is the proof."},
	},

	UnderTheHood: []string{
		"Today: one process, one account, one hardcoded balance, printed to the screen.",
		"A system earns its structure one real problem at a time.",
	},

	BuildIt: BuildIt{
		Intro:     "Build the Vault: the first component, in its smallest useful form.",
		Technique: "Zero-shot prompting",
		Why:       "No examples, no reasoning scaffold, no role: just the instruction. The task has one obvious shape, so scaffolding would cost tokens and buy nothing. Every technique after this is a departure from this baseline, which is why it comes first.",
		Source:    "The Prompt Report: Zero-Shot",
		Prompt: "Build the Vault. The component that holds accounts and what's in them. Nothing else is allowed to change a balance.\n\n" +
			"An account is an owner and an amount. Seed it with one account, alice, holding 100. Print the Vault's contents on startup, run until interrupted, then print a shutdown message.\n\n" +
			"One file, standard library only. No persistence, no network. The Vault earns structure later, when a real problem demands it.\n\n" +
			"Done when running it prints alice's balance and it exits cleanly on Ctrl+C.",
	},

	BreakIt: BreakIt{
		Intro: "This is the simplest way a system can fail: turn it off.",
		Exercises: []string{
			"Run it with `go run .`. Confirm it prints the account's starting balance.",
			"Stop it with Ctrl+C, run it again. The balance resets to the hardcoded value.",
			"Nothing survived, and today that is expected: the balance only ever lived in memory, and memory goes when the process does.",
		},
	},
}
