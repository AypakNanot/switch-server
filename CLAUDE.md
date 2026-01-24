<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

---

# Git Workflow Instructions

**DO NOT automatically commit or push code.**

When working with this repository:

1. **NEVER** automatically execute `git commit`
2. **NEVER** automatically execute `git push`
3. **ASK** the user before performing any commit/push operations
4. **SHOW** the proposed commit message and await confirmation
5. **OR** display the git commands for the user to execute manually

If the user asks to "commit" or "push" changes, first show them what will be committed and ask for confirmation.