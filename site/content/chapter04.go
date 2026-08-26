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
		"Behind the API sit the Teller (decides what to do), the Vault (where balances live), and the Ledger (the record of every transfer).",
		"The caller never learns how peyva stores a balance, only what to send and what comes back.",
	},

	BuildIt: BuildIt{
		Intro:     "Build the Teller, and give the Gateway a real payment request to forward.",
		Technique: "Few-shot (multishot) prompting",
		Why:       "Prose leaves field names and status codes up for grabs. Concrete input and output pairs pin the contract harder than a paragraph ever will, and the examples have to include the failures, or only the happy path gets built.",
		Source:    "The Prompt Report: In-Context Learning; Anthropic, Use examples effectively",
		Prompt: "Requests reach the Gateway but nothing acts on them. Have the Gateway speak HTTP, and build the Teller behind it.\n\n" +
			"The Teller handles one payment end to end: validate the request, check the payer has enough, move the amount between accounts in the Vault, and return a reference the caller can quote later. The Teller is the only thing allowed to move money: the Gateway parses and forwards, and never touches a balance itself.\n\n" +
			"Match these exactly:\n\n" +
			"  {\"from\": \"alice\", \"to\": \"bob\", \"amount\": 20}\n" +
			"  -> 200 {\"status\": \"success\", \"reference\": \"tx_7f3b9c2a\"}\n\n" +
			"  {\"from\": \"alice\", \"to\": \"bob\"}\n" +
			"  -> 400 {\"error\": \"amount must be greater than zero\"}\n\n" +
			"  not json\n" +
			"  -> 400 {\"error\": \"invalid JSON body\"}\n\n" +
			"Validate before touching any account. No authentication, no retries, no persistence: later chapters.\n\n" +
			"Done when those three requests produce exactly those three responses, and alice's balance drops by 20 only on the first.",
	},

	BreakIt: BreakIt{
		Intro: "A well-designed API should be hard to misuse. Try to misuse it.",
		Exercises: []string{
			"Send a body that is not JSON at all. Expect 400 Bad Request, not a crash.",
			"Send a request missing the 'amount' field. Expect 400, with a clear reason.",
			"Send duplicate requests back to back. peyva has no way yet to tell they're the same transfer. That's Chapter 8's problem.",
		},
	},
}
