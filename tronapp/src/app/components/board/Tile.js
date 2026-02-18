"use client";

import { useEffect, useState } from "react";
import Pips from "./Pips";

const TileHalf = ({ pips, back, orientation, tileRotation }) => {
	var outercnm = "w-full aspect-square";
	var halfRotation = 0;
	switch (orientation) {
		case "down": // down
			outercnm = `${outercnm} rotate-0`;
			halfRotation = 0;
			break;
		case "up": // down
			outercnm = `${outercnm} rotate-180`;
			halfRotation = 180;
			break;
	}
	return (
		<div className={outercnm}>
			<div className="w-full h-full bg-transparent flex items-center justify-center">
				{!back && <Pips pips={pips} parentRotation={tileRotation + halfRotation} />}
			</div>
		</div>
	);
}

export default function Tile({pipsa, pipsb, orientation, back = false, color = "white", dead = false, selected = false, lineHeads, indicated, setIndicated, hintedTiles, roundLeader = undefined, freeLeaders = undefined, interactive = true, last = false }) {
	var bar = <div className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-[70%] h-[4px] bg-gray-300" />;

	var height = "";
	var rotate = 'rotate-0'
	var tileRotation = 0;

	if (orientation == "down") {
		height = "h-[200%]";
		rotate = 'rotate-0'
		tileRotation = 0;
	}
	if (orientation == "up") {
		height = "h-[200%]";
		rotate = 'rotate-180'
		tileRotation = 180;
	}
	if (orientation == "left") {
		height = "h-[200%]";
		rotate = 'rotate-90'
		tileRotation = 90;
	}
	if (orientation == "right") {
		height = "h-[200%]";
		rotate = '-rotate-90'
		tileRotation = -90;
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
		const selectedColorMap = {
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
		const colorMap = {
			red: "bg-red-700",
			blue: "bg-blue-700",
			green: "bg-green-700",
			yellow: "bg-yellow-600",
			orange: "bg-orange-700",
			fuchsia: "bg-fuchsia-700",
			white: "bg-gray-200"
		};
		const lineHeadColorMap = selectedColorMap;
		
		if (dead) {
			setBgcolor("bg-gray-500");
			setBordercolor(borderColorMap[color]);
		} else if (indicated?.a == pipsa && indicated?.b == pipsb) {
			// this is only for within the hand.
			setBgcolor(selected ? selectedColorMap[color] : colorMap[color]);
			setBordercolor("border-white")
		} else if (isLineHead) {
			setBgcolor(lineHeadColorMap[color]);
			setBordercolor("border-black")
		} else {
			setBgcolor(selected ? selectedColorMap[color] : colorMap[color]);
			setBordercolor("border-black")
		}
	}, [selected, dead, isLineHead, indicated, color, pipsa, pipsb]);

	function tileClicked() {
		if (!isLineHead || dead) {
			return;
		}
		setIndicated({ a: pipsa, b: pipsb });
	}

	return (
		<div className={`h-full w-full ${rotate} ${hinted && "-translate-y-2"}`}>
			<div className={height + " w-full"}>
				<div
					className={`relative overflow-hidden w-full h-full ${bgcolor} ${bordercolor} rounded-lg border-2 ${last ? '[box-shadow:0_0_0_2px_rgba(251,191,36,0.9),0_0_12px_4px_rgba(245,158,11,0.6)]' : ''}`}
					onClick={interactive ? () => tileClicked() : undefined}
					style={interactive ? undefined : { pointerEvents: 'none' }}
				>
					<table className="w-full h-full table-fixed">
						<tbody>
							<tr><td>
								<TileHalf pips={back ? 0 : pipsa} back={back} orientation={"down"} tileRotation={tileRotation} />
								{!back && bar}
							</td></tr>
							<tr><td>
								<TileHalf pips={back ? 0 : pipsb} back={back} orientation={"up"} tileRotation={tileRotation} />
							</td></tr>
						</tbody>
					</table>
				</div>
			</div>
		</div>
	);

}
