"use client";

import Pips from "./Pips";

function TileHalf({ pips, orientation }) {
	var outercnm = "w-full aspect-square";
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
			<div className="w-full h-full bg-transparent flex items-center justify-center p-1">
				<Pips pips={pips} />
			</div>
		</div>
	);
}

function Tile({pipsa, pipsb, orientation, back=false, color="white", dead=false}) {
	const colorMap = {
        red: "bg-red-100",
        blue: "bg-blue-100",
        green: "bg-green-100",
        indigo: "bg-indigo-100",
        orange: "bg-orange-100",
        fuchsia: "bg-fuchsia-100",
        white: "bg-white"
    };

	var squareBar = <div className="absolute bottom-[-2px] left-[15%] w-[70%] h-[4px] bg-gray-300"/>;
	var bar = <div className="absolute left-[15%] w-[70%] h-[4px] bg-gray-300"/>;
	
	var height = "";
	var rotate = 'rotate-0'

	if (orientation == "down") {
		height = "h-[200%]";
		rotate = 'rotate-0'
		bar = squareBar;
	}
	if (orientation == "up") {
		height = "h-[200%]";
		rotate = 'rotate-180'
		bar = squareBar;
	}
	if (orientation == "left") {
		height = "h-[200%]";
		rotate = 'rotate-90'
		bar = squareBar;
	}
	if (orientation == "right") {
		height = "h-[200%]";
		rotate = '-rotate-90'
		bar = squareBar;
	}
	return (
		<div className={`h-full w-full ${rotate}`}>
			<div className={height+" w-[100%] p-1"}>
				<div className={`w-full h-full ${colorMap[color]} ${dead ? "border-gray-300" : "border-black"} rounded-lg border-4`}>
					<table className="w-full h-full table-fixed">
						<tbody>
							<tr><td>
								<TileHalf pips={back ? 0 : pipsa} orientation={"down"}/>
								{!back && bar}
							</td></tr>
							<tr><td>
								<TileHalf pips={back ? 0 : pipsb} orientation={"up"}/>
							</td></tr>
						</tbody>
					</table>
				</div>
			</div>
		</div>
	);

}

export default Tile;
