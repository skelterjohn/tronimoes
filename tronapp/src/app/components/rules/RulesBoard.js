"use client";

import { useRef, useEffect } from "react";
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
 * @param {string} [className] - Optional extra classes for the board container
 */
export default function RulesBoard({ height, tiles, roundLeader, lineHeads, activePlayer, className = "" }) {
	const width = height - 1;
	const containerRef = useRef(null);

	useEffect(() => {
		const el = containerRef.current;
		if (!el) return;
		const handleWheel = (e) => {
			e.preventDefault();
			e.stopPropagation();
		};
		el.addEventListener("wheel", handleWheel, { passive: false, capture: true });
		return () => el.removeEventListener("wheel", handleWheel, { capture: true });
	}, []);

	return (
		<div
			ref={containerRef}
			className={`w-full max-w-full mx-auto aspect-square ${className}`.trim()}
		>
			<TipProvider>
				<Board
					width={width}
					height={height}
					tiles={tiles}
					spacer={undefined}
					lineHeads={lineHeads}
					roundLeader={roundLeader}
					freeLeaders={new Set()}
					selectedTile={undefined}
					playTile={noop}
					playSpacer={noop}
					clearSpacer={noop}
					chickenFeet={{}}
					chickenFeetURLs={{}}
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
				/>
			</TipProvider>
		</div>
	);
}
