"use client";

import { useMemo } from "react";
import Board from "@/app/components/board/Board";
import { TipProvider } from "@/app/components/tutorial/InnerTip";
import { GameContext } from "@/app/components/GameState";
import { useGameState } from "@/app/components/GameState";

const RULES_BOARD_WIDTH = 6;
const RULES_BOARD_HEIGHT = 7;

const rulesTiles = {
	"3,2": { a: 3, b: 3, orientation: "down", color: "white", dead: false },
	"4,2": { a: 3, b: 1, orientation: "right", color: "white", dead: false },
	"5,2": { a: 1, b: 4, orientation: "right", color: "white", dead: false },
	"2,2": { a: 3, b: 6, orientation: "left", color: "white", dead: false },
	"1,2": { a: 6, b: 0, orientation: "left", color: "white", dead: false },
	"3,3": { a: 3, b: 2, orientation: "down", color: "white", dead: false },
	"3,4": { a: 2, b: 5, orientation: "down", color: "white", dead: false },
};

const rulesRoundLeader = { pips_a: 3, pips_b: 3 };

const rulesLineHeads = [
	{ tile: { pips_a: 1, pips_b: 4 }, coord: { x: 5, y: 2 } },
	{ tile: { pips_a: 2, pips_b: 5 }, coord: { x: 3, y: 4 } },
	{ tile: { pips_a: 6, pips_b: 0 }, coord: { x: 1, y: 2 } },
];

const noop = () => {};

export default function RulesPage() {
	const gameState = useGameState();
	const stateWithConfig = useMemo(
		() => ({ ...gameState, config: gameState?.config ?? { tileset: "beehive" } }),
		[gameState]
	);

	return (
		<main className="min-h-screen w-full bg-slate-800 text-slate-100">
			<div className="mx-auto px-6 py-10">
				<header className="mb-8">
					<h1 className="text-3xl font-bold tracking-tight">game rules</h1>
				</header>

				<div
					className="w-[75vw] mx-auto aspect-square"
					style={{ maxHeight: "75vw" }}
				>
					<GameContext.Provider value={stateWithConfig}>
						<TipProvider>
							<Board
						width={RULES_BOARD_WIDTH}
						height={RULES_BOARD_HEIGHT}
						tiles={rulesTiles}
						spacer={undefined}
						lineHeads={rulesLineHeads}
						roundLeader={rulesRoundLeader}
						freeLeaders={new Set()}
						selectedTile={undefined}
						playTile={noop}
						playSpacer={noop}
						clearSpacer={noop}
						chickenFeet={{}}
						chickenFeetURLs={{}}
						indicated={undefined}
						setIndicated={noop}
						activePlayer={undefined}
						hints={{}}
						playA={undefined}
						setPlayA={noop}
						spacerHints={undefined}
						hoveredSquares={new Set()}
						setMouseIsOver={noop}
						dropCallback={noop}
						setSquareSpan={noop}
					/>
						</TipProvider>
					</GameContext.Provider>
				</div>
			</div>
		</main>
	);
}
