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
												<Tile pipsa={0} pipsb={1} orientation={"down"} />
											</div>
										)}
										{ x == 0 && y == 2 && (
											<div className="w-full h-full z-20 absolute">
												<Tile pipsa={2} pipsb={3} orientation={"right"} />
											</div>
										)}
										{ x == 1 && y == 3 && (
											<div className="w-full h-full z-20 absolute">
												<Tile pipsa={4} pipsb={5} orientation={"down"} />
											</div>
										)}
										{ x == 1 && y == 5 && (
											<div className="w-full h-full z-20 absolute">
												<Tile pipsa={6} pipsb={7} orientation={"left"} />
											</div>
										)}

										{ x == 4 && y == 0 && (
											<div className="w-full h-full z-20 absolute">
												<Tile pipsa={8} pipsb={9} orientation={"down"} />
											</div>
										)}
										{ x == 4 && y == 2 && (
											<div className="w-full h-full z-20 absolute">
												<Tile pipsa={10} pipsb={11} orientation={"right"} />
											</div>
										)}
										{ x == 5 && y == 3 && (
											<div className="w-full h-full z-20 absolute">
												<Tile pipsa={12} pipsb={13} orientation={"down"} />
											</div>
										)}
										{ x == 5 && y == 5 && (
											<div className="w-full h-full z-20 absolute">
												<Tile pipsa={14} pipsb={15} orientation={"left"} />
											</div>
										)}
										{ x == 7 && y == 7 && (
											<div className="w-full h-full z-20 absolute">
												<Tile pipsa={16} pipsb={16} orientation={"left"} />
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