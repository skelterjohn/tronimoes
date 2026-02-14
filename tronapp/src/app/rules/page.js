export const metadata = {
	title: "Rules | tronimoes",
	description: "How to play tronimoes",
};

export default function RulesPage() {
	return (
		<main className="min-h-screen w-full bg-slate-800 text-slate-100">
			<div className="max-w-2xl mx-auto px-6 py-10">
				<header className="mb-10">
					<h1 className="text-3xl font-bold tracking-tight">game rules</h1>
				</header>

				<div className="space-y-8 font-[var(--font-geist-sans)]">
					<section>
						<h2 className="text-xl font-semibold mb-2">Overview</h2>
						<p className="text-slate-300 leading-relaxed">
							Tronimoes is a dominoes-style game. Players take turns placing tiles on the board,
							starting from the round leader and building lines. Match tile ends by pip count.
						</p>
					</section>

					<section>
						<h2 className="text-xl font-semibold mb-2">Tiles</h2>
						<p className="text-slate-300 leading-relaxed">
							Each tile has two halves with pips (0–6). You can rotate a tile before placing it
							so the correct half lines up with the board. Click or tap a tile to select it and
							rotate; when it’s oriented correctly, drag it onto the board.
						</p>
					</section>

					<section>
						<h2 className="text-xl font-semibold mb-2">Round leader</h2>
						<p className="text-slate-300 leading-relaxed">
							One tile is the round leader. Your line must start adjacent to this tile. The game
							shows which tile is the round leader so you know where you’re allowed to start.
						</p>
					</section>

					<section>
						<h2 className="text-xl font-semibold mb-2">Lines and line heads</h2>
						<p className="text-slate-300 leading-relaxed">
							Lines grow from the round leader. The last tile on each line is the “line head.”
							Only the line head is playable: you place your next tile adjacent to it, matching
							pips. Raised (hinted) tiles in your hand are playable at the current line heads.
						</p>
					</section>

					<section>
						<h2 className="text-xl font-semibold mb-2">Chicken foot</h2>
						<p className="text-slate-300 leading-relaxed">
							When a double is played and creates a “chicken foot,” the player who played it is
							“footed.” That player can only play on their own line until the foot is covered.
							Other players can also play on that line. Once the foot is gone, normal play
							resumes.
						</p>
					</section>

					<section>
						<h2 className="text-xl font-semibold mb-2">Free lines</h2>
						<p className="text-slate-300 leading-relaxed">
							You can start a free line by placing a double next to a playable line, with the
							other end of the double five spaces away (a “spacer”). The double must be higher
							than any other leader on the board. Free lines can be played on by anyone.
						</p>
					</section>

					<section>
						<h2 className="text-xl font-semibold mb-2">Dead lines</h2>
						<p className="text-slate-300 leading-relaxed">
							If a line is blocked (no one can play on it), that line is dead. Tiles on dead
							lines are shown gray. If it was your line, you can still play on free lines or
							chicken-footed lines. The game indicates which tiles and lines are dead.
						</p>
					</section>

					<section>
						<h2 className="text-xl font-semibold mb-2">How to play</h2>
						<ul className="text-slate-300 leading-relaxed list-disc list-inside space-y-1">
							<li>Enter your name and join a game with a code, or start a pick-up game.</li>
							<li>Select a playable tile (raised) and rotate it so the correct half will match the board.</li>
							<li>Drag the tile onto a valid square next to a line head (or to start a free line).</li>
							<li>Play continues until no one can play; the round may end and a new round leader can be set.</li>
						</ul>
					</section>
				</div>
			</div>
		</main>
	);
}
