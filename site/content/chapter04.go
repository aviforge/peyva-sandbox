package content

var Chapter04 = ChapterContent{
	Number:     4,
	Slug:       "chapter-4",
	Title:      "Designing the API",
	Subtitle:   "You ask for something. They check, do the work, and bring it back. That is a request and a response.",
	Category:   "Design",
	Difficulty: "Beginner",
	EstTime:    "20 min",
	QuickTip:   "Validate before you touch any data. Reject bad requests instead of guessing what the caller meant.",

	HeroImage:   "images/chapter-4.webp",
	HeroCaption: "Simple. Predictable. Hard to misuse.",

	Intuition: []string{
		"Alice doesn't care about ports and packets: she cares about sending $20.",
		"An API is the menu: a fixed set of requests, in a fixed shape, with a predictable answer.",
		"Order POST /transfer, get back 200 OK and a transaction id.",
	},

	Concepts: []ConceptItem{
		{Term: "Endpoint", Description: "A specific thing the API lets you do: POST /transfer, for moving money."},
		{Term: "Request Body", Description: "The details of what you're asking for, e.g. {\"from\": \"alice\", \"to\": \"bob\", \"amount\": 20}."},
		{Term: "Status Codes", Description: "A predictable signal of what happened: 200 OK, 400 Bad Request, 500 Server Error."},
		{Term: "Validation", Description: "Checking the request makes sense before doing any work, rejecting bad orders instead of guessing."},
		{Term: "Teller", Description: "The component that handles a payment end to end, and the only thing allowed to move money. The Gateway parses and forwards; it never touches a balance itself."},
	},

	UnderTheHood: []string{
		"Client sends 'Request: POST /transfer {...}' to peyva's API; peyva sends back 'Response: 200 OK {...}'.",
		"Behind the API sit the Teller, which decides what to do, and the Vault, which is the only thing that changes a balance.",
		"The caller never learns how peyva stores a balance, only what to send and what comes back.",
	},

	BuildIt: BuildIt{
		Intro:     "Build the Teller, and give the Gateway a real payment request to forward.",
		Technique: "Few-shot (multishot) prompting",
		Why:       "Prose leaves field names and status codes up for grabs. Examples pin them down, failures included.",
		Source:    "The Prompt Report: In-Context Learning; Anthropic, Use examples effectively",
		Prompts: []Prompt{
			{Label: "Build", Text: `Requests reach the Gateway but nothing acts on them. Have it speak HTTP, and build the Teller behind it.

The Teller handles one payment end to end: validate the request, check the payer has enough, move the amount between accounts in the Vault, and return a reference the caller can quote later. The Teller is the only thing allowed to move money: the Gateway parses and forwards, and never touches a balance itself.

Add one more endpoint: open an account for a handle that does not exist yet, starting at zero.

Match these exactly:

  {"from": "alice", "to": "bob", "amount": 20}
  -> 200 {"status": "success", "reference": "tx_7f3b9c2a"}

  {"from": "alice", "to": "bob"}
  -> 400 {"error": "amount must be greater than zero"}

  not json
  -> 400 {"error": "invalid JSON body"}

Validate before touching any account. No authentication, no retries, no persistence: later chapters.

Done when those three requests produce exactly those three responses, and alice's balance drops by 20 only on the first.`},
			{Label: "Portal", Portal: true, Intro: "The menu gains what a wallet is for: sending money, from whoever the switcher says.", Text: `The Portal's menu has one entry, Balance, and no way to move money out of the account it is showing. Add Send to it, and a way to open a new account by handle. A new account joins the switcher.

There is no From field. Money leaves whoever the switcher names.

The forms post to the endpoints you just built and render what comes back. Match these:

  open an account   -> the new handle appears with a balance of 0.00
  send to a handle  -> the reference comes back and the payer's balance drops
  amount left blank -> the page says which field is wrong, and nothing moves
  unknown handle    -> the page says so, and nothing moves

Every failure shows on the page, never a blank screen or a raw error.

Done when those four cases each produce exactly what is written above, and switching to the recipient shows the money arrived.`},
		},
	},

	BreakIt: BreakIt{
		Intro: "A well-designed API should be hard to misuse. Try to misuse it.",
		Exercises: []string{
			"Send a body that is not JSON at all. Expect 400 Bad Request, not a crash.",
			"Send a request missing the 'amount' field. Expect 400, with a clear reason.",
			"Send duplicate requests back to back. Both go through, and peyva has no way yet to tell they were meant to be one payment.",
		},
	},
}
