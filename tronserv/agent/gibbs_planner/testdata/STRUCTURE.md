# Gibbs planner testdata structure

These JSON files are the same format as `tronserv/game/testdata`: serialized `game.Game` from `tronserv/game/game.go`. See **tronserv/game/testdata/STRUCTURE.md** for the full Game/Player/Round/LaidTile/Coord/Spacer schema.

## Loading in tests

- `loadGameFromTestdata(t, "label")` loads `testdata/<label>.json` and returns `*game.Game` (label has no `.json`).

## Example: oneshot.json

Used by `TestOneshot` in `gibbs_test.go`:

- **Scenario:** Current player can win the round in one move (play a tile that empties their hand and ends the round).
- **Assertions:** `GetMove` returns a lay (not pass); after applying that move with `g.LayTile`, `g.CurrentRound(ctx)` is nil (round is done).

So this file is a mid-round game state where `g.Players[g.Turn]` has a winning play available.

## Adding new testdata

Do not invent cases; create files only when asked for a specific scenario. When you do:

1. Reuse the same JSON shape as existing files (see game testdata STRUCTURE.md).
2. Ensure `turn` and `players` are consistent; the agent acts as `g.Players[g.Turn]`.
3. Keep round state consistent: `laid_tiles`, `player_lines`, and (if used) `free_lines`/`spacer` must match the game rules in `tronserv/game/game.go`.
