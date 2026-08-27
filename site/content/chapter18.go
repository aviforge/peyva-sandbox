package content

var Chapter18 = ChapterContent{
	Number:     18,
	Slug:       "chapter-18",
	Title:      "Lock It Down: Security",
	Subtitle:   "Security protects peyva and its users from attackers, keeps data safe, and builds trust.",
	Category:   "Operations",
	Difficulty: "Advanced",
	QuickTip:   "Authentication proves who's calling; authorisation checks what they're allowed to do. You need both.",

	HeroImage:   "images/chapter-18.webp",
	HeroCaption: "Security is not a feature, it's a foundation. We build it in, not bolt it on.",

	Why: []string{
		"The 'from' field is just something the caller typed. Believe it and anyone can spend from anyone's account.",
		"One check proves who is calling. A second checks whether they own the account. You need both.",
		"Check once, at the front door, before any money moves. Spread the check around and one place will forget it.",
		"Passwords and keys come from outside the code. The program refuses to start without them.",
		"A caller is not trusted for being inside. Anything on the network can reach the Vault's port too.",
		"Review by asking how this system breaks, not by reading a checklist that fits any system.",
	},

	Concepts: []ConceptItem{
		{Term: "AuthN (Authentication)", Description: "Proving who is calling: that alice really is alice."},
		{Term: "AuthZ (Authorisation)", Description: "Checking what that person is allowed to do. Give the least access that works."},
		{Term: "Encryption in Transit", Description: "Scrambling data as it crosses the network (TLS), so it is useless to anyone reading along the way."},
		{Term: "Secrets Management", Description: "Keeping passwords and keys outside the code, never typed into a file you commit."},
		{Term: "Trust Boundary", Description: "The line where you stop believing what a caller says about itself. The front door is one; the Vault's port is another."},
	},

	BuildIt: BuildIt{
		Technique: "Chain-of-Verification (CoVe)",
		Why:       "Asking for secure code gets you the checklist. Verifying your own draft finds the holes the checklist does not mention.",
		Source:    "The Prompt Report: Self-Criticism, Chain-of-Verification",
		Prompts: []Prompt{
			{Label: "Draft", Text: `The Gateway believes the "from" field, so any caller can spend from any account. The Vault takes a write from anything that can reach its port.

Make the Gateway prove who the caller is, and check they own the account they are spending from, before the Teller sees the request. Make the Vault take writes only from callers presenting a shared secret. Move every password and key out of the code into environment variables the processes refuse to start without.

Done when a request with no credential is refused, a caller signed in as one owner cannot spend from another's account, and a write sent straight to the Vault's port without the secret is refused.`},
			{Label: "Plan the checks", Thinking: true, Text: `You put a sign-in in front of a payments API, an ownership check between the caller and the money, and a shared secret in front of the store.

Write the questions that would show that work to be broken. Specific to this system, these fields and these checks. No generic categories, nothing that would apply to any application.

Do not answer them yet.

Done when I have questions that are all about this system, and none of them could be asked of any other.`},
			{Label: "Answer and revise", Text: `You wrote the questions that would show your sign-in and ownership checks to be broken.

Take each one against the code you actually wrote, one at a time. Do not soften an answer because of what you concluded on another.

Then fix what the answers exposed, and say plainly what is still exploitable, including anything you left out on purpose.

Done when every question has an answer, the fixable ones are fixed, and I have your list of what remains exploitable.`},
			{Label: "Portal", Portal: true, Text: `The Portal's switcher takes whoever it is told, so anyone at the keyboard can pick alice and spend her money. That was fine on one laptop and is not fine now.

Put a sign-in in front of it. Switching account means signing in as that account, and the switcher offers only accounts already signed in. Signing out removes one.

A signed-in customer sees their own account and nobody else's, and can only send from their own.

Done when signing out as alice and in as bob changes every screen to bob, and closing the tab leaves neither signed in.`},
			{Label: "Portal checks", Portal: true, Text: `You put a sign-in in front of a wallet page, where switching account means signing in as that account.

Write the questions that would show that page to be broken: specific to this page and these forms, not generic advice. Answer each against what you built, and fix what they expose.

Done when nothing I can type in the browser shows or spends bob's money without bob's own sign-in, and I have your list of what you checked.`},
		},
	},
}
