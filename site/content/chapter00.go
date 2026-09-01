package content

var Chapter00 = ChapterContent{
	Number:     0,
	Slug:       "chapter-0",
	Title:      "What Are We Building?",
	Subtitle:   "A peer-to-peer wallet, built one honest layer at a time.",
	Category:   "Foundations",
	Difficulty: "Beginner",
	QuickTip:   "Every system starts as one honest process. Resist the urge to design for scale you don't have yet.",

	HeroImage:   "images/chapter-0.webp",
	HeroCaption: "Alice sends Bob $20: the story, and what's really happening underneath.",

	Why: []string{
		"A wallet is the smallest system where every mistake costs someone real money.",
		"Five rules never break: money is never created or lost, never negative, never paid twice, only the Vault changes it, every move is written down.",
		"You can check them. All the balances added up equal what you started with, and each account's history adds up to its balance.",
		"Never store money as a fraction: a computer cannot hold 0.10 exactly. Count whole pennies.",
		"One program holding everything in memory is the right start. Add a piece when a failure demands it, not before.",
		"There is no finished code to copy, on purpose. You prompt for it, run it, and watch it fail as the chapter said.",
	},

	Concepts: []ConceptItem{
		{Term: "User", Description: "A person with an identity in the system. Owns exactly one account."},
		{Term: "Account", Description: "A handle, an owner and an amount. The handle is how everyone else refers to it."},
		{Term: "Balance", Description: "How much an account holds. Never negative, never lost, never spent twice."},
		{Term: "Vault", Description: "The component holding every account and what is in it. Nothing else may change a balance."},
		{Term: "Transfer", Description: "A request to move money between accounts."},
		{Term: "History", Description: "The record of every movement in and out of an account, in the order it happened."},
		{Term: "Invariant", Description: "A rule that must be true at every moment, whatever else is going on. The five in goal.md are the point of this book."},
	},

	BuildIt: BuildIt{
		Technique: "Zero-shot prompting",
		What:      "Asking for what you want in plain words, with no examples and no step-by-step instructions.",
		Why:       "The Vault has one obvious shape, so anything more would be scaffolding around a simple job.",
		Source:    "The Prompt Report: Zero-Shot",
		SourceURL: PromptReportURL + "#Ch2.S2.SS1.SSS3",
		Prompts: []Prompt{
			{Label: "Build", Text: `Build the Vault: the component that holds every account, and the only thing allowed to change a balance.

An account is a handle, an owner and an amount. Seed alice with 100 and bob with 0. Print them on startup, run until interrupted, then print a goodbye.

One file. No disk, no network.

Done when it prints alice's balance and exits cleanly on Ctrl+C.`},
			{Label: "Portal", Portal: true, Text: `Write peyva/portal/index.html each time the program starts. It shows one customer's wallet, never everyone's.

A switcher at the top says whose. A menu down the side has one entry, Balance, with room for more later. Put both accounts in the page so switching needs no server.

I open the file from disk. This is where the look is set, so commit to the one visual idea the brief asks for.

Done when the page opens on alice at 100.00 and the switcher shows bob at 0.00.`},
			{Label: "Try", Reader: true, Text: `The page you just opened has the balance written into it. Prove it. With the program running, run this: it changes 100.00 to 999.00 in the file. Reload the page. Then stop the program with Ctrl+C, start it again the way you did before, and reload once more.

You should see: alice at 999.00 after the first reload, and back at 100.00 after the restart. The page is a copy the program writes at startup, not a view of what the Vault holds.`,
				Commands: Commands(
					`(Get-Content peyva\portal\index.html -Raw) -replace '100\.00', '999.00' | Set-Content peyva\portal\index.html`,
					`powershell -Command "(Get-Content peyva\portal\index.html -Raw) -replace '100\.00', '999.00' | Set-Content peyva\portal\index.html"`,
					`sed -i.bak 's/100\.00/999.00/g' peyva/portal/index.html && rm peyva/portal/index.html.bak`,
				)},
		},
	},
}
