"use client";

import Board from "@/app/components/board/Board";
import { TipProvider } from "@/app/components/tutorial/InnerTip";

const noop = () => {};

/**
 * Renders a non-interactive rules board with fixed tiles. Use inside GameContext.Provider.
 * @param {number} height - Board row count (width becomes height - 1)
 * @param {Record<string, { a: number, b: number, orientation: string, color?: string, dead?: boolean }>} tiles - Map of "x,y" to tile data
 * @param {{ pips_a: number, pips_b: number }} roundLeader - The round leader tile
 * @param {{ tile: { pips_a: number, pips_b: number }, coord: { x: number, y: number } }[]} lineHeads - Line head entries
 * @param {{ color: string }} [activePlayer] - Optional; if set, gutter uses this player color (e.g. { color: "red" }). Omit for completed-round look (black gutter).
 * @param {Record<string, string>} [chickenFeet] - Optional; map of "x,y" to color for chicken-foot markers
 * @param {Record<string, string>} [chickenFeetURLs] - Optional; map of "x,y" to image URL for chicken-foot markers
 * @param {{ a: { x: number, y: number }, b?: { x: number, y: number } }} [spacer] - Optional; spacer placement (Board shows Spacer at spacer.a)
 * @param {string} [className] - Optional extra classes for the board container
 */
export default function RulesBoard({ height, tiles, roundLeader, lineHeads, activePlayer, chickenFeet = {}, chickenFeetURLs = {}, spacer, className = "" }) {
	const width = height - 1;

	return (
		<div
			className={`w-full max-w-full mx-auto aspect-square ${className}`.trim()}
		>
			<TipProvider>
				<Board
					width={width}
					height={height}
					tiles={tiles}
					spacer={spacer}
					lineHeads={lineHeads}
					roundLeader={roundLeader}
					freeLeaders={new Set()}
					selectedTile={undefined}
					playTile={noop}
					playSpacer={noop}
					clearSpacer={noop}
					chickenFeet={chickenFeet}
					chickenFeetURLs={chickenFeetURLs}
					indicated={undefined}
					setIndicated={noop}
					activePlayer={activePlayer}
					hints={{}}
					playA={undefined}
					setPlayA={noop}
					spacerHints={undefined}
					hoveredSquares={new Set()}
					setMouseIsOver={noop}
					dropCallback={noop}
					setSquareSpan={noop}
					interactive={false}
				/>
			</TipProvider>
		</div>
	);
}
