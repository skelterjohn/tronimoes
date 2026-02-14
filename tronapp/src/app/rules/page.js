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
						A game of murder
					</h2>
					<p>Tronimoes is a game with tiles.</p>
					<RulesBoard
						height={7}
						tiles={{
							"2,3": { a: 3, b: 3, orientation: "right", color: "white", dead: false },
						}}
						roundLeader={{ pips_a: 3, pips_b: 3 }}
						lineHeads={[
							{ tile: { pips_a: 3, pips_b: 3 }, coord: { x: 2, y: 3 } },
						]}
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
