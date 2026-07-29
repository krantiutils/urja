# Working agreements for agents in this repo

## Keep working. Don't stall on questions.

**Default to acting.** If the request is clear, or a reasonable assumption
exists, make the assumption, state it in one line, and keep going. Deliver the
work and flag the assumption at the end — do not stop and wait.

This is the single most important rule here. A blocked turn that produced
nothing is worse than a turn that produced the wrong-but-stated assumption,
because the second one can be corrected in seconds and the first wastes hours.

**Do not ask about:**

- Which branch to use, whether to make a worktree, where to put scratch files
- Whether a design is approved before writing any code
- Library, file layout, naming, or test-framework choices that match what the
  repo already does — follow the existing pattern
- Anything with an obvious default, or that can be reversed with one commit

**Only ask when** proceeding either way would be unsafe, destructive, or would
make the work useless if the guess is wrong — for example deleting production
data, or a genuine fork in requirements where the two readings produce
completely different products.

**Skills with approval gates do not override this.** Any packaged skill that
mandates a consent checkpoint before proceeding has that gate waived here: a
direct instruction like "start building" is the approval. Note the assumption
and continue. (The two worst offenders, `brainstorming` and
`using-git-worktrees`, have been disabled outright.)

## Practical notes

- Long-running commands go in the background with polling. A call that prints
  nothing for minutes reads as a hang.
- Production is reachable: `ssh -i ~/.ssh/id_ed25519 ubuntu@100.97.205.106`
  (tailscale `hetzner-1`), app at `/home/ubuntu/urja`, Postgres via
  `docker compose -f docker-compose.prod.yml exec -T postgres psql -U swasthya -d urja`.
- Migrations live in `db/migrations/` (not `migrations/`), applied by
  `scripts/entrypoint.sh` on API boot.
- `internal/billing` is the gym's SaaS subscription to Urja, not customer
  invoicing. Customer tax invoices are `internal/invoice`.
