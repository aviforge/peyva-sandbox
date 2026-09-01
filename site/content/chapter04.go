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
		What:      "Showing a few examples of the input and the answer you want, so the assistant copies the pattern.",
		Why:       "Four requests with their exact replies settle field names and status codes that prose would leave open.",
		Source:    "The Prompt Report: In-Context Learning; Anthropic, Use examples effectively",
		SourceURL: PromptReportURL + "#Ch2.S2.SS1",
		Prompts: []Prompt{
			{Label: "Build", Text: `Requests reach the Gateway but nothing acts on them. Have it speak HTTP, and build the Teller behind it.

The Teller handles a payment start to finish: check the request, check the payer has enough, move the money in the Vault, return a reference. Only the Teller moves money.

The Gateway takes POST /pay for a payment, POST /accounts with {"handle": "carol"} to open one, and GET /accounts/carol answers {"handle": "carol", "balance": "0.00"}. Match these exactly:

  {"from": "alice", "to": "bob", "amount": 20}
  -> 200 {"status": "success", "reference": "tx_4c8a1f6b2e9d05374a1c8f2b6d0e9a35"}

  {"from": "alice", "to": "bob"}
  -> 400 {"error": "amount must be greater than zero"}

  {"from": "alice", "to": "nobody", "amount": 20}
  -> 404 {"error": "unknown account: nobody"}

  not json
  -> 400 {"error": "invalid JSON body"}

The reference is sixteen random bytes written as thirty-two hex characters, so no two payments ever share one.

Check everything before touching an account. No sign-in, no retries, no disk yet.

Done when those four requests give those four answers, and alice drops by 20 only on the first.`},
			{Label: "Try", Reader: true, Text: `The four requests the assistant tested, sent by you. With the program running, run this: it sends the good payment, the missing amount, the unknown account and the broken JSON in turn, then reads alice's balance. If your paths are not /pay and /accounts, ask the assistant to rename them to these. The rest of the book uses them.

You should see: the four answers from the prompt, in order, with their status codes, and alice at 80.00 after them. Only the first one moved money.`,
				Commands: Commands(
					`curl.exe -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{\"from\":\"alice\",\"to\":\"bob\",\"amount\":20}' -w ' -> %{http_code}\n'
curl.exe -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{\"from\":\"alice\",\"to\":\"bob\"}' -w ' -> %{http_code}\n'
curl.exe -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{\"from\":\"alice\",\"to\":\"nobody\",\"amount\":20}' -w ' -> %{http_code}\n'
curl.exe -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d 'not json' -w ' -> %{http_code}\n'
curl.exe -s http://127.0.0.1:9310/accounts/alice -w '\n'`,
					`curl.exe -s -X POST http://127.0.0.1:9310/pay -H "Content-Type: application/json" -d "{\"from\":\"alice\",\"to\":\"bob\",\"amount\":20}" -w " -> %{http_code}\n"
curl.exe -s -X POST http://127.0.0.1:9310/pay -H "Content-Type: application/json" -d "{\"from\":\"alice\",\"to\":\"bob\"}" -w " -> %{http_code}\n"
curl.exe -s -X POST http://127.0.0.1:9310/pay -H "Content-Type: application/json" -d "{\"from\":\"alice\",\"to\":\"nobody\",\"amount\":20}" -w " -> %{http_code}\n"
curl.exe -s -X POST http://127.0.0.1:9310/pay -H "Content-Type: application/json" -d "not json" -w " -> %{http_code}\n"
curl.exe -s http://127.0.0.1:9310/accounts/alice -w "\n"`,
					`curl -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{"from":"alice","to":"bob","amount":20}' -w ' -> %{http_code}\n'
curl -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{"from":"alice","to":"bob"}' -w ' -> %{http_code}\n'
curl -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d '{"from":"alice","to":"nobody","amount":20}' -w ' -> %{http_code}\n'
curl -s -X POST http://127.0.0.1:9310/pay -H 'Content-Type: application/json' -d 'not json' -w ' -> %{http_code}\n'
curl -s http://127.0.0.1:9310/accounts/alice -w '\n'`,
				)},
			{Label: "Portal", Portal: true, Text: `The Portal shows a balance and nothing else. Add Send, and a way to open an account by handle.

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
