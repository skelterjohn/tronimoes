import Square from "./Square";

function Board({ width = 10, height = 10 }) {
	return (
			<table className="h-[80%] aspect-square table-fixed">
				<tbody>
					{Array.from({length: height}, (_, y) => (
						<tr key={y}>
							{Array.from({length: width}, (_, x) => (
								<td key={y*width+x} className="p-0 border-0 ">
									<div className="w-full pb-[100%] relative">
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