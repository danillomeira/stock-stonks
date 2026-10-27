---
name: find-skills
description: Use when the user asks "how do I do X", says "find a skill for X", asks whether Claude can do something, or when a task in this repo would clearly benefit from a capability (a CLI, a workflow, a checklist) that no current skill covers. Searches the open agent-skills ecosystem (skills.sh) and installs the best match instead of improvising from scratch.
---

# Find Skills

A meta-skill: before building ad-hoc tooling for a recurring kind of task,
check whether a maintained skill already does it well.

## When to use

- The user asks "how do I do X" / "can you do X" / "find a skill for X".
- A task in this repo (spec writing, ADRs, Go tooling, Next.js patterns,
  release notes, changelogs, etc.) would clearly benefit from a dedicated
  skill that isn't installed yet.
- The user expresses interest in extending what this agent can do.

## CLI reference

- `npx skills find [query] [--owner <owner>]` — search by keyword.
- `npx skills add <owner/repo@skill>` — install from GitHub or another source.
- `npx skills update` — update installed skills.
- `npx skills init <name>` — scaffold a new skill when nothing matches.

## Workflow

1. Identify the domain and the specific task. Ask: does a common, well-known
   skill likely exist for this (e.g. "writing ADRs", "generating OpenAPI
   docs", "Go linting conventions")?
2. Check the leaderboard at https://skills.sh before searching broadly — it
   ranks skills by install count.
3. Run `npx skills find` with specific keywords. Try alternate terms if the
   first search is thin (e.g. "ADR" vs "architecture decision record" vs
   "decision log").
4. Verify quality before recommending:
   - Prefer skills with 1K+ installs.
   - Treat anything under 100 installs with caution.
   - Favor official/verified sources over random forks.
   - Be skeptical of repos with fewer than 100 GitHub stars.
   - **Do not recommend a skill based solely on search results** — check the
     source repo briefly.
5. Present options to the user: name, install count, source, install
   command, and the skills.sh link. Let them choose.
6. Offer to install with:
   ```
   npx skills add <owner/repo@skill> -g -y
   ```
   (`-g` installs globally, `-y` skips confirmation prompts.)

## Category cheat sheet

| Domain | Example queries |
|---|---|
| Web development | "react forms", "next.js routing", "accessibility audit" |
| Testing | "test coverage", "e2e playwright", "mocking" |
| DevOps | "github actions", "docker multi-stage", "terraform" |
| Documentation | "api docs", "changelog", "readme generator" |
| Code quality | "lint config", "code review checklist" |
| Design | "design tokens", "figma sync" |
| Productivity | "task tracking", "meeting notes" |

## Fallback

If nothing matches after trying alternate keywords, say so plainly, offer to
help with the task directly, and suggest `npx skills init <name>` if this
looks like something worth turning into a reusable skill for next time.
