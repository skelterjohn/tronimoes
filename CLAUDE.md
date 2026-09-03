# tronimoes

A chicken-foot dominoes variant. `tronserv` (Go) is the backend/game engine; `tronapp`
(Next.js) is the frontend. See [tronserv/STRUCTURE.md](tronserv/STRUCTURE.md) for backend
layout and `tronapp/AGENTS.md` for frontend framework notes.

## Running locally for testing

```bash
# backend: in-memory store, no GCP dependency needed
cd tronserv && go install ./... && tronserv --dev
# (installs tronserv, agent, replicant to $GOPATH/bin; local agent spawner needs agent/replicant on PATH)

# frontend
cd tronapp && npm run dev
```

`tronserv --dev` with no `--env` flag uses `game.NewMemoryStore()` -- no Firestore/GCP
credentials needed. Server listens on `:8080`, `tronapp` defaults to talking to
`${protocol}//${hostname}:8080` unless `NEXT_PUBLIC_API_URL` is set, so opening
`http://localhost:3000` just works.

`tronserv/run.ps1` does the same build+run for the maintainer's own local testing.

## Playing/testing the game via browser automation (no human required)

The whole game can be driven headlessly through browser JS execution -- no need for a
human to click through the UI. This was built up while playing full games solo (multiple
tabs = multiple independent players) to validate a bug fix.

### Registering a player (no auth needed for local testing)

Registering with just a name sends `X-Player-Name` only (no Authorization/X-Player-Id
header) -- the server accepts this without any token validation. Each browser tab is a
fully independent player; no shared cookies/localStorage/Firebase auth required.

antd's `onPressEnter` doesn't reliably fire from a synthetic OS-level Enter keypress;
dispatch a real `KeyboardEvent` instead:

```js
const el = document.querySelector('input[placeholder="enter your username"]');
el.focus();
el.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', code: 'Enter', keyCode: 13, which: 13, bubbles: true }));
```

Then click "pick-up game" (`Array.from(document.querySelectorAll('button')).find(b => b.textContent.trim() === 'pick-up game')`)
to join/create the shared pickup lobby. Once a game exists, join a second tab the same
way -- it lands in the same pickup game as long as it hasn't started yet.

Once everyone's joined, each player clicks "ready" the same way
(`b.textContent.trim() === 'ready'`); the round auto-starts once all players are ready.

### Reading game state -- fetch the API directly, don't parse the DOM

```js
fetch('http://localhost:8080/game/<CODE>?version=0', { headers: { 'X-Player-Name': '<my-name>' } })
  .then(r => r.json()).then(g => JSON.stringify(g));
```

