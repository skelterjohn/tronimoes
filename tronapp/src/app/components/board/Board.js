"use client";

import { useState, useEffect } from 'react';

import Square from "./Square";
import Tile from "./Tile";
import ChickenFoot from './ChickenFoot';

export default function Board({ width = 10, height = 11, tiles, lineHeads, selectedTile, playTile, chickenFeet }) {
	const [playA, setPlayA] = useState(undefined);

	function rightClick(evt) {
		evt.preventDefault();
		setPlayA(undefined);
	}

	useEffect(() => {
		setPlayA(undefined);
	}, [selectedTile]);

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
		<div onContextMenu={rightClick} className="max-w-[75vw] max-h-[75vh] w-fit">
			<table className="w-full table-fixed">
				<tbody>
					{Array.from({length: height}, (_, y) => (
						<tr key={y}>
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
													lineHeads={lineHeads} />
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
						</tr>
					))}
				</tbody>
			</table>
		</div>
	);
}

