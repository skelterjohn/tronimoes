import Pips from "./Pips";

function TileHalf({ pips, orientation }) {
	var outercnm = "w-full aspect-square";
	var innercnm = "w-full h-full bg-white flex items-center justify-center ";
		switch (orientation) {
			case "down": // down
				outercnm = `${outercnm} rotate-0`;
				break;
			case "up": // down
				outercnm = `${outercnm} rotate-180`;
				break;
	}
	return (
		<div className={outercnm}>
			<div className={innercnm}>
				<Pips pips={pips} />
			</div>
		</div>
	);
}

function Tile({pipsa, pipsb, orientation}) {
	var rotate = 'rotate-0'
	if (orientation == "up") {
		rotate = 'rotate-180'
	}
	if (orientation == "left") {
		rotate = 'rotate-90'
	}
	if (orientation == "right") {
		rotate = '-rotate-90'
	}
		return (
			<div className={`h-full w-full ${rotate}`}>
				<div className="h-[200%] w-[100%] p-1">
					<div className="w-full h-full bg-white border-black rounded-lg border-2">
						<table className="w-full h-full table-fixed">
							<tbody>
								<tr><td>
									<TileHalf pips={pipsa} orientation={"down"}/>
								</td></tr>
								<tr><td>
									<TileHalf pips={pipsb} orientation={"up"}/>
								</td></tr>
							</tbody>
						</table>
					</div>
				</div>
			</div>
		);

}

export default Tile;