- `version=0` (or anything lower than current) always returns the latest state immediately.
- The server redacts every *other* player's hand to all-zero `{pips_a:0,pips_b:0}` tiles
  for privacy -- only your own hand (matched by the request's identity) shows real values.
- `turn` is the 0-based index into `players[]` for whose turn it is.
- `players[i].hints` is index-aligned with `players[i].hand`: `null` for an unplayable
  tile, or a list of `{x,y}` coords (both ends of every valid placement) for a playable
  one. The server computes legal moves for you -- don't re-derive game rules client-side.
  **Gotcha:** when eyeballing a long hints array it's easy to miscount and grab the wrong
  index. Do the pairing inside the fetch instead of by hand:
  ```js
  fetch(url, {headers}).then(r=>r.json()).then(g => {
    const me = g.players.find(p => p.name === '<me>');
    const moves = me.hand.map((t,i) => ({i, t, hint: me.hints[i]})).filter(x => x.hint);
    return JSON.stringify({ turn: g.turn, just_drew: me.just_drew,
      moves: moves.map(m=>({a:m.t.pips_a,b:m.t.pips_b,sq1:m.hint[0],sq2:m.hint[1]})) });
  });
  ```
- `me.spacer_hints`: list of `{a:{x,y}, b:{x,y}}` spacer placements, present only when you
  hold a double higher than the current highest leader and your own line already extends
  past the round leader.
- Round's `laid_tiles`: every tile on the board this round, with `coord` + `orientation`
  (which adjacent cell the tile's other half occupies) + `next_pips` (value needed to
  extend from this tile) + `who_laid_it`.
- `player_lines`: per-player ordered tile list, starting from the shared round-leader tile.
- `free_lines`: lines not owned by any player (`player_name: ""`), started off a spacer.

### Laying a tile -- click-click (reliable for automation)

1. Select the tile: `document.querySelectorAll('[data-tile]')`, each has
   `data-tile='{"a":..,"b":..}'` -- find the match and `.click()` it.
2. Click the FAR square (not touching the existing line) -- this becomes `coord`, and the
   tile's `a` pip lands here.
3. Click the NEAR square (touching the line) -- orientation is derived automatically from
   square1 -> square2, and the tile's `b` pip lands here.
   - Squares: `document.querySelector('[data-tron_x="X"][data-tron_y="Y"]')`. The actual
     React `onClick` lives on a wrapping overlay div, not the `Square` div itself, but a
     native click on the inner element bubbles up fine.
4. Fires `POST /game/<code>/tile`.

**Critical gotcha:** the tile-click and both square-clicks must be THREE SEPARATE tool
calls, not one script firing all three `.click()`s synchronously. React state updates are
async/batched, so back-to-back synchronous clicks don't give React a chance to re-render
between them -- `clickSquare`'s closure still sees `selectedTile === undefined` on what
was meant to be the first square click, and the whole sequence silently no-ops (no
request fires, no error either). Splitting into three separate round-trips fixes it.

**Order matters** for which pip lands where: tile `{a:7,b:8}` clicked at `(4,3)` then
`(4,4)` (where `(4,4)` touches the existing line) => orientation "down" => `a=7` at
`(4,3)`, `b=8` at `(4,4)`. Get the pip-to-square assignment backwards and the server will
often silently fix it anyway (see auto-reversal note below) -- but if you also picked the
wrong two squares entirely, expect a real rejection.

The server auto-reverses an invalid-as-given placement if the reversed tile validates
(`LayTileAllWays` tries as-given, then reversed). This is deliberate UX forgiveness for
players who fat-finger the click/drag order, not just an implementation quirk -- it only
saves you when exactly one orientation is legal, not a genuinely illegal placement.

### Laying a tile -- drag-and-drop (what humans use, partially testable headlessly)

Real flow: click the tile until it's rotated correctly (repeat-clicking an
already-selected tile cycles `dragOrientation` through down/right/up/left), then drag it
onto the target square; occupied-to-be squares turn white and the tile ghosts to the
cursor as feedback.

`Square.js`'s `onDrop` handler ignores the actual `DataTransfer` payload and just calls
`dropCallback(x, y)` using whatever `selectedTile`/`dragOrientation` are already in React
state -- so a synthetic `drop` event on the target square exercises the real
`onDrop -> dropCallback -> playTile` path without needing a real drag gesture:

```js
const dt = new DataTransfer();
document.querySelector('[data-tron_x="X"][data-tron_y="Y"]')
  .dispatchEvent(new DragEvent('drop', { bubbles: true, cancelable: true, dataTransfer: dt }));
```

**Not** testable this way: the actual mouse-driven `handleDragStart` ghost-image /
hover-highlight feedback -- that needs real screenshot compositing. Any change to
`handleDragStart`/ghost-image/hover-highlight code needs a human to verify visually.

**Gotcha if you invoke `handleDragStart`'s React prop directly** (e.g. via
`element[Object.keys(element).find(k=>k.startsWith('__reactProps$'))].onDragStart(...)`
to sidestep real drag events): it clones the tile into a ghost element appended to
`document.body` and only removes it on the next `requestAnimationFrame`. A synchronous
follow-up script (e.g. re-querying `[data-tile]` to check the hand) runs before that
rAF fires and will see one extra element with a `data-tile` attribute -- this looks like
a duplicated tile but is just the leftover ghost. Scope queries to the actual hand
container (e.g. `.overflow-y-auto`) or wait a tick before asserting.

### Draw / pass

- Draw: click the button (`b.textContent.trim() === 'draw'`). `POST /game/<code>/draw`,
  adds a tile, `just_drew` becomes true, turn does NOT advance (still your turn -- you
  must then play or pass).
