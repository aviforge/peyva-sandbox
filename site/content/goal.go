package content

// GoalSpec is what the reader saves as peyva/goal.md before writing any code.
// Every prompt afterwards points at it rather than restating the goal.
//
// It lives here, and is rendered into chapter 0 with a copy button, rather than
// sitting as a file in the repository. A reader following the published site
// should never have to leave it to fetch something, and a file on GitHub is a
// second place the site cannot reach.
//
// It carries the invariants because nothing else does. The prompts say where
// code goes and how to hold money; none of them says money must never be
// created or lost, which is the property the whole book exists to teach.
const GoalSpec = `# peyva

A peer-to-peer wallet. Alice sends Bob $20, and the money arrives exactly once.

## What it must do

- Hold an account for each user: a handle, an owner and an amount.
- Move money between accounts on request.
- Answer what an account holds, and how it came to hold it.
- Give one customer a page of their own: what they hold, sending money to a
  handle, and the history behind the number. Reached from a menu, not a
  command line.
- Let whoever is at the keyboard change which customer that is.

## What must never happen

These hold in every chapter. If a change would break one, the change is wrong.

- Money is never created and never lost. Every debit has a matching credit.
- No balance goes negative.
- No payment is applied twice, however many times it is submitted.
- Only the Vault changes a balance. Nothing else writes one.
- Every movement of money is recorded. The balance is the answer, the record is
  the proof.

## Constraints

- Standard library only. No frameworks, no brokers, no third-party libraries.
  The one exception is the runner, a script you are given rather than build.
- Every process is configured only by environment: PEYVA_PORT for what it
  listens on, PEYVA_VAULT for the copies, PEYVA_PEERS for the proxy,
  PEYVA_PRIMARY for a replica, PEYVA_WARDEN for whoever asks who may write.
- Settings that differ between runs or machines are config. Settings with only
  one correct value, the money rules above among them, are code and stay in it.
- Runs on one laptop. No containers, no cloud, no deployment.
- One process until chapter 10, several after it. A transaction never spans
  two of them.
- Code lives in peyva/<component>/, one folder per component.
- Money is exact: a decimal type, or integer minor units where the language has
  none. Never floating point. Two decimal places.
- The Portal is plain HTML and CSS. No build step, no dependencies.
- The Portal shows one customer's wallet at a time, never everyone's at once.
  A switcher at the top says whose.
- Structure is earned. A component appears when a problem needs it, not before.

## Components

Each appears in the chapter that builds it. Until then it does not exist.

- Vault: holds every account and what is in it. Chapter 0.
- Portal: one customer's own wallet, with a switcher for whose. Chapter 0, and
  a little more in most chapters after it.
- Gateway: takes requests from outside. Chapter 2.
- Teller: handles one payment end to end, and is the only thing that moves
  money. Chapter 4.
- Ledger: the append-only record behind every balance. Chapter 7.
- Runner: starts every process, says which are alive, and stops all of them.
  Chapter 10.
- Courier: carries out work after a payment clears. Chapter 12.
- Warden: grants a time-limited lease to the one Vault allowed to write.
  Chapter 16.
- Config: reads every setting from outside the code and checks it. Chapter 19.
- Reconciler: proves the Vault and the Ledger still agree. Chapter 21.
`
