"use client";

import Square from "./Square";
import Tile from "./Tile";

function Board({ width = 10, height = 11, tiles }) {
	return (
			<table className="w-full aspect-square table-fixed">
				<tbody>
					{Array.from({length: height}, (_, y) => (
						<tr key={y}>
							{Array.from({length: width}, (_, x) => (
								<td key={y*width+x} className="p-0 border-0 ">
									<div className="w-full pb-[100%] relative">
										{ tiles[`${x},${y}`] && (
											<div className="w-full h-full z-20 absolute">
												<Tile 
													pipsa={tiles[`${x},${y}`].a}
													pipsb={tiles[`${x},${y}`].b}
													orientation={tiles[`${x},${y}`].orientation}
													color={tiles[`${x},${y}`].color}
													dead={tiles[`${x},${y}`].dead} />
											</div>
										)}
										<div className="z-10 absolute inset-0">
											<Square
												x={x} y={y}
												center={y == (height-1)/2 && (x == (width/2)-1 || x == (width/2))}/>
										</div>
									</div>
								</td>
							))}
						</tr>
					))}
				</tbody>
			</table>
	)
}

export default Board;