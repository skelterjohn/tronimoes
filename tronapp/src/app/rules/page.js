"use client";

import { useMemo, useState } from "react";
import RulesBoard from "@/app/components/rules/RulesBoard";
import { GameContext } from "@/app/components/GameState";
import { useGameState } from "@/app/components/GameState";
import Settings from "@/app/components/settings/Settings";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faGear } from "@fortawesome/free-solid-svg-icons";

function slugify(title) {
	return title.toLowerCase().replace(/\s+/g, "-").replace(/[^a-z0-9-]/g, "");
}

function Section({ title, children }) {
	const id = slugify(title);
	return (
		<div id={id} className="mx-auto px-6 py-5 max-w-2xl space-y-10 scroll-mt-24">
			<h2 className="text-2xl font-bold tracking-tight text-white border-b border-slate-500 pb-2 mb-1">{title}</h2>
			{children}
		</div>
	);
}

const SECTIONS = [
	{
		title: "Leaders and lines",
		content: (
			<>
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
			</>
		),
	},
	{
		title: "A game of murder",
		content: (
			<>
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
			</>
		),
	},
	{
		title: "A dead line",
		content: (
			<>
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
			</>
		),
	},
	{
		title: "Drawing and passing",
		content: (
			<>
				<p>
					On your turn, you may opt to draw a tile, if you haven't already.
				</p>
				<p>
					Once you've drawn a tile, you may opt to pass your turn.
				</p>
				<p>
					You can do this draw/pass maneuver whether or not you have a tile
					that could have been laid.
				</p>
			</>
		),
	},
	{
		title: "The dreaded chicken-foot",
		content: (
			<>
				<p>
					If you pass, and you weren't chicken-footed already, you become
					chicken-footed.
				</p>
				<p>
					When you are chicken-footed, other players may lay tiles on your
					line, and you may only play on your own line.
				</p>
				<p>
					Other players are likely to do bad things to you. If they have your
					line box itself in, you're dead and they get the credit (it's the
					player who laid the tile that made it so you couldn't continue that
					gets points for the kill).
				</p>
				<p>
					Once you are finally able to play a tile on your own line again, you
					are no longer chicken-footed.
				</p>
			</>
		),
	},
	{
		title: "A vision quest",
		content: (
			<>
				<p>
					The first time you become chicken-footed, you must complete a vision quest.
				</p>
				<p>
					On tronimoes.com, that means doing a quick image search for something that
					makes you feel like you're chicken-footed. The resulting image is your chicken-foot.
				</p>
				<p>
					When you are chicken-footed, your chicken-foot is displayed on the end of your
					line, making it clear to you and your opponents that you are chicken-footed.
				</p>
			</>
		),
	}
];

export default function RulesPage() {
	const [showSettingsModal, setShowSettingsModal] = useState(false);
	const gameState = useGameState();
	const stateWithConfig = useMemo(
		() => ({ ...gameState, config: gameState?.config ?? { tileset: "beehive" } }),
		[gameState]
	);

	return (
		<GameContext.Provider value={stateWithConfig}>
			<main className="flex flex-col h-screen w-full bg-slate-800 text-slate-100 overflow-hidden">
				<header className="sticky top-0 z-10 flex items-center justify-between px-6 py-4 bg-slate-800 border-b border-slate-600">
					<div className="flex items-baseline gap-4">
						<h1 className="text-3xl font-bold tracking-tight">Tronimoes: the rules of play</h1>
						<a
							href="https://docs.google.com/document/d/e/2PACX-1vQJivrpZZ14fF60BqqVWWtj5D_3ZH3b-1KU42FXMevsrVjC034QxnRc0a7pYraCnQ-vuYdjmrm9OT8A/pub"
							className="cursor-pointer underline underline-offset-2 text-slate-400 hover:text-slate-200 text-sm"
						>
							original google doc
						</a>
					</div>
					<button
						type="button"
						onClick={() => setShowSettingsModal(true)}
						className="text-slate-300 hover:text-white cursor-pointer p-1"
						aria-label="Settings"
					>
						<FontAwesomeIcon icon={faGear} className="text-xl" />
					</button>
				</header>
				<div className="flex flex-1 min-h-0">
					<aside
						aria-label="Table of contents"
						className="w-52 flex-shrink-0 overflow-y-auto border-r border-slate-600 bg-slate-800 py-4"
					>
						<nav className="px-4 space-y-1">
							{SECTIONS.map((section) => (
								<a
									key={section.title}
									href={`#${slugify(section.title)}`}
									className="block py-1.5 text-sm text-slate-300 hover:text-slate-100 focus:text-slate-100 focus:outline-none underline-offset-2 hover:underline"
								>
									{section.title}
								</a>
							))}
						</nav>
					</aside>
					<div className="flex-1 overflow-y-auto">
						{SECTIONS.map((section) => (
							<Section key={section.title} title={section.title}>
								{section.content}
							</Section>
						))}
					</div>
				</div>
				<Settings
					isOpen={showSettingsModal}
					onClose={() => setShowSettingsModal(false)}
				/>
			</main>
		</GameContext.Provider>
	);
}
