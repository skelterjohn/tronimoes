"use client";

import { useMemo, useState } from "react";
import RulesBoard from "@/app/components/rules/RulesBoard";
import { GameContext } from "@/app/components/GameState";
import { useGameState } from "@/app/components/GameState";
import Settings from "@/app/components/settings/Settings";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faGear } from "@fortawesome/free-solid-svg-icons";

export default function RulesPage() {
	const [showSettingsModal, setShowSettingsModal] = useState(false);
	const gameState = useGameState();
	const stateWithConfig = useMemo(
		() => ({ ...gameState, config: gameState?.config ?? { tileset: "beehive" } }),
		[gameState]
	);

	return (
		<GameContext.Provider value={stateWithConfig}>
			<main className="min-h-screen w-full bg-slate-800 text-slate-100">
				<header className="sticky top-0 z-10 flex items-center justify-between px-6 py-4 bg-slate-800 border-b border-slate-600">
					<h1 className="text-3xl font-bold tracking-tight">Tronimoes: the rules of play</h1>
					<button
						type="button"
						onClick={() => setShowSettingsModal(true)}
						className="text-slate-300 hover:text-white cursor-pointer p-1"
						aria-label="Settings"
					>
						<FontAwesomeIcon icon={faGear} className="text-xl" />
					</button>
				</header>
				<div className="mx-auto px-6 py-10 max-w-2xl space-y-10">
					<h2 className="text-xl font-semibold tracking-tight text-slate-200">
						Leaders and lines
					</h2>
					<p>
						In tronimoes, players take turns laying tiles on the board. They start
						from the central tile, known as the "round leader". Lines are built by
						matching pips from tile to tile.
					</p>
					<p>
						Here, there are two players: red and blue. They both have lines beginning
						from the 3:3 round leader. All players must begin their line from the
						round leader.
					</p>
					<p>
						The goal is to win the round, which can be accomplished in a two ways: 
						being first to run out of tiles, or being the last player standing.
					</p>
					<RulesBoard
						height={7}
						tiles={{
							"2,3": { a: 3, b: 3, orientation: "right", color: "white", dead: false },
							"2,2": { a: 3, b: 5, orientation: "up", color: "red", dead: false },
							"4,3": { a: 3, b: 7, orientation: "up", color: "blue", dead: false },
							"2,0": { a:5, b:1, orientation: "left", color: "red", dead: false },
							"4,1": { a: 7, b:0, orientation: "up", color: "blue", dead: false },
						}}
						roundLeader={{ pips_a: 3, pips_b: 3 }}
						lineHeads={[
							{ tile: { pips_a: 3, pips_b: 3 }, coord: { x: 2, y: 3 } },
						]}
						activePlayer={{ color: "red" }}
					/>
				</div>
				<div className="mx-auto px-6 py-10 max-w-2xl space-y-10">
					<h2 className="text-xl font-semibold tracking-tight text-slate-200">
						A game of murder
					</h2>
					<p>
						If a player lays a tile that makes it impossible for another player to
						continue their line, that player is "killed", and their line is "dead".
					</p>
					<p>
						Their tiles remain on the board, blocking others.
					</p>
					<p>What move can blue make to win this round?</p>

					<RulesBoard
						height={7}
						tiles={{
							"2,3": { a: 3, b: 3, orientation: "right", color: "white", dead: false },
							"2,2": { a: 3, b: 5, orientation: "up", color: "red", dead: false },
							"3,2": { a: 3, b: 7, orientation: "up", color: "blue", dead: false },
							"2,0": { a: 5, b:3, orientation: "right", color: "red", dead: false },
						}}
						roundLeader={{ pips_a: 3, pips_b: 3 }}
						lineHeads={[
							{ tile: { pips_a: 3, pips_b: 3 }, coord: { x: 2, y: 3 } },
						]}
						activePlayer={{ color: "blue" }}
					/>
				</div>
				<div className="mx-auto px-6 py-10 max-w-2xl space-y-10">
					<h2 className="text-xl font-semibold tracking-tight text-slate-200">
						A dead line
					</h2>
					<p>
						That's right (probably): blue can place their next tile blocking
						red from continuing.
					</p>
					<p>
						Unfortunately for red, the edges of the board block them on one side,
						and blue's tiles block the other. There is no way to play more tiles
						on that line, so it's dead.
					</p>
					<RulesBoard
						height={7}
						tiles={{
							"2,3": { a: 3, b: 3, orientation: "right", color: "white", dead: false },
							"2,2": { a: 3, b: 5, orientation: "up", color: "red", dead: true },
							"3,2": { a: 3, b: 7, orientation: "up", color: "blue", dead: false },
							"2,0": { a: 5, b:3, orientation: "right", color: "red", dead: true },
							"4,1": { a: 7, b:0, orientation: "up", color: "blue", dead: false },
						}}
						roundLeader={{ pips_a: 3, pips_b: 3 }}
						lineHeads={[
							{ tile: { pips_a: 3, pips_b: 3 }, coord: { x: 2, y: 3 } },
						]}
						activePlayer={false}
					/>
				</div>
				<Settings
					isOpen={showSettingsModal}
					onClose={() => setShowSettingsModal(false)}
				/>
			</main>
		</GameContext.Provider>
	);
}
