import Pips from "./Pips";

function Tile({ pips, orientation }) {
	var outercnm = "pt-1 pr-1 pl-1 w-full aspect-square";
	var innercnm = "w-full h-full bg-white border-black flex items-center justify-center rounded-t-lg  border-t-2 border-l-2 border-r-2 ";
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
				<div className="p-1">	
					<Pips pips={pips} />
				</div>
			</div>
		</div>
	);
}

export default Tile;
