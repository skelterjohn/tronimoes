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
											<div className="w-full h-full z-20 pt-1 pl-1 pr-1 absolute">
												<Tile pipsa={4} pipsb={5} orientation={0} />
											</div>
										)}
										{ x == 1 && y == 1 && (
											<div className="w-full h-full z-20 pt-1 pl-1 pr-1 absolute">
												<Tile pipsa={3} pipsb={6} orientation={1} />
											</div>
										)}
										{ x == 1 && y == 3 && (
											<div className="w-full h-full z-20 pt-1 pl-1 pr-1 absolute">
												<Tile pipsa={2} pipsb={1} orientation={2} />
											</div>
										)}
										{ x == 3 && y == 5 && (
											<div className="w-full h-full z-20 pt-1 pl-1 pr-1 absolute">
												<Tile pipsa={2} pipsb={1} orientation={3} />
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