package content

var Chapter03 = ChapterContent{
	Number:     3,
	Slug:       "chapter-3",
	Title:      "Across the Wire (Networking)",
	Subtitle:   "A message reaches someone when you know where to go, which door to use, and how the package travels.",
	Category:   "Foundations",
	Difficulty: "Beginner",
	QuickTip:   "TCP gives you bytes in order, not messages. Where one request ends and the next begins is your job.",

	HeroImage:   "images/chapter-3.webp",
	HeroCaption: "peyva sends data between machines over the network.",

	Why: []string{
		"TCP gives you bytes, not messages. One send can arrive as three reads, or three sends as one.",
		"Marking where a request ends is your job: a length at the front, an end marker, or a header like Content-Length.",
		"TCP keeps bytes in order and resends lost ones. It cannot promise delivery if the connection drops.",
		"Packets get lost, duplicated, delayed and reordered. TCP hides that inside one connection, and nowhere else.",
		"Three things must allow a connection: the address you listen on, the firewall, and the network. All three failures look the same: a timeout.",
	},

	Concepts: []ConceptItem{
		{Term: "IP Address", Description: "Where a machine lives on the network. A street address for computers."},
		{Term: "Protocol", Description: "The agreed way of talking. peyva speaks TCP/IP, with HTTP on top of it later."},
		{Term: "Packets", Description: "The chunks the network moves data in. Any one can be lost, delayed, duplicated or arrive out of order."},
		{Term: "Stream", Description: "What TCP hands your program: bytes in order, with nothing marking where one message ends."},
		{Term: "Framing", Description: "How you mark where a message ends: a length at the front, an end marker, or a header saying how many bytes follow."},
	},

	BuildIt: BuildIt{
		Technique: "Add context and motivation",
		Why:       "Told where you are and why, it picks commands that run on your machine and warns you about the firewall prompt.",
		Source:    "Anthropic: Prompting best practices, Add context to improve performance",
		Prompts: []Prompt{
			{Label: "Build", Text: `My environment: {os}. If you can run commands yourself, this machine is the one to run them on: act on it directly rather than assuming a separate machine you can only advise about. The Gateway is listening on TCP port 9310 on this machine, and I have a phone on the same Wi-Fi.

I want to reach the Gateway from the phone instead of from localhost. I'm doing this to understand how a service becomes reachable beyond the machine it runs on, so tell me what is actually happening rather than only what to type.

Tell me the command for my OS to find this machine's network address, and how to know which of the addresses it prints is the right one. Then tell me whether the Gateway needs a code change to accept connections from another machine, or whether it already does.

If a firewall prompt is likely on my OS, warn me before I hit it. If the phone can't reach it, ask me one thing, whether this is a locked-down network, and go straight to having me check from a second terminal on this machine against that address instead. Don't make me report back symptoms first.

Then show me that TCP is a stream: have the Gateway read what arrives and print each read with its byte count, and send it one message from a client in two separate writes with a pause between them. Tell me what the Gateway would need in order to know where that message ended.

Done when the phone (or that second terminal, if the network blocks it) reaches the Gateway, the Gateway logs the connection, and I have seen one message arrive as more than one read.`},
			{Label: "Portal", Portal: true, Text: `The Gateway accepts connections but serves nothing. Have it return peyva/portal/index.html to anything that connects, so the page I could only open locally is now reachable from my phone on the same Wi-Fi.

Tell me what changes about how the page loads its stylesheet once it arrives over a connection instead of from disk, before I hit it.

Leave the Gateway running once you've checked it works. The point is reaching it from my phone next, not just confirming it starts.

Done when the phone (or a second terminal on this machine, if the network blocks the phone) shows alice's balance.`},
		},
	},
}
