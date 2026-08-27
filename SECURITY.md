# Security Policy

Supported: 0.1.x

Report via https://github.com/Pastalikek65/txs/security/advisories/new — 48h response, 14d fix.

Token handling: `GITHUB_TOKEN` > `GH_TOKEN` > `gh auth token` (2s timeout), never logged. Cache stores only public metadata.
