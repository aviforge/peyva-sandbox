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
			{Label: "Reach it", Text: `My environment: {os}. If you can run commands, run them here rather than advising about another machine. The Gateway is listening on TCP port 9310, and I have a phone on the same Wi-Fi.

I want to reach the Gateway from the phone. Tell me what is happening, not only what to type.

Give me the command to find this machine's network address, and how to tell which of the addresses it prints is the right one. Say whether the Gateway needs a change to accept connections from another machine. Warn me before any firewall prompt. If the phone cannot reach it, ask once whether this is a locked-down network, then have me try from a second terminal instead.

Done when the phone, or that second terminal, reaches the Gateway and it logs the connection.`},
			{Label: "Stream", Text: `The Gateway listens on TCP port 9310, accepts a connection, logs one line with the caller's address, then closes it.

Show me that TCP is a stream and not messages. Have it print each read with its byte count, then send one message from a client in two writes with a pause between them. Say what the Gateway would need in order to know where that message ended.

Done when I have seen one message arrive as more than one read.`},
			{Label: "Portal", Portal: true, Text: `The Gateway accepts connections but serves nothing. Have it return peyva/portal/index.html to anything that connects, so the page is reachable from my phone.

Tell me what changes about how the page loads its stylesheet now it arrives over a connection instead of from disk.

Leave the Gateway running when you are done. The point is reaching it from the phone next.

Done when the phone, or a second terminal, shows alice's balance.`},
		},
	},
}
