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
		Why:       "Draft an answer, plan the questions that would catch it being wrong, answer those independently, then revise. Asking for secure code gets you the checklist; verifying your own draft finds the holes the checklist doesn't mention.",
		Source:    "The Prompt Report: Self-Criticism, Chain-of-Verification",
		Prompt: "The Gateway trusts the \"from\" field on a payment request completely. Any caller can move money out of any account.\n\n" +
			"Work in four passes and show me each one.\n\n" +
			"1. Draft. Make the Gateway prove who the caller is, and confirm they own the account they're spending from before the Teller ever sees the request. Move any credential out of source code into an environment variable.\n" +
			"2. Plan the checks. Write the list of questions that would expose your draft as broken. Be specific to this system: no generic OWASP categories.\n" +
			"3. Answer them. Take each question against the code you actually wrote, one at a time, and don't soften an answer because of what you concluded on another.\n" +
			"4. Revise. Fix what the answers exposed, then state plainly what is still exploitable, including anything you left out of scope on purpose.\n\n" +
			"Done when a request with no credential is refused, a caller authenticated as one owner spending from another's account is refused, no credential is left in source, and I have your list of what remains exploitable.",
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
