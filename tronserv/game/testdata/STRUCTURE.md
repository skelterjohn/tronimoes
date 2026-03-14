# Game testdata structure

These JSON files are serialized `game.Game` values from `tronserv/game/game.go`. Tests load them with `decodeGame(t, "label")` (no `.json`); the label is the filename without extension.

## Root: Game

Relevant fields for testdata (other fields like `created`, `version`, `pickup`, `code`, `history` may be present in JSON but are irrelevant for tests):

| Field | Type | Description |
|-------|------|-------------|
| `done` | bool | Game over |
| `players` | []*Player | Players in turn order |
| `turn` | int | Index into `players` for current turn |
| `rounds` | []*Round | Completed + current round(s). Only the last may be in progress (`done: false`). |
| `bag` | []*Tile | Remaining tiles (not yet drawn) |
| `board_width` | int | Board width (depends on player count: 1→6, 2→8, 3→10, 4→12, 5→14, 6→16) |
| `board_height` | int | Board height (1→7, 2→9, 3→11, 4→13, 5→15, 6→17) |
| `max_pips` | int | Highest pip value in this game (1→6, 2→7, 3→8, 4→10, 5→11, 6→12) |

## Player

Relevant fields for testdata (other fields like `hints`, `spacer_hints`, `chicken_foot_url`, `react_url`, `kills` may be present but are irrelevant for tests):

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Unique player name |
| `ready` | bool | Ready for next round |
| `score` | int | Round score |
| `hand` | []*Tile | Tiles in hand (typically 7–10) |
| `chicken_foot` | bool | Player is “on the foot” (must play from chicken foot next) |
| `dead` | bool | Player is out this round |
| `just_drew` | bool | Drew from bag this turn (must pass or play) |
| `chicken_foot_coord` | Coord | Where their chicken foot is; meaningful when `chicken_foot` is true |

## Tile

| Field | Type | Description |
|-------|------|-------------|
| `pips_a` | int | One side pips (0..max_pips) |
| `pips_b` | int | Other side pips |

Doubles have `pips_a == pips_b`. Round leader is a double.

## Coord

| Field | Type |
|-------|------|
| `x` | int |
| `y` | int |

Board coordinates; origin and axes are defined by the game (e.g. round leader often near center).

## LaidTile (tile on the board)

| Field | Type | Description |
|-------|------|-------------|
| `tile` | *Tile | The tile |
| `coord` | Coord | Position of the **A** end of the tile |
| `orientation` | string | Which adjacent cell holds the **B** end. From `coord`: `"right"` = (x+1, y), `"left"` = (x-1, y), `"down"` = (x, y+1), `"up"` = (x, y-1). |
| `player_name` | string | Owner of the line this tile extends; "" for round leader |
| `next_pips` | int | Pips at the open end(s) for the next tile to match |
| `dead` | bool | Tile is dead (e.g. round over, line killed) |
| `indicated` | *Tile or null | Optional disambiguation; use `{"pips_a":-1,"pips_b":-1}` for “no indication” |
| `who_laid_it` | string | Player who laid this tile |

**Coord and orientation:** `coord` is the A side; the B side occupies the adjacent cell given by orientation (right/left = x±1, down/up = y±1).

**Consistency:** Each tile appears in `round.laid_tiles` and again in the appropriate `round.player_lines[name]` (or free line). The round leader appears in every player’s line as the first tile.

## Round

| Field | Type | Description |
|-------|------|-------------|
| `laid_tiles` | []*LaidTile | All tiles laid this round, in order |
| `spacer` | *Spacer or null | Current spacer for starting free lines; null if none |
| `done` | bool | Round finished |
| `history` | []string | Round-level history |
| `player_lines` | map[string][]*LaidTile | Per-player lines. Keys = player names. Each line starts with the round leader, then that player’s tiles. |
| `free_lines` | [][]*LaidTile | Lines started from a spacer (no single owner) |
| `bagless_passes` | int | Consecutive passes with empty bag (stalemate when ≥ len(players)) |
| `highest_leader` | int | Highest double that started a line this round (round leader or free line) |

## Spacer

| Field | Type | Description |
|-------|------|-------------|
| `a` | Coord | One end (must be adjacent to a line head) |
| `b` | Coord | Other end; exactly 5 steps in one axis from `a`; straight line with no tiles in between |

## Test usage (game package)

- `decodeGame(t, "basic_report")` → loads `testdata/basic_report.json`, returns `*Game`.
- Tests typically use `game.CurrentRound(ctx)`, `round.FindLegalMoves(...)`, etc., and assert on moves or state.

When adding new testdata, keep `player_lines` and `laid_tiles` consistent, and ensure `turn` indexes the correct player in `players`.

**Prettify:** After editing JSON testdata (or when a file is minified), run from `tronserv`: `go run ./prettify game/testdata/<filename>.json` to rewrite the file with tab-indented formatting.
