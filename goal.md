# peyva

A peer-to-peer wallet. Alice sends Bob $20, and the money arrives exactly once.

Copy this file to `peyva/goal.md` before you start. Every chapter's prompt
points at it, so your assistant can read the goal and the rules without being
told them again.

## What it must do

- Hold an account for each user: an owner and an amount.
- Move money between accounts on request.
- Answer what an account holds, and how it came to hold it.

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
- One process on a laptop. One user. No deployment.
- Code lives in `peyva/<component>/`, one folder per component.
- Money is exact: a decimal type, or integer minor units where the language has
  none. Never floating point. Two decimal places.
- Structure is earned. A component appears when a problem needs it, not before.

## Components

Each appears in the chapter that builds it. Until then it does not exist.

| Component | What it is for | Chapter |
| --- | --- | --- |
| Vault | Holds every account and what is in it | 0 |
| Gateway | Takes requests from outside | 2 |
| Teller | Handles one payment end to end, and is the only thing that moves money | 4 |
| Ledger | The append-only record behind every balance | 7 |
| Courier | Carries out work after a payment clears | 12 |
| Reconciler | Proves the Vault and the Ledger still agree | 20 |
