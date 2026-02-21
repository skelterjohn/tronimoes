"use client";

import { useMemo, useState } from "react";
import RulesBoard from "@/app/components/rules/RulesBoard";
import { GameContext } from "@/app/components/GameState";
import { useGameState } from "@/app/components/GameState";
import Settings from "@/app/components/settings/Settings";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faGear, faList } from "@fortawesome/free-solid-svg-icons";
import VisionQuest from "@/app/components/visionquest/VisionQuest";

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
		title: "leaders and lines",
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
					The goal is to win the round, which can be accomplished in two ways: 
					being first to run out of tiles, or being the last player standing.
				</p>
				<RulesBoard
					height={7}
					tiles={{
						"2,3": { a: 3, b: 3, orientation: "right", color: "white", dead: false },
						"2,2": { a: 3, b: 5, orientation: "up", color: "red", dead: false },
						"4,3": { a: 3, b: 7, orientation: "up", color: "blue", dead: false },
						"2,0": { a:5, b:1, orientation: "left", color: "red", dead: false },
						"4,1": { a: 7, b:0, orientation: "up", color: "blue", dead: false, last: true },
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
		title: "leader election",
		content: (
			<>
				<p>
					The round leader, or the tile at the center of the board,
					is chosen and placed automatically.
				</p>
				<p>
					It is the highest double available from a player's hand that
					is lower than all previous round leaders.
				</p>
				<p>
					If no one has such a double, every player draws a tile until
					the double is found. This tile-drawing is why in some rounds, especially
					later ones, players may begin the round with a higher-than-expected
					number of tiles.
				</p>
				<p>
					The last round is the one that is led by the double-zero.
				</p>
			</>
		),
	},
	{
		title: "scoring and victory",
		content: (
			<>
				<p>
					The winner of a round gets 2 points.
				</p>
				<p>
					If you kill another player (more on this soon), you get 1 point.
				</p>
				<p>
					If you are killed, you lose 1 point.
				</p>
				<p>
					(If you kill your own line, it's net zero points.)
				</p>
				<p>
					After the last round is played, the player with the most points wins.
				</p>
			</>
		),
	},
	{
		title: "a game of murder",
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
						"2,0": { a: 5, b:3, orientation: "right", color: "red", dead: false, last: true },
					}}
					roundLeader={{ pips_a: 3, pips_b: 3 }}
					lineHeads={[
						{ tile: { pips_a: 3, pips_b: 3 }, coord: { x: 2, y: 3 } },
					]}
					activePlayer={{ color: "blue" }}
				/>
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
						"4,1": { a: 7, b:0, orientation: "up", color: "blue", dead: false, last: true },
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
		title: "drawing and passing",
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
		title: "the dreaded chicken-foot",
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
					player who laid the tile that made it so you couldn't continue who
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
		title: "a vision quest",
		content: (openVisionQuest, chickenFoot) => (
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
				<p>
					<button
						type="button"
						onClick={openVisionQuest}
						className="rounded bg-slate-600 px-4 py-2 text-slate-100 hover:bg-slate-500 focus:outline-none focus:ring-2 focus:ring-slate-400"
					>
						Go on a vision quest
					</button>
				</p>
				<p>You're red, and it's blue's turn. What do you think they might do?</p>
				<RulesBoard
					height={7}
					tiles={{
						"2,3": { a: 3, b: 3, orientation: "right", color: "white", dead: false },
						"2,2": { a: 3, b: 5, orientation: "up", color: "red", dead: false },
						"4,3": { a: 3, b: 7, orientation: "up", color: "blue", dead: false },
						"4,1": { a: 7, b:0, orientation: "up", color: "blue", dead: false, last: true },
					}}
					roundLeader={{ pips_a: 3, pips_b: 3 }}
					lineHeads={[
						{ tile: { pips_a: 3, pips_b: 3 }, coord: { x: 2, y: 3 } },
					]}
					activePlayer={{ color: "blue" }}
					chickenFeet={{ "2,2": "red" }}
					chickenFeetURLs={{ "2,2": chickenFoot }}
				/>
			</>
		),
		contentIsFunction: true,
	},
	{
		title: "your own worst enemy",
		content: (openVisionQuest, chickenFoot) => (
			<>
				<p>
					In this case, they backed you into a corner, which killed your line
					and ended the round.
				</p>
				<RulesBoard
					height={7}
					tiles={{
						"2,3": { a: 3, b: 3, orientation: "right", color: "white", dead: false },
						"2,2": { a: 3, b: 5, orientation: "up", color: "red", dead: true },
						"4,3": { a: 3, b: 7, orientation: "up", color: "blue", dead: false },
						"4,1": { a: 7, b:0, orientation: "up", color: "blue", dead: false },
						"3,1": { a: 5, b: 1, orientation: "down", color: "red", dead: true, last: true },
					}}
					roundLeader={{ pips_a: 3, pips_b: 3 }}
					lineHeads={[
						{ tile: { pips_a: 3, pips_b: 3 }, coord: { x: 2, y: 3 } },
					]}
					activePlayer={false}
				/>
			</>
		),
		contentIsFunction: true,
	},
	{
		title: "starting chicken-footed",
		content: (openVisionQuest, chickenFoot) => (
			<>
				<p>
					If you find yourself in the unfortunate position of being unable to
					play off the round leader, you still become chicken-footed but must
					choose a square, adjacent to the round leader, for your foot.
				</p>
				<p>
					You choose this square by clicking on it before you pass.
				</p>
				<p>
					When you are finally able to play, it must go through this chosen square.
					If someone attempts to play on your line, it must go through this square.
				</p>
				<p>
					It is illegal for another player to lay a tile in such a way that you can 
					never play the first tile of your line. However, they may leave a box for you
					which results in line-death immediately once the first tile is played.
				</p>
				<RulesBoard
					height={7}
					tiles={{
						"2,3": { a: 3, b: 3, orientation: "right", color: "white", dead: false },
						"4,3": { a: 3, b: 7, orientation: "up", color: "blue", dead: false, last: true },
					}}
					roundLeader={{ pips_a: 3, pips_b: 3 }}
					lineHeads={[
						{ tile: { pips_a: 3, pips_b: 3 }, coord: { x: 2, y: 3 } },
					]}
					activePlayer={{ color: "red" }}
					chickenFeet={{ "2,2": "red" }}
					chickenFeetURLs={{ "2,2": chickenFoot }}
				/>
			</>
		),
		contentIsFunction: true,
	},
	{
		title: "double tile, double turn",
		content: (
			<>
				<p>
					If you lay a double tile, you get to go again.
				</p>
				<RulesBoard
					height={7}
					tiles={{
						"2,3": { a: 3, b: 3, orientation: "right", color: "white", dead: false },
						"2,2": { a: 3, b: 5, orientation: "up", color: "red", dead: false },
						"4,3": { a: 3, b: 7, orientation: "up", color: "blue", dead: false, last: true },
					}}
					roundLeader={{ pips_a: 3, pips_b: 3 }}
					lineHeads={[
						{ tile: { pips_a: 3, pips_b: 3 }, coord: { x: 2, y: 3 } },
					]}
					activePlayer={{ color: "red" }}
				/>
				<p>If red has the double-5...</p>
				<RulesBoard
					height={7}
					tiles={{
						"2,3": { a: 3, b: 3, orientation: "right", color: "white", dead: false },
						"2,2": { a: 3, b: 5, orientation: "up", color: "red", dead: false },
						"4,3": { a: 3, b: 7, orientation: "up", color: "blue", dead: false },
						"3,1": { a: 5, b: 5, orientation: "right", color: "red", dead: false, last: true },
					}}
					roundLeader={{ pips_a: 3, pips_b: 3 }}
					lineHeads={[
						{ tile: { pips_a: 3, pips_b: 3 }, coord: { x: 2, y: 3 } },
					]}
					activePlayer={{ color: "red" }}
				/>
				<p>Red gets one more turn for the sniper-rifle kill.</p>
				<RulesBoard
					height={7}
					tiles={{
						"2,3": { a: 3, b: 3, orientation: "right", color: "white", dead: false },
						"2,2": { a: 3, b: 5, orientation: "up", color: "red", dead: false },
						"4,3": { a: 3, b: 7, orientation: "up", color: "blue", dead: true },
						"3,1": { a: 5, b: 5, orientation: "right", color: "red", dead: false },
						"5,1": { a: 5, b: 2, orientation: "down", color: "red", dead: false, last: true },
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
		title: "dangling doubles",
		content: (
			<>
				<p>
					When you play off a double, you can choose either side. Both of
					these examples are legal plays (even though they missed the sniper-rifle kill).
				</p>
				<RulesBoard
					height={7}
					tiles={{
						"2,3": { a: 3, b: 3, orientation: "right", color: "white", dead: false },
						"2,2": { a: 3, b: 5, orientation: "up", color: "red", dead: false },
						"4,3": { a: 3, b: 7, orientation: "up", color: "blue", dead: false },
						"3,1": { a: 5, b: 5, orientation: "up", color: "red", dead: false },
						"4,1": { a: 5, b: 2, orientation: "right", color: "red", dead: false, last: true },
					}}
					roundLeader={{ pips_a: 3, pips_b: 3 }}
					lineHeads={[
						{ tile: { pips_a: 3, pips_b: 3 }, coord: { x: 2, y: 3 } },
					]}
					activePlayer={{ color: "blue" }}
				/>
				<RulesBoard
					height={7}
					tiles={{
						"2,3": { a: 3, b: 3, orientation: "right", color: "white", dead: false },
						"2,2": { a: 3, b: 5, orientation: "up", color: "red", dead: false },
						"4,3": { a: 3, b: 7, orientation: "up", color: "blue", dead: false },
						"3,1": { a: 5, b: 5, orientation: "up", color: "red", dead: false },
						"4,0": { a: 5, b: 2, orientation: "right", color: "red", dead: false, last: true },
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
		title: "free lines",
		content: (
			<>
				<p>
					If you have a double that is higher than any leader on the board, you
					can use it to start a "free line".
				</p>
				<p>
					The free-line spacer is a special non-played tile that helps determine
					where a new free line can be started. 
				</p>
				<p>
					The spacer is six squares long, and can be placed beginning on the end of
					any living line on the board. The leader is then placed off the other end of 
					the spacer. Then the spacer is removed from the board and does not impact
					future turns. 
				</p>
				<p>
					Later free lines must have a leader whose pips are higher than the previous
					free line leaders.
				</p>
				<RulesBoard
					height={7}
					tiles={{
						"2,3": { a: 3, b: 3, orientation: "right", color: "white", dead: false },
						"2,2": { a: 3, b: 5, orientation: "up", color: "red", dead: false },
						"4,3": { a: 3, b: 7, orientation: "up", color: "blue", dead: false, last: true },
					}}
					roundLeader={{ pips_a: 3, pips_b: 3 }}
					lineHeads={[
						{ tile: { pips_a: 3, pips_b: 3 }, coord: { x: 2, y: 3 } },
					]}
					activePlayer={{ color: "red" }}
					spacer={{ a: { x: 1, y: 1 }, b: { x: 1, y: 6 } }}
				/>
				<RulesBoard
					height={7}
					tiles={{
						"2,3": { a: 3, b: 3, orientation: "right", color: "white", dead: false },
						"2,2": { a: 3, b: 5, orientation: "up", color: "red", dead: false },
						"4,3": { a: 3, b: 7, orientation: "up", color: "blue", dead: false },
						"2,6": { a: 6, b: 6, orientation: "right", color: "white", dead: false, last: true },
					}}
					roundLeader={{ pips_a: 3, pips_b: 3 }}
					lineHeads={[
						{ tile: { pips_a: 3, pips_b: 3 }, coord: { x: 2, y: 3 } },
					]}
					activePlayer={{ color: "red" }}
				/>
				<p>
					The double-6 is the leader of a new free line.
				</p>
				<p>
					Any player who is not chicken-footed can play on the free line, but unlike
					the round leader, only a single line is created.
				</p>
				<p>
					Free lines can be used to kill player lines, and the player who used it to
					do so is the one who gets credit.
				</p>
				<p>
					Free lines can also be killed, but no one gets points (what kind of reprobate 
					would do that?).
				</p>
				<p>
					Upon starting a new free line, the player immediately goes again (since they
					just played a double).
				</p>
			</>
		),
	},
	{
		title: "moving a tile to the board",
		content: (
			<>
				<p>
					To play a tile, you first select it from your hand, which is next to
					the board. Your hand's background color is the same as your player color.
				</p>
				<p>
					There are two methods to move it to the board.
				</p>
				<p>
					First, you can click a tile to select it, then click the board square for
					the tile's "upper" pips, and click the board square for the tile's "lower" pips.
				</p>
				<p>
					Second, once a tile is selected, you can click it again to change its orientation.
					When the desired orientation is reached, the tile can be dragged into position
					on the board.
				</p>
				<p>
					These methods also work on mobile, for touch rather than click.
				</p>
			</>
		),
	},
	{
		title: "legal move hints",
		content: (
			<>
				<p>
					Playable tiles are raised slightly in your hand.
				</p>
				<p>
					Once you've selected a tile, playable squares are indicated on the board
					by a small white outline. Not all ways to place a tile on those squares are
					legal. But all legal ways to place a tile will fall within those squares.
				</p>
			</>
		),
	},
	{
		title: "placing a free line spacer",
		content: (
			<>
				<p>
					If you can begin a free line, the free line spacer will be slightly 
					raised in your hand. Select it by clicking.
				</p>
				<p>
					Then, chose one of the hinted squares adjacent to the current end of some
					existing line (not necessarily one you could play on otherwise).
				</p>
				<p>
					Then, choose another hinted square that is 6 spaces away from the first, with
					no laid tiles interposed.
				</p>
				<p>
					Once the spacer is in place, select the double you want to use to begin a new
					free line. The double must be placed on the end of the spacer opposite the
					existing line head.
				</p>
			</>
		),
	},
	{
		title: "organizing your hand",
		content: (
			<>
				<p>
					If you want to change the order of tiles in your hand (for instance, you
					want to track a list of linked tiles), they can be dragged within the hand.
				</p>
				<p>
					Simply click (or touch) a tile, and move it to the tile currently in the
					position you want. The rest of the tiles will be shifted to make room.
				</p>
			</>
		),
	},
	{
		title: "managing a large hand",
		content: (
			<>
				<p>
					If you have so many tiles that they don't all appear on the screen at once,
					you can scroll through them by clicking/touching the colored field behind the
					tiles and dragging.
				</p>
				<p>
					Alternatively, there are small up and down arrows to the right of the "react"
					button, which will scroll the hand up or down one row.
				</p>
			</>
		),
	},
	{
		title: "reactions",
		content: (
			<>
				<p>
					No one ever wants to read chat text written by their opponents in an online game.
				</p>
				<p>
					Instead, you can "react" by clicking the "react" button, and doing a quick
					image search for something that gets the message across.
				</p>
				<p>
					The reaction disappears on its own after 10 seconds. Alternatively, you can
					dismiss a reaction (for yourself only) by clicking on the image.
				</p>
			</>
		),
	},
	{
		title: "signing in and registering your name",
		content: (
			<>
				<p>
					If you want to be anonymous, you can stay signed out and enter an arbitrary
					designation and play. If you lose the page holding this game, you will not
					be able to rejoin.
				</p>
				<p>
					You can also sign in and choose "Anonymous", which will store your username and
					session in the current browser's cookies.
				</p>
				<p>
					Or, you can use Google or Facebook to manage your log-in.
				</p>
				<p>
					When signed in, your username is permanent and cannot be used by others.
				</p>
			</>
		),
	},
	{
		title: "starting and joining games",
		content: (
			<>
				<p>
					If you want to play a game against the first-available opponents, you can
					click "pick-up game". If there are no pick-up games available, this will
					create one.
				</p>
				<p>
					If you want to create a private game, enter a 6-character code of your choice.
					Other players that use this same 6-character code will join your game.
				</p>
				<p>
					Once a game is created, players can click the "ready" button. Once all players
					are ready, the game begins and no others may join.
				</p>
			</>
		),
	},
	{
		title: "observing games",
		content: (
			<>
				<p>
					If you receive the full URL of a game, including the 6-character code extension
					(beyond the initial 6-character code), you can visit that URL to observe the game
					without playing.
				</p>
			</>
		),
	},
	{
		title: "reporting issues",
		content: (
			<>
				<p>
					Buried in the settings menu, you can report issues with the game. When you report,
					a copy of the game is saved, along with your name and the information you provide.
				</p>
				<p>
					These reports are used to create tests that keep the game working properly.
				</p>
			</>
		),
	},
];

export default function RulesPage() {
	const [showSettingsModal, setShowSettingsModal] = useState(false);
	const [showVisionQuestModal, setShowVisionQuestModal] = useState(false);
	const [chickenFoot, setChickenFoot] = useState(null);
	const [showToc, setShowToc] = useState(false);
	const gameState = useGameState();
	const stateWithConfig = useMemo(
		() => ({ ...gameState, config: gameState?.config ?? { tileset: "classic" } }),
		[gameState]
	);

	return (
		<GameContext.Provider value={stateWithConfig}>
			<main className="font-game flex flex-col h-screen w-full bg-slate-800 text-slate-100 overflow-hidden">
				<header className="sticky top-0 z-10 flex items-center justify-between px-6 py-4 bg-slate-800 border-b border-slate-600">
					<div className="flex items-baseline gap-4">
						<h1 className="text-3xl font-bold tracking-tight">how to tronimo</h1>
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
				<div className="relative flex-1 min-h-0">
					{showToc && (
						<aside
							aria-label="Table of contents"
							className="absolute left-0 top-0 bottom-0 w-52 z-20 flex flex-col overflow-hidden border-r border-slate-600/70 bg-slate-800/80 shadow-lg"
						>
							<button
								type="button"
								onClick={() => setShowToc(false)}
								className="flex-shrink-0 p-2 m-2 self-start cursor-pointer rounded text-slate-100 bg-slate-700/90 hover:bg-slate-600/90"
								aria-label="Hide table of contents"
								title="Hide contents"
							>
								<FontAwesomeIcon icon={faList} className="text-xl" />
							</button>
							<nav className="flex-1 overflow-y-auto px-4 pb-4 space-y-1">
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
					)}
					<div className="relative w-full h-full overflow-y-auto">
						<button
							type="button"
							onClick={() => setShowToc(true)}
							className={`sticky top-0 left-0 z-10 mt-2 ml-2 cursor-pointer rounded p-2 text-slate-400 hover:text-slate-200 bg-slate-700/70 hover:bg-slate-600/90 ${showToc ? 'invisible pointer-events-none' : ''}`}
							aria-label="Show table of contents"
							title="Show contents"
						>
							<FontAwesomeIcon icon={faList} className="text-xl" />
						</button>
						{SECTIONS.map((section) => (
							<Section key={section.title} title={section.title}>
								{section.contentIsFunction
									? section.content(() => setShowVisionQuestModal(true), chickenFoot)
									: section.content}
							</Section>
						))}
					</div>
				</div>
				{showVisionQuestModal && (
					<VisionQuest
						title="Choose an image for your vision quest"
						isOpen={showVisionQuestModal}
						onClose={() => setShowVisionQuestModal(false)}
						setURL={setChickenFoot}
					/>
				)}
				<Settings
					isOpen={showSettingsModal}
					onClose={() => setShowSettingsModal(false)}
				/>
			</main>
		</GameContext.Provider>
	);
}
