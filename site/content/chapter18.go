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
		"The 'from' field is a claim. Acting on it unchecked lets any caller spend from any account.",
		"Authentication says who is calling. Authorisation says whether alice may spend from bob. Both, or neither counts.",
		"Check at the edge, in one place, before the money. A spread-out check has one handler that forgot.",
		"Secrets come from the environment, never source. The process refuses to start without them.",
		"Internal calls are not trusted for being internal. Anything on the network can reach the Vault's port.",
		"Review by asking how this system breaks, not by applying a generic checklist.",
	},

	Concepts: []ConceptItem{
		{Term: "AuthN (Authentication)", Description: "Verifying who someone is. Proving Alice is really Alice."},
		{Term: "AuthZ (Authorisation)", Description: "Checking what an authenticated user is allowed to do, least privilege by default."},
		{Term: "Encryption in Transit", Description: "Protecting data as it travels over the network (TLS), so it can't be read if intercepted."},
		{Term: "Secrets Management", Description: "Storing credentials and keys securely, never hardcoded in source code."},
		{Term: "Trust Boundary", Description: "The line across which a caller's claims stop being believed. The Gateway is one; the Vault's port is another, because anything on the network can reach it."},
	},

	BuildIt: BuildIt{
		Technique: "Chain-of-Verification (CoVe)",
		Why:       "Asking for secure code gets you the checklist. Verifying your own draft finds the holes the checklist does not mention.",
		Source:    "The Prompt Report: Self-Criticism, Chain-of-Verification",
		Prompts: []Prompt{
			{Label: "Draft", Text: `The Gateway trusts the "from" field on a payment request completely. Any caller can move money out of any account by naming it. The Vault accepts a write from anything that can reach its port.

Make the Gateway prove who the caller is, and confirm they own the account they are spending from before the Teller ever sees the request. Make the Vault accept writes only from callers presenting a shared secret the copies hold. Move every credential out of source code into environment variables the processes refuse to start without.

Done when a request with no credential is refused, a caller authenticated as one owner cannot spend from another's account, and a write sent straight to the Vault's port without the secret is refused.`},
			{Label: "Plan the checks", Thinking: true, Text: `You put authentication in front of a payments API, an ownership check between the caller and the money, and a shared secret in front of the store.

Write the list of questions that would expose that work as broken. Be specific to this system, these fields and these checks: no generic OWASP categories, no advice that would apply to any application.

Don't answer them yet.

Done when I have a list of questions that are all about this system, and none of them could be asked of any other.`},
			{Label: "Answer and revise", Text: `You wrote a list of questions that would expose your authentication and ownership checks as broken.

Take each one against the code you actually wrote, one at a time. Don't soften an answer because of what you concluded on another question.

Then fix what the answers exposed, and state plainly what is still exploitable, including anything you left out of scope on purpose.

Done when every question has an answer, the fixable ones are fixed, and I have your list of what remains exploitable.`},
			{Label: "Portal", Portal: true, Text: `The Portal's switcher takes whoever it is told. Anyone at the keyboard can pick alice and send her money, which was fine while peyva ran on one laptop and is not fine now.

Put a sign-in in front of it. Switching account means signing in as that account, and the switcher offers only accounts already signed in. Signing out removes one.

A signed-in customer sees their own account and nobody else's, and can only send from their own.

Done when signing out as alice and back in as bob changes every screen to bob, and closing the tab does not leave either of them signed in.`},
			{Label: "Portal checks", Portal: true, Text: `You put a sign-in in front of a wallet page, where switching account means signing in as that account.

Write the list of questions that would expose that page as broken: specific to this page and these forms, not generic security advice. Answer each one against what you built, and fix what the answers expose.

Done when nothing I can type in the browser makes it show or spend bob's money without bob's own sign-in, and I have your list of what you checked.`},
		},
	},
}
