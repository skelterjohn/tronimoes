<!-- BEGIN:nextjs-agent-rules -->

# This is NOT the Next.js you know

This version has breaking changes — APIs, conventions, and file structure may all differ from your training data. Read the relevant guide in `node_modules/next/dist/docs/` (resolved from this file's directory; in monorepos the `next` package may not be visible from the repo root) before writing any code. Heed deprecation notices.

This block is written and re-added by `next dev` — verify at `node_modules/next/dist/server/lib/generate-agent-files.js`. Removing it from a diff only re-creates the uncommitted change; committing it with your work keeps the tree clean.

<!-- END:nextjs-agent-rules -->

# react-hooks/set-state-in-effect: remediation patterns

This repo's `eslint-config-next` flags any `setState` called synchronously inside a
`useEffect` body (not inside a `.then()`/`.catch()`/other async callback -- those are
fine). Fixing ~50 of these across the frontend surfaced two real patterns and one sharp
gotcha; read this before "fixing" another one.

## Pick the right pattern

1. **Pure derivation** -- the effect *only* computes `X` from props/other state, with no
   other mutation path for `X` (grep for every `setX(` call in the file first). Delete the
   `useState`/`useEffect` entirely; compute `X` as a plain `const`/`let` during render (or
   via `useMemo` if it's expensive, or if `X` is itself a dependency of another
   `useCallback`/`useMemo` -- otherwise that dependent hook gets a new reference every
   render and ESLint will flag *that* as unstable). Example: `Board.js`, `Tile.js`,
   `Square.js`, most of `Opponent.js`.
2. **Sync-then-diverge** ("adjust state during render", React's own documented pattern) --
   the state is *also* set by something else (a click handler, an async response, a
   drag/reorder callback, another component via a passed-down setter), so the effect's job
   is only to re-initialize/re-sync it when some *other* value changes, not to own it
   fully. Track the previous value of the trigger in its own `useState`; compare during
   render; call `setState` directly in the render body (not inside `useEffect`) when it
   changed. Example: `Options.js`'s `roundOut` (re-inits from shared `options` on open, but
   the `Select` also sets it directly), `Hand.js`'s `handOrder` (re-synced from
   `player.hand`, but drag-reorder also sets it directly).

Effects with real side effects beyond `setState` (event listeners, `ResizeObserver`,
`audio.play()`, network calls, `router.push`) are not this rule's target and are correctly
left alone -- ESLint only flags a *synchronous* `setState` call in the effect body. If an
effect both plays audio/does another side effect *and* calls `setState` for pure
bookkeeping (e.g. "was the player count higher or lower than last time"), don't try to
force the whole thing into render -- convert just the bookkeeping value to a `useRef`
(mutated inside the still-effect-housed side effect) instead of `useState`. See
`Game.js`'s `lastPlayerCountRef` next to the pre-existing `prevTilesInBagRef` for the
established pattern.

## The sentinel gotcha (caused a real bug -- read this)

When seeding the "previous value" tracker for pattern 2, **do not** initialize it from the
current prop/state value:

```js
// WRONG -- looks reasonable, silently breaks first-mount sync:
const [prevPlayer, setPrevPlayer] = useState(player);
if (player !== prevPlayer) { setPrevPlayer(player); syncSomethingFrom(player); }
```

An effect with `useEffect(() => { syncSomethingFrom(player); }, [player])` **always** runs
once after the initial mount, regardless of what `player` was "before." But the pattern
above treats "the value I already have" as already-synced from render 1 onward, so if the
component ever mounts with the sync already needed (e.g. `player.hand` already populated,
which is the *normal* case for `Hand.js` -- it doesn't mount until game data has already
loaded), the first sync silently never happens and the derived state stays at its empty
initial value forever.

Fix: seed the tracker with a sentinel that can never equal a real value, so the sync is
guaranteed to run at least once:

```js
const NEVER_SYNCED = Symbol('never-synced'); // module scope
const [prevPlayer, setPrevPlayer] = useState(NEVER_SYNCED);
if (player !== prevPlayer) { setPrevPlayer(player); syncSomethingFrom(player); }
```

This exact mistake shipped a real bug in `Hand.js` (hand rendered empty on mount) that
only showed up under live testing, not lint/build -- caught before commit, but it's why
`Game.js`'s many render-adjustment blocks all seed from a shared `NEVER_SYNCED` constant
rather than the current value. When you *can* prove the trigger value is guaranteed to
start in a "no sync needed" state (e.g. a modal's `isOpen` that always mounts `false`
because the modal component is unconditionally rendered but the boolean starts `false`;
or a Firebase `useAuthState()` hook that's documented to start `loading: true`
synchronously, so nothing needs to happen on that first tick anyway), seeding from the
current value is fine -- but that's a fact about the specific data flow, not something to
assume by default.

## SSR guard when a client-only effect becomes a `useMemo`

If the effect being replaced only ran client-side for a reason (reads `window`, DOM APIs,
etc.), converting it to a bare `useMemo`/render-time computation breaks SSR, since render
(unlike effects) also executes on the server. `GameState.js`'s `client` derivation
(`clientFor()` reads `window.location`) needs an explicit `typeof window === 'undefined'`
guard inside the `useMemo` for exactly this reason -- caught via a real SSR 500 on first
page load, not lint.

## antd 6.x internals (relevant if you ever touch Modal/Select/etc.)

This version vendors its dialog/motion/trigger implementations as scoped
`@rc-component/*` packages (`@rc-component/dialog`, `@rc-component/motion`,
`@rc-component/trigger`, ...) rather than the old standalone `rc-dialog`/`rc-motion`/etc.
packages from antd 5. `Modal`'s close sequence depends on a CSS transition actually
completing (`@rc-component/motion`'s `CSSMotion` fires `onVisibleChanged` off a real
`transitionend`/`animationend` event) before the dialog unmounts from the DOM -- see the
root `CLAUDE.md`'s "Testing-harness gotchas" section if this looks broken while testing
headlessly.
