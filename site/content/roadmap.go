package content

type RoadmapEntry struct {
	Number int
	Title  string
}

// Roadmap lists every planned chapter for the sidebar, in order. Entries
// whose Number has no matching ChapterContent in All render as
// disabled/non-clickable in the sidebar. Titles match the exact wording
// baked into each chapter's hero image (site/assets/images/chapter-<N>.webp).
var Roadmap = []RoadmapEntry{
	{0, "What Are We Building?"},
	{1, "Inside One Computer"},
	{2, "Finding Peyva (Processes & Ports)"},
	{3, "Across the Wire (Networking)"},
	{4, "Designing the API"},
	{5, "Storing Money (Databases)"},
	{6, "Finding Things Fast (Indexes)"},
	{7, "Making It Safe (Transactions)"},
	{8, "Exactly Once (Idempotency)"},
	{9, "How Big Is PEYVA? (Capacity Estimation)"},
	{10, "Growing the Team: Scale Out"},
	{11, "Sharing Work: Caching"},
	{12, "Decoupling with Messages: Queues"},
	{13, "Reliability Patterns: Transactional Outbox"},
	{14, "Big Changes Safely: Sagas"},
	{15, "Data Copies: Replication"},
	{16, "When Things Fail: CAP / Consistency"},
	{17, "See Everything: Observability"},
	{18, "Lock It Down: Security"},
	{19, "Operating in Production"},
	{20, "Putting It All Together"},
}

// Labs are non-numbered bonus pages, listed separately in the sidebar.
var Labs = []RoadmapEntry{
	{-1, "Break It on Purpose"},
}
