"use client";

import { useState, useEffect } from 'react';

import Square from "./Square";
import Tile from "./Tile";
import ChickenFoot from './ChickenFoot';

const bgColorMap = {
	red: "bg-red-300",
	blue: "bg-blue-300",
	green: "bg-green-300",
	indigo: "bg-indigo-300",
	orange: "bg-orange-300",
	fuchsia: "bg-fuchsia-300",
	white: "bg-white"
};

export default function Board({ width = 10, height = 11, tiles, lineHeads, selectedTile, playTile, chickenFeet, indicated, setIndicated, activePlayer }) {
	const [playA, setPlayA] = useState(undefined);

	function rightClick(evt) {
		evt.preventDefault();
		setPlayA(undefined);
		setIndicated(undefined);
	}

	useEffect(() => {
		setPlayA(undefined);
	}, [selectedTile]);

	const [gutterColor, setGutterColor] = useState("bg-gray")
	useEffect(() => {
		if (activePlayer !== undefined) {
			setGutterColor(bgColorMap[activePlayer.color]);
		} else {
			setGutterColor("bg-gray-300");
		}
	}, [activePlayer])

	function clickSquare(x, y) {
		if (selectedTile===undefined) {
			return;
		}
		if (playA===undefined) {
			setPlayA({x:x, y:y});
		} else {
			var orientation = undefined;
			if (x === playA.x+1 && y === playA.y) {
				orientation = "right";
			} else if (x === playA.x-1 && y === playA.y) {
				orientation = "left";
			} else if (x === playA.x && y === playA.y+1) {
				orientation = "down";
			} else if (x === playA.x && y === playA.y-1) {
				orientation = "up";
			} else {
				setPlayA({x:x, y:y});
				return;
			}
			playTile({
				a:selectedTile.a, b:selectedTile.b,
				x:playA.x, y:playA.y,
				orientation:orientation,
				dead:false,
			});
			setPlayA(undefined);
		}
	}

	return (
		<div onContextMenu={rightClick} className="h-full w-full flex items-center justify-center overflow-hidden">
			<div className="aspect-square min-w-0 min-h-0" style={{ maxHeight: '100%', maxWidth: '100%' }}>
				<table className="w-full h-full table-fixed">
					<tbody>
						{Array.from({length: height}, (_, y) => (
							<tr key={y}>
								<td className={`p-0 ${gutterColor} border-0 w-[4.76%]`}>
								</td>
								{Array.from({length: width}, (_, x) => (
									<td key={y*width+x} className="p-0 border-0">
										<div className="w-full pb-[100%] relative">
											{ tiles[`${x},${y}`] && (
												<div className="w-full h-full z-20 absolute">
													<Tile 
														pipsa={tiles[`${x},${y}`].a}
														pipsb={tiles[`${x},${y}`].b}
														orientation={tiles[`${x},${y}`].orientation}
														color={tiles[`${x},${y}`].color}
														dead={tiles[`${x},${y}`].dead}
														lineHeads={lineHeads}
														indicated={indicated}
														setIndicated={setIndicated} />
												</div>
											)}
											{ chickenFeet[`${x},${y}`] && (
												<div className="w-full h-full z-30 absolute pointer-events-none">
													<ChickenFoot 
														color={chickenFeet[`${x},${y}`]} />
												</div>
											)}
											<div
												className="z-10 absolute inset-0"
												onClick={()=>clickSquare(x, y)}
											>
												<Square
													x={x} y={y}
													center={y == (height-1)/2 && (x == (width/2)-1 || x == (width/2))}
													clicked={playA!==undefined && playA.x==x && playA.y==y}
													pips={selectedTile?.a}
												/>
											</div>
										</div>
									</td>
								))}
								<td className={`p-0 ${gutterColor} border-0 w-[4.76%]`}>
								</td>
							</tr>
						))}
					</tbody>
				</table>
			</div>
		</div>
	);
}

