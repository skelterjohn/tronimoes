"use client";

import { useEffect, useState, useRef } from "react";
import Pips from "./Pips";
import Tip, { makeTipBundle } from "@/app/components/tutorial/Tip";

const TileHalf = ({ pips, back, orientation }) => {
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
			<div className="w-full h-full bg-transparent flex items-center justify-center">
				{!back && <Pips pips={pips} />}
			</div>
		</div>
	);
}

export default function Tile({ pipsa, pipsb, orientation, back = false, color = "white", dead = false, selected = false, lineHeads, indicated, setIndicated, hintedTiles }) {



	var squareBar = <div className="absolute bottom-[-2px] left-[15%] w-[70%] h-[4px] bg-gray-300" />;
	var bar = <div className="absolute left-[15%] w-[70%] h-[4px] bg-gray-300" />;

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

	const [isLineHead, setIsLineHead] = useState(false);
	useEffect(() => {
		setIsLineHead(false);
		lineHeads?.forEach((lt) => {
			if (lt.tile.pips_a !== pipsa || lt.tile.pips_b !== pipsb) {
				return;
			}
			setIsLineHead(true);
		});
	}, [pipsa, pipsb, lineHeads])

	const [bgcolor, setBgcolor] = useState("bg-gray-500");
	const [bordercolor, setBordercolor] = useState("border-black");

	const [hinted, setHinted] = useState(false);
	useEffect(() => {
		let h = false;
		hintedTiles && hintedTiles.forEach((ht) => {
			if (ht.a == pipsa && ht.b == pipsb) {
				h = true;
			}
		});
		setHinted(h);
	}, [hintedTiles, pipsa, pipsb]);

	useEffect(() => {
		const colorMap = {
			red: "bg-red-500",
			blue: "bg-blue-500",
			green: "bg-green-500",
			yellow: "bg-yellow-400",
			orange: "bg-orange-500",
			fuchsia: "bg-fuchsia-500",
			white: "bg-white"
		};
		const borderColorMap = {
			red: "border-red-500",
			blue: "border-blue-500",
			green: "border-green-500",
			yellow: "border-yellow-400",
			orange: "border-orange-500",
			fuchsia: "border-fuchsia-500",
			white: "border-white"
		};
		const lineHeadBorderColorMap = {
			red: "border-red-700",
			blue: "border-blue-700",
			green: "border-green-700",
			yellow: "border-yellow-600",
			orange: "border-orange-700",
			fuchsia: "border-fuchsia-700",
			white: "border-white"
		};
		const selectedColorMap = {
			red: "bg-red-700",
			blue: "bg-blue-700",
			green: "bg-green-700",
			yellow: "bg-yellow-600",
			orange: "bg-orange-700",
			fuchsia: "bg-fuchsia-700",
			white: "bg-gray-200"
		};
		setBgcolor(dead ? "bg-gray-500" : (selected ? selectedColorMap[color] : colorMap[color]));
		if (dead) {
			setBordercolor(borderColorMap[color]);
		} else if (indicated?.a == pipsa && indicated?.b == pipsb) {
			setBordercolor("border-white")
		} else if (isLineHead) {
			setBordercolor(lineHeadBorderColorMap[color])
		} else {
			setBordercolor("border-black")
		}
	}, [selected, dead, isLineHead, indicated, color, pipsa, pipsb]);

	function tileClicked() {
		if (!isLineHead || dead) {
			return;
		}
		setIndicated({ a: pipsa, b: pipsb });
	}

	const playATileBundle = makeTipBundle("This tile is raised, which means it is currently playable.");
	useEffect(() => {
		if (hinted) {
			playATileBundle.setShow(true);
		}
	}, [hinted]);

	const rotateBundle = makeTipBundle("You can rotate the tile by clicking or tapping on it.");
	useEffect(() => {
		if (selected) {
			rotateBundle.setShow(true);
		}
	}, [selected]);

	const dragBundle = makeTipBundle("Once it's oriented, you can drag it to the board.");
	useEffect(() => {
		if (rotateBundle.done && selected && orientation !== "down") {
			dragBundle.setShow(true);
		}
	}, [selected, orientation, rotateBundle.done]);

	return (
		<div className={`h-full w-full ${rotate} ${hinted && "-translate-y-2"}`}>
			<Tip bundle={playATileBundle} />
			<Tip bundle={dragBundle} />
			<Tip bundle={rotateBundle} />
			<div className={height + " w-[100%]"}>
				<div className={`w-full h-full ${bgcolor} ${bordercolor} rounded-lg border-2`} onClick={() => tileClicked()}>
					<table className="w-full h-full table-fixed">
						<tbody>
							<tr><td>
								<TileHalf pips={back ? 0 : pipsa} back={back} orientation={"down"} />
								{!back && bar}
							</td></tr>
							<tr><td>
								<TileHalf pips={back ? 0 : pipsb} back={back} orientation={"up"} />
							</td></tr>
						</tbody>
					</table>
				</div>
			</div>
		</div>
	);

}
