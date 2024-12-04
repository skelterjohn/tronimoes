"use client";

import { useState, useEffect } from 'react';

import Spacer from "./Spacer";
import Square from "./Square";
import Tile from "./Tile";
import ChickenFoot from './ChickenFoot';
import Hint from './Hint';

const bgColorMap = {
	red: "bg-red-500",
	blue: "bg-blue-500",
	green: "bg-green-500",
	yellow: "bg-yellow-500",
	orange: "bg-orange-500",
	fuchsia: "bg-fuchsia-500",
	white: "bg-white"
};

export default function Board({ 
		width = 10, height = 11,
		tiles, spacer,
		lineHeads,
		roundLeader, freeLeaders,
		selectedTile,
		playTile, playSpacer,
		chickenFeet, chickenFeetURLs,
		indicated, setIndicated,
		activePlayer, hints,
		playA, setPlayA,
		spacerHints, clearSpacer,
		hoveredSquares, setMouseIsOver,
		dropCallback
	}) {
	function rightClick(evt) {
		evt.preventDefault();
		setPlayA(undefined);
		setIndicated(undefined);
		clearSpacer();
	}

	useEffect(() => {
		setPlayA(undefined);
	}, [selectedTile, setPlayA]);

	const [gutterColor, setGutterColor] = useState("bg-gray-900")
	useEffect(() => {
		if (activePlayer !== undefined) {
			setGutterColor(bgColorMap[activePlayer.color]);
		} else {
			setGutterColor("bg-green-900");
		}
	}, [activePlayer])

	const [cellSpan, setCellSpan] = useState("");
	const [gutterSpan, setGutterSpan] = useState("");
	useEffect(() => {
		setCellSpan(`${100 / width + 1}%`);
		setGutterSpan(`${50 / width + 1}%`);
	}, [width]);

	function clickSquare(x, y) {
		// useful for choosing your chicken-foot
		if (selectedTile === undefined) {
			setPlayA({ x: x, y: y });
			return;
		}
		if (playA === undefined) {
			setPlayA({ x: x, y: y });
			return;
		}

		if (selectedTile.a == -1 && selectedTile.b == -1) {
			clickForSpacer(x, y);
			return;
		}

		var orientation = undefined;
		if (x === playA.x + 1 && y === playA.y) {
			orientation = "right";
		} else if (x === playA.x - 1 && y === playA.y) {
			orientation = "left";
		} else if (x === playA.x && y === playA.y + 1) {
			orientation = "down";
		} else if (x === playA.x && y === playA.y - 1) {
			orientation = "up";
		} else {
			setPlayA({ x: x, y: y });
			return;
		}
		playTile({
			a: selectedTile.a, b: selectedTile.b,
			coord: {
				x: playA.x,
				y: playA.y,
			},
			orientation: orientation,
			dead: false,
		});
		setPlayA(undefined);
	}

	const [spacerHintPrefix, setSpacerHintPrefix] = useState({});
	useEffect(() => {
		let prefix = {};
		if (selectedTile?.a == -1 && selectedTile?.b == -1 && spacerHints) {
			spacerHints.forEach(hint => {
				const [first, second] = hint.split("-");
			prefix[first] = second;
			prefix[second] = first;
			});
		}
		setSpacerHintPrefix(prefix);
	}, [spacerHints, selectedTile, hints]);

	function clickForSpacer(x, y) {
		if (!(`${x},${y}` in spacerHintPrefix)) {
			setPlayA(undefined);
			return;
		}
		playSpacer({
			a: playA,
			b: {x: x, y: y},
		});
		setPlayA(undefined);
	}

	const [spacerA, setSpacerA] = useState(undefined);
	useEffect(() => {
		setSpacerA(`${spacer?.a.x},${spacer?.a.y}`);
	}, [spacer]);

	return (
		<div onContextMenu={rightClick} className={`aspect-square h-full border-8 border-gray-500 flex items-center justify-center ${gutterColor}`}>
			<div className="aspect-square pb-[100%] min-w-0 min-h-0" style={{ maxHeight: '100%', maxWidth: '100%' }}>
				<div className="aspect-square">
					<table className="w-full h-full table-fixed">
						<tbody>
							{Array.from({ length: height }, (_, y) => (
								<tr key={y}>
									<td className={`p-0 border-0 ${gutterColor}`} style={{ height: cellSpan, width: gutterSpan }}>
									</td>
									{Array.from({ length: width }, (_, x) => (
										<td key={y * width + x} className="p-0 border-0 bg-green-900" style={{ height: cellSpan, width: cellSpan }}>
											<div className="w-full pb-[100%] relative">
												{hints[`${x},${y}`] && (
													<div className="w-full h-full z-20 absolute pointer-events-none">
														<Hint />
													</div>
												)}
												{spacerHintPrefix[`${x},${y}`] && (
													<div className="w-full h-full z-20 absolute pointer-events-none">
														<Hint />
													</div>
												)}
												{spacerA === `${x},${y}` && (
													<div className="w-full h-full z-20 absolute">
														<Spacer spacer={spacer} />
													</div>
												)}
												{tiles[`${x},${y}`] && (
													<div className="w-full h-full z-20 absolute">
														<Tile
															pipsa={tiles[`${x},${y}`].a}
															pipsb={tiles[`${x},${y}`].b}
															orientation={tiles[`${x},${y}`].orientation}
															color={tiles[`${x},${y}`].color}
															dead={tiles[`${x},${y}`].dead}
															lineHeads={lineHeads}
															roundLeader={roundLeader}
															freeLeaders={freeLeaders}
															indicated={indicated}
															setIndicated={setIndicated} />
													</div>
												)}
												{chickenFeet[`${x},${y}`] && (
													<div className="w-full h-full z-30 absolute pointer-events-none">
														<ChickenFoot
															url={chickenFeetURLs[`${x},${y}`]}
															color={chickenFeet[`${x},${y}`]} />
													</div>
												)}
												<div
													className="z-10 absolute inset-0"
													onClick={() => clickSquare(x, y)}
												>
													<Square
														x={x} y={y}
														hoveredSquares={hoveredSquares}
														setMouseIsOver={setMouseIsOver}
														dropCallback={dropCallback}
														center={y == (height - 1) / 2 && (x == (width / 2) - 1 || x == (width / 2))}
														clicked={playA !== undefined && playA.x == x && playA.y == y}
														pips={selectedTile?.a}
													/>
												</div>
											</div>
										</td>
									))}
									<td className={`p-0 border-0 ${gutterColor}`} style={{ height: cellSpan, width: gutterSpan }}>
									</td>
								</tr>
							))}
						</tbody>
					</table>
				</div>
			</div>
		</div>
	);
}

