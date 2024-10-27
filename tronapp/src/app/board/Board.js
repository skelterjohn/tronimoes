import Square from "./Square";
import Tile from "./Tile";

function Board({ width = 10, height = 10 }) {
	return (
			<table className="h-[80%] aspect-square table-fixed">
				<tbody>
					{Array.from({length: height}, (_, y) => (
						<tr key={y}>
							{Array.from({length: width}, (_, x) => (
								<td key={y*width+x} className="p-0 border-0 ">
									<div className="w-full z-10 pb-[100%] relative">
										{ x == 3 && y == 3 && (
											<div className="w-full h-full z-20 pb-[100%] absolute">
												<Tile pips={6} orientation={0} />
											</div>
										)}
										{ x == 3 && y == 4 && (
											<div className="w-full h-full z-20 pb-[100%] absolute">
												<Tile pips={5} orientation={1} />
											</div>
										)}
										{ x == 5 && y == 6 && (
											<div className="w-full h-full z-20 pb-[100%] absolute">
												<Tile pips={6} orientation={2} />
											</div>
										)}
										{ x == 6 && y == 6 && (
											<div className="w-full h-full z-20 pb-[100%] absolute">
												<Tile pips={5} orientation={3} />
											</div>
										)}
										<div className="absolute inset-0">
											<Square x={x} y={y} />
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