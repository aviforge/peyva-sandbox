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
- Let whoever is at the keyboard change which customer that is, so both ends
  of a payment can be seen.

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
  The one exception is the runner, which is a shell script.
- Runs on one laptop. No containers, no cloud, no deployment.
- One process until chapter 10, several after it, started by a runner in the
  repo rather than by hand.
- Code lives in peyva/<component>/, one folder per component.
- Money is exact: a decimal type, or integer minor units where the language has
  none. Never floating point. Two decimal places.
- The portal is plain HTML and CSS. No build step, no dependencies.
- The portal shows one customer's wallet at a time, never everyone's at once.
  A switcher at the top says whose.
- Structure is earned. A component appears when a problem needs it, not before.

## Components

Each appears in the chapter that builds it. Until then it does not exist.

- Vault: holds every account and what is in it. Chapter 0.
- Gateway: takes requests from outside. Chapter 2.
- Teller: handles one payment end to end, and is the only thing that moves
  money. Chapter 4.
- Ledger: the append-only record behind every balance. Chapter 7.
- Courier: carries out work after a payment clears. Chapter 12.
- Reconciler: proves the Vault and the Ledger still agree. Chapter 20.
- Runner: a script that starts a given number of copies, says which are alive,
  and stops all of them. Chapter 10. A tool for working on peyva, not a part
  of it, so it is the one thing written for the operating system rather than in
  the project's language.
- Portal: one customer's own wallet, with a switcher for whose. Its menu grows
  a chapter at a time, from chapter 0.
`
