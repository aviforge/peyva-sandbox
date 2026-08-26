package content

var Chapter18 = ChapterContent{
	Number:     18,
	Slug:       "chapter-18",
	Title:      "Lock It Down: Security",
	Subtitle:   "Security protects peyva and its users from attackers, keeps data safe, and builds trust.",
	Category:   "Operations",
	Difficulty: "Advanced",
	EstTime:    "25 min",
	QuickTip:   "Authentication proves who's calling; authorization checks what they're allowed to do. You need both.",

	HeroImage:   "images/chapter-18.webp",
	HeroCaption: "Security is not a feature, it's a foundation. We build it in, not bolt it on.",

	Intuition: []string{
		"Everything peyva does so far assumes every caller is honestly Alice or Bob.",
		"A real bank verifies who you are, checks what you're allowed to do, and logs everything.",
		"Security applies that same discipline: prove identity, limit access, protect data, watch for trouble.",
	},

	Concepts: []ConceptItem{
		{Term: "AuthN (Authentication)", Description: "Verifying who someone is. Proving Alice is really Alice."},
		{Term: "AuthZ (Authorization)", Description: "Checking what an authenticated user is allowed to do, least privilege by default."},
		{Term: "Encryption in Transit", Description: "Protecting data as it travels over the network (TLS), so it can't be read if intercepted."},
		{Term: "Secrets Management", Description: "Storing credentials and keys securely, never hardcoded in source code."},
	},

	UnderTheHood: []string{
		"Users -> an encrypted connection -> the Gateway, which checks who is calling, what they are allowed to do, and how often they may ask -> the Teller -> stored data, itself encrypted and backed up.",
		"Both checks run before the Teller is called, so a request that fails either one never reaches a balance.",
		"Credentials live outside the code, in the environment or a secrets store, so a copy of the repository is never a copy of the keys.",
	},

	BuildIt: BuildIt{
		Intro:     "The Gateway learns to prove who's calling before the Teller sees a payment.",
		Technique: "Chain-of-Verification (CoVe)",
		Why:       "Asking for secure code gets you the checklist. Verifying your own draft finds the holes the checklist does not mention.",
		Source:    "The Prompt Report: Self-Criticism, Chain-of-Verification",
		Prompts: []Prompt{
			{Label: "Draft", Intro: "The authentication and the ownership check.", Text: `The Gateway trusts the "from" field on a payment request completely. Any caller can move money out of any account by naming it.

Make the Gateway prove who the caller is, and confirm they own the account they are spending from before the Teller ever sees the request. Move any credential out of source code into an environment variable.

Done when a request with no credential is refused, and a caller authenticated as one owner cannot spend from another's account.`},
			{Label: "Plan the checks", Thinking: true, Intro: "The questions, written before the answers.", Text: `You put authentication in front of a payments API and an ownership check between the caller and the money.

Write the list of questions that would expose that work as broken. Be specific to this system, these fields and these checks: no generic OWASP categories, no advice that would apply to any application.

Don't answer them yet.

Done when I have a list of questions that are all about this system, and none of them could be asked of any other.`},
			{Label: "Answer and revise", Intro: "Answer each one, then fix what they expose.", Text: `You wrote a list of questions that would expose your authentication and ownership checks as broken.

Take each one against the code you actually wrote, one at a time. Don't soften an answer because of what you concluded on another question.

Then fix what the answers exposed, and state plainly what is still exploitable, including anything you left out of scope on purpose.

Done when every question has an answer, the fixable ones are fixed, and I have your list of what remains exploitable.`},
			{Label: "Portal", Portal: true, Intro: "The sign-in.", Text: `The Portal's switcher takes whoever it is told. Anyone at the keyboard can pick alice and send her money, which was fine while peyva ran on one laptop and is not fine now.

Put a sign-in in front of it. Switching account means signing in as that account, and the switcher offers only accounts already signed in. Signing out removes one.

A signed-in customer sees their own account and nobody else's, and can only send from their own.

Done when signing in as alice shows alice, and the switcher offers nobody she has not signed in as.`},
			{Label: "Portal checks", Portal: true, Intro: "Try to get past it.", Text: `You put a sign-in in front of a wallet page, where switching account means signing in as that account.

Write the list of questions that would expose that page as broken: specific to this page and these forms, not generic security advice. Answer each one against what you built, and fix what the answers expose.

Done when nothing I can type in the browser makes it show or spend bob's money without bob's own sign-in, and I have your list of what you checked.`},
		},
	},

	BreakIt: BreakIt{
		Intro: "Try to misuse peyva the way an attacker would.",
		Exercises: []string{
			"Send a /transfer request with no auth token. Confirm it's rejected with 401, not silently processed.",
			"Authenticate as Alice but set 'from' to Bob's account. Confirm it's rejected with 403, proving authentication alone isn't authorization.",
			"Search the codebase for any hardcoded credentials. There should be none by the end of this chapter.",
		},
	},
}
