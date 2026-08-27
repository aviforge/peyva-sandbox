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
			{Label: "Draft", Text: `The Gateway believes the "from" field, so anyone can spend from any account. The Vault takes writes from anything that reaches its port.

Make the Gateway prove who is calling and check they own the account, before the Teller sees the request. Make the Vault take writes only from callers with a shared secret. Every password and key comes from the environment, and a process refuses to start without it.

Done when a request with no credential is refused, alice cannot spend from bob's account, and a write sent straight to the Vault without the secret is refused.`},
			{Label: "Plan the checks", Thinking: true, Text: `You put a sign-in in front of a payments API, an ownership check before the money, and a shared secret in front of the store.

Write the questions that would show that work to be broken. Specific to this system, these fields and these checks. Nothing that could be asked of any application.

Do not answer them yet.

Done when I have questions that are all about this system and no other.`},
			{Label: "Answer and revise", Text: `You wrote the questions that would show your sign-in and ownership checks to be broken.

Answer each one against the code you actually wrote. Then fix what the answers exposed, and say plainly what is still exploitable.

Done when every question has an answer, the fixable ones are fixed, and I have your list of what remains.`},
			{Label: "Portal", Portal: true, Text: `The Portal's switcher takes whoever it is told, so anyone at the keyboard can spend alice's money.

Put a sign-in in front of it. Switching means signing in as that account, and the switcher offers only accounts already signed in. A signed-in customer sees only their own account and sends only from it.

Done when signing out as alice and in as bob changes every screen to bob, and closing the tab leaves neither signed in.`},
			{Label: "Portal checks", Portal: true, Text: `You put a sign-in in front of a wallet page, where switching account means signing in as that account.

Write the questions that would show that page to be broken, specific to this page and these forms. Answer each against what you built, and fix what they expose.

Done when nothing I can type in the browser shows or spends bob's money without bob's sign-in.`},
		},
	},
}
