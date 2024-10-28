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
									<div className="w-full pb-[100%] relative">
										{ x == 0 && y == 0 && (
											<div className="w-full h-full z-20 absolute">
												<Tile pipsa={4} pipsb={5} orientation={"down"} />
											</div>
										)}
										{ x == 1 && y == 1 && (
											<div className="w-full h-full z-20 absolute">
												<Tile pipsa={3} pipsb={6} orientation={"up"} />
											</div>
										)}
										{ x == 1 && y == 3 && (
											<div className="w-full h-full z-20 absolute">
												<Tile pipsa={2} pipsb={1} orientation={"left"} />
											</div>
										)}
										{ x == 3 && y == 5 && (
											<div className="w-full h-full z-20 absolute">
												<Tile pipsa={2} pipsb={1} orientation={"right"} />
											</div>
										)}
										<div className="z-10 absolute inset-0">
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