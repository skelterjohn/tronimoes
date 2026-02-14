"use client";

import { useMemo, useState, useRef, useEffect } from "react";
import Board from "@/app/components/board/Board";
import { TipProvider } from "@/app/components/tutorial/InnerTip";
import { GameContext } from "@/app/components/GameState";
import { useGameState } from "@/app/components/GameState";
import Settings from "@/app/components/settings/Settings";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faGear } from "@fortawesome/free-solid-svg-icons";

const RULES_BOARD_HEIGHT = 7;
const RULES_BOARD_WIDTH = RULES_BOARD_HEIGHT - 1;

const rulesTiles = {
	"2,3": { a: 3, b: 3, orientation: "right", color: "white", dead: false },
};

const rulesRoundLeader = { pips_a: 3, pips_b: 3 };

const rulesLineHeads = [
	{ tile: { pips_a: 3, pips_b: 3 }, coord: { x: 3, y: 2 } },
];

const noop = () => {};

export default function RulesPage() {
	const [showSettingsModal, setShowSettingsModal] = useState(false);
	const boardContainerRef = useRef(null);
	const gameState = useGameState();
	const stateWithConfig = useMemo(
		() => ({ ...gameState, config: gameState?.config ?? { tileset: "beehive" } }),
		[gameState]
	);

	useEffect(() => {
		const el = boardContainerRef.current;
		if (!el) return;
		const handleWheel = (e) => {
			e.preventDefault();
			e.stopPropagation();
		};
		el.addEventListener("wheel", handleWheel, { passive: false, capture: true });
		return () => el.removeEventListener("wheel", handleWheel, { capture: true });
	}, []);

	return (
		<main className="min-h-screen w-full bg-slate-800 text-slate-100">
			<header className="sticky top-0 z-10 flex items-center justify-between px-6 py-4 bg-slate-800 border-b border-slate-600">
				<h1 className="text-3xl font-bold tracking-tight">game rules</h1>
				<button
					type="button"
					onClick={() => setShowSettingsModal(true)}
					className="text-slate-300 hover:text-white cursor-pointer p-1"
					aria-label="Settings"
				>
					<FontAwesomeIcon icon={faGear} className="text-xl" />
				</button>
			</header>
			<div className="mx-auto px-6 py-10">

				<div
					ref={boardContainerRef}
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
			<GameContext.Provider value={stateWithConfig}>
				<Settings
					isOpen={showSettingsModal}
					onClose={() => setShowSettingsModal(false)}
				/>
			</GameContext.Provider>
		</main>
	);
}
