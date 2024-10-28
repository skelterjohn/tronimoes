import Pips from "./Pips";

function TileHalf({ pips, orientation }) {
	var outercnm = "w-full aspect-square";
	var innercnm = "w-full h-full bg-white flex items-center justify-center ";
		switch (orientation) {
			case 0: // down
				outercnm = `${outercnm} rotate-0`;
				break;
			case 1: // down
				outercnm = `${outercnm} rotate-180`;
				break;
			case 2: // right
				outercnm = `${outercnm} -rotate-90`;
				break;
			case 3: // left
				outercnm = `${outercnm} rotate-90`;
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
	if (orientation == 0) {
		return (
			<div className="h-[200%] w-full p-1 bg-white border-black rounded-lg border-2 ">
				<table className="w-full h-full table-fixed ">
					<tbody className="h-full w-full">
						<tr className="w-full"><td>
							<TileHalf pips={pipsa} orientation={0}/>
						</td></tr>
						<tr className="w-full"><td>
							<TileHalf pips={pipsb} orientation={1}/>
						</td></tr>
					</tbody>
				</table>
			</div>
		);
	}
	if (orientation == 1) {
		return (
			<div className="h-[200%] w-full p-1 bg-white border-black rounded-lg border-2 ">
				<table className="w-full h-full table-fixed ">
					<tbody className="h-full w-full">
						<tr className="w-full"><td>
							<TileHalf pips={pipsb} orientation={0}/>
						</td></tr>
						<tr className="w-full"><td>
							<TileHalf pips={pipsa} orientation={1}/>
						</td></tr>
					</tbody>
				</table>
			</div>
		);
	}
	if (orientation == 2) {
		return (
			<div className="w-[200%] h-full p-1 bg-white border-black rounded-lg border-2 ">
				<table className="w-full h-full table-fixed ">
					<tbody className="h-full w-full">
						<tr className="w-full"><td>
							<TileHalf pips={pipsa} orientation={2}/>
						</td><td>
							<TileHalf pips={pipsb} orientation={3}/>
						</td></tr>
					</tbody>
				</table>
			</div>
		);
	}
	if (orientation == 3) {
		return (
			<div className="w-[200%] h-full p-1 bg-white border-black rounded-lg border-2 ">
				<table className="w-full h-full table-fixed ">
					<tbody className="h-full w-full">
						<tr className="w-full"><td>
							<TileHalf pips={pipsb} orientation={2}/>
						</td><td>
							<TileHalf pips={pipsa} orientation={3}/>
						</td></tr>
					</tbody>
				</table>
			</div>
		);
	}
}

export default Tile;
