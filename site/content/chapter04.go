package content

var Chapter04 = ChapterContent{
	Number:     4,
	Slug:       "chapter-4",
	Title:      "Designing the API",
	Subtitle:   "You ask for something. They check, do the work, and bring it back. That is a request and a response.",
	Category:   "Design",
	Difficulty: "Beginner",
	QuickTip:   "Validate before you touch any data. Reject bad requests instead of guessing what the caller meant.",

	HeroImage:   "images/chapter-4.webp",
	HeroCaption: "Simple. Predictable. Hard to misuse.",

	Why: []string{
		"An API promises what happens when things fail, not only when they work.",
		"4xx means the same request will fail again. 5xx means it might work later. Swap them and callers retry the wrong things.",
		"Check everything before you change anything. A debit followed by 'no such recipient' has created money out of an error.",
		"The reference is all the caller has to ask about a payment later. Make a new one each time, never reuse one.",
		"200 means the money moved and was written down. Sending it sooner is a promise you might not keep.",
	},

	Concepts: []ConceptItem{
		{Term: "Endpoint", Description: "One specific thing the API lets you do: POST /transfer, for moving money."},
		{Term: "Request Body", Description: "The details of what you are asking for, e.g. {\"from\": \"alice\", \"to\": \"bob\", \"amount\": 20}."},
		{Term: "Status Codes", Description: "What happened, as a number: 2xx it worked, 4xx your request is wrong, 5xx the server failed and it might work later."},
		{Term: "Validation", Description: "Checking the request makes sense before doing any work, instead of guessing what was meant."},
		{Term: "Teller", Description: "The component that handles a payment from start to finish. The Gateway passes requests on and never touches a balance."},
	},

	BuildIt: BuildIt{
		Technique: "Few-shot (multishot) prompting",
		Why:       "Prose leaves field names and status codes up for grabs. Examples pin them down, failures included.",
		Source:    "The Prompt Report: In-Context Learning; Anthropic, Use examples effectively",
		Prompts: []Prompt{
			{Label: "Build", Text: `Requests reach the Gateway but nothing acts on them. Have it speak HTTP, and build the Teller behind it.

The Teller handles a payment start to finish: check the request, check the payer has enough, move the money in the Vault, return a reference. Only the Teller moves money. Add an endpoint that opens an account at zero.

Match these exactly:

  {"from": "alice", "to": "bob", "amount": 20}
  -> 200 {"status": "success", "reference": "tx_7f3b9c2a"}

  {"from": "alice", "to": "bob"}
  -> 400 {"error": "amount must be greater than zero"}

  {"from": "alice", "to": "nobody", "amount": 20}
  -> 404 {"error": "unknown account: nobody"}

  not json
  -> 400 {"error": "invalid JSON body"}

Check everything before touching an account. No sign-in, no retries, no disk yet.

Done when those four requests give those four answers, and alice drops by 20 only on the first.`},
			{Label: "Page", Portal: true, Text: `The Portal shows a balance and nothing else. Add Send, and a way to open an account by handle.

No From field: money leaves whoever the switcher names. The forms post to the endpoints you built and show what comes back:

  open an account   -> the handle appears at 0.00
  send to a handle  -> the reference shows and the balance drops
  amount left blank -> the page says which field is wrong, nothing moves
  unknown handle    -> the page says so, nothing moves

Never a blank screen or a raw error.

Done when those four cases behave as written, and switching to the recipient shows the money arrived.`},
		},
	},
}
