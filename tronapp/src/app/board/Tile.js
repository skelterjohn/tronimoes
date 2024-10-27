function Tile({ pips, orientation }) {
	var outercnm = "w-full aspect-square";
	var innercnm = "w-full h-full border-black";
		switch (orientation) {
			case 0: // down
				outercnm = `${outercnm} pt-1 pr-1 pl-1 `;
				innercnm = `${innercnm} rounded-t-lg  border-t-2 border-l-2 border-r-2 `;
				break;
			case 1: // up
				outercnm = `${outercnm} pb-1 pr-1 pl-1 `;
				innercnm = `${innercnm} rounded-b-lg  border-b-2 border-l-2 border-r-2 `;
				break;
			case 2: // right
				outercnm = `${outercnm} pt-1 pb-1 pl-1 `;
				innercnm = `${innercnm} rounded-l-lg  border-t-2 border-b-2 border-l-2 `;
				break;
			case 3: // left
				outercnm = `${outercnm} pt-1 pb-1 pr-1 `;
				innercnm = `${innercnm} rounded-r-lg  border-t-2 border-b-2 border-r-2 `;
				break;
	}
	return (
		<div className={outercnm}>
			<div className={innercnm}>
				Tile
			</div>
		</div>
	);
}

export default Tile;