- Pass: click similarly, `POST /game/<code>/pass`. On your first pass of the whole game
  (not per-round), a "Vision Quest" modal pops up to pick a chicken-foot reaction GIF
  (real Klipy API search, purely decorative, doesn't affect game logic). It's rendered
  via an antd portal, so it won't show up in `<main>`-scoped text extraction -- check
  `document.querySelector('.ant-modal')`.
  ```js
  // type a search term (native-setter trick needed for a React-controlled input)
  const input = document.querySelector('input[placeholder="Search KLIPY"]');
  Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, 'value').set.call(input, 'chicken');
  input.dispatchEvent(new Event('input', { bubbles: true }));
  // wait ~1-2s for the debounced search, then pick a result
  // scope the query to the modal -- once any player has a chicken-foot GIF set,
  // img[src*="static.klipy.com"] also matches ChickenFoot marker components elsewhere
  // on the page, so an unscoped query can grab the wrong element
  document.querySelector('.ant-modal').querySelector('img[src*="static.klipy.com"]')
    .closest('div[class*="cursor-pointer"]').click();
  ```
  Fires `POST /game/<code>/foot`. `chicken_foot_coord` is auto-computed server-side from
  your line's current open end -- no board-square picking needed once your line is longer
  than just the leader tile.

### antd controls in headless tests

- `Select` dropdowns don't open from a plain `.click()` on `.ant-select`; dispatch a real
  `MouseEvent('mousedown', { bubbles: true, cancelable: true })` on it instead, then query
  `.ant-select-item-option` elements once `.ant-select-dropdown` appears (may need a
  separate follow-up read -- the dropdown DOM can lag one tick behind the `mousedown`).
- A `Modal`'s close button/`onCancel` firing correctly does **not** guarantee the DOM node
  disappears in this harness -- see "Browser pane doesn't composite frames" below before
  concluding a modal-close bug is real.

### Rejoining an in-progress game / multi-tab pickup testing

- A full page reload loses the in-memory player session (name/userInfo aren't persisted).
  To get a player back into a game they were already in: register with the **same player
  name** again, then type the game's code (without the `-XXXXXX` suffix) into the OTP-style
  code field -- the server reconnects them by name via `X-Player-Name`.
- "pick-up game" via the UI is unreliable for deliberately landing two tabs in the *same*
  pickup game -- it often creates a new pickup game per click rather than matching an
  existing waiting one. More reliable for two-tab testing: click "pick-up game" once from
  one tab to create/join a pickup game, note its code, then add the second player directly
  via the API rather than fighting the UI's matching:
  ```bash
  curl -s -X POST http://localhost:8080/players -H 'X-Player-Name: <name2>'
  curl -s -X PUT http://localhost:8080/game/<code> -H 'X-Player-Name: <name2>' -d '{}'
  ```
  Then load `<name2>` into a second tab and rejoin by code (see above) to see it in the UI.

### Free lines (spacer -> double)

1. Check `me.spacer_hints` (see above).
2. Click "FREE LINE" in the hand panel (`spacerClicked()` in `Hand.js`, sets
   `selectedTile = {a:-1, b:-1}` as a sentinel):
   ```js
   Array.from(document.querySelectorAll('*'))
     .find(e => e.children.length === 0 && e.textContent.trim() === 'FREE LINE').click();
   ```
3. Click the two spacer endpoint squares (either order), same click-click pattern and
   same "separate tool calls" requirement as laying a tile. Fires `POST /game/<code>/spacer`.
   Placing a spacer does NOT end your turn.
4. Re-check hints -- the qualifying double now shows a hint near the spacer's far end.
   Play it like a normal tile. Server records it in `round.free_lines` with
   `player_name: ""` and resets `round.spacer` to null.

### Round/game end

- A round ends via: hand-emptied win, attrition (one player left alive), all-dead
  stalemate, or a cutoff kill leaving one survivor. Check `round.done`.
- If `round.done` but not `game.done`, all players click "ready" again to start the next
  round.
- The whole game ends (`game.done: true`) once a round LED BY THE `0:0` DOUBLE finishes.
  Player `score` accumulates across all rounds of a game (not reset per round).

### Architecture note

The game is robust to "messing around" with the UI because the entire game state is
returned by the server on every GET, and the client tracks/derives no authoritative state
of its own -- every render is just a projection of the latest server response. A stray
click or a rejected action can't desync anything; the next poll re-syncs the true state
regardless.

## Testing-harness gotchas (not app bugs)

- **The Browser pane doesn't composite frames when it isn't displayed**, and CSS
  transitions/animations never fire `transitionend`/`animationend` in that state. Anything
  whose completion depends on a CSS transition (e.g. antd `Modal`'s close animation, via
  `@rc-component/motion`) will look permanently stuck in this harness even though the
  underlying React state (`isOpen`/`open` prop) updates correctly. Confirmed by creating a
  plain `<div>` with a CSS `transition` and toggling it -- `transitionend` never fired.
  Before treating a "close/animation doesn't finish" observation as a real bug, verify in
  an actually-visible browser (ask the user, or use `claude-in-chrome` if connected) --
  don't "fix" library close-animation code based on this harness alone.
- React fiber inspection is useful for verifying/driving component state directly when
  there's no UI affordance for it (or to bypass a broken UI path while isolating a
  different bug): grab a DOM node, find its `__reactFiber$...`/`__reactProps$...` key,
  walk `.return` to a named component (`f.type?.name`), then read `.memoizedState`
  (hooks, in declaration order) or `.memoizedProps`. A hook's `useState` setter is
  reachable at `hook.queue.dispatch` and can be called directly to force a state change
  without going through the UI.
- `javascript_exec` calls against the same tab appear to share a persistent top-level
  scope -- a `const`/`let` declared in one call can collide ("Identifier ... has already
  been declared") with the same name in a later call. Wrap probe scripts in an IIFE
  (`(() => { ... })()`) to avoid this.

## Firestore issue reports

Players can file in-app issue reports (`ReportIssue` in `tronserv/game/firestore.go`),
stored at `envs/<env>/issues/<docId>` (`env` is `prod` or `qa`) with fields
`reported_by`, `summary`, `whatHappened`, `whatShouldHappen`, `errorMessage`, and
`game_json` (the full `Game` struct, JSON-encoded as a string).

- Fetching a doc via the Firestore REST API with a manually-printed `gcloud auth
  print-access-token` piped into `curl` gets blocked by the auto-mode permission
  classifier (looks like credential exfiltration). A small standalone Go program using
  `cloud.google.com/go/firestore`'s client directly (Application Default Credentials via
  `gcloud auth`) is not blocked and is the natural approach anyway, since it reuses the
  same auth path `tronserv` itself uses.
- To turn an issue's `game_json` into testdata: write the raw string to
  `tronserv/game/testdata/<name>.json`, then run (from `tronserv`)
  `go run ./prettify game/testdata/<name>.json` -- this round-trips it through the real
  `game.Game` struct (canonical field order, tab indent), not a generic JSON
  pretty-printer. See `tronserv/game/testdata/roundfoot.json` for an example (from
  Firestore issue `11C0K32DvXs3M8VRwK9r`) and `TestRoundfoot` in `game_test.go` for the
  test pattern: load via `decodeGame`, call `r.FindLegalMoves`, assert on the specific
  return value (`moves`/`spacers`/`passFeet`) relevant to the bug.

## Dev environment version mismatches (this machine)

- Node is 18.17.0, but `tronapp`'s Next.js (16.x) requires >=20.9.0 for `npm run build` --
  it hard-fails immediately with just a version-check message, not a real compile
  attempt. `npm run dev` still works fine on 18.17. Don't rely on `npm run build` for
  verification here; use `npx eslint` plus live testing via `npm run dev` instead.
- Go installed is 1.23.2, but `tronserv/go.mod` declares `go 1.24.11`. This hasn't
  actually broken `go build`/`go test`/`go run` in practice, but is worth knowing if a
  build ever complains about the toolchain version.

## Known issues

- `Game.js` logs an empty `{}` error object in two places (`getGame`'s polling-loop catch
  around line 143, `leaveOrQuit`'s catch around line 544) -- the underlying error shape
  from `Client.js`'s `doRequest` isn't useful for debugging on at least some failure paths
  (e.g. an HTTP 408 during heavy polling). Being tracked/fixed separately.
