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

	Concepts: []ConceptItem{
		{Term: "User", Description: "A person with an identity in the system. Owns exactly one account."},
		{Term: "Account", Description: "A handle, an owner and an amount. The handle is how everyone else refers to it."},
		{Term: "Balance", Description: "How much an account holds. Must never go negative, be lost, or be double-spent."},
		{Term: "Vault", Description: "The component holding every account and what's in it. Nothing else is allowed to change a balance."},
		{Term: "Transfer", Description: "A request to move money between accounts."},
		{Term: "History", Description: "The record of every transfer. If the balance is the answer, history is the proof."},
	},

	BuildIt: BuildIt{
		Intro:     "Build the Vault: the first component, in its smallest useful form.",
		Technique: "Zero-shot prompting",
		Why:       "The task has one obvious shape. Scaffolding it would cost tokens and buy nothing.",
		Source:    "The Prompt Report: Zero-Shot",
		Prompts: []Prompt{
			{Label: "Build", Text: `Build the Vault. The component that holds accounts and what's in them. Nothing else is allowed to change a balance.

An account is a handle, an owner and an amount. Seed two: alice holding 100, bob holding 0. Print the Vault's contents on startup, run until interrupted, then print a shutdown message.

peyva/goal.md holds the goal and the rules money must never break. Read it first.

One file. No persistence, no network.

Done when running it prints alice's balance and it exits cleanly on Ctrl+C.`},
			{Label: "Portal", Portal: true, Intro: "The Portal starts as one customer's wallet, with a way to change which customer.", Text: `Write peyva/portal/index.html each time the program starts. It shows one account at a time, not a list of everyone: whose it is, and what they hold.

Give it the shell the rest of the book fills in. A switcher at the top naming whose wallet is on screen, and a menu down one side with Balance in it. Leave room in the menu for the entries later chapters add, and show nothing that is not built yet.

Both seeded accounts go into the page, so switching between them needs nothing from a server. Say what you used to do the switching and why.

No server. I open the file from disk. Style it so it looks like something you would show a customer, not a debug dump.

Done when the page opens on alice holding 100.00 and the switcher shows me bob holding 0.00.`},
		},
	},
}
