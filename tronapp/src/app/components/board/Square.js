import { useState, useEffect } from "react";
import Pips from "./Pips";

export default function Square({ x, y, center = false, clicked = false, pips, hoveredSquares }) {
	const [hovered, setHovered] = useState(false);
	useEffect(() => {
		setHovered(hoveredSquares.has(`${x},${y}`));
	}, [hoveredSquares]);

	const [bgColor, setBgColor] = useState("");
	useEffect(() => {
		if (center) {
			setBgColor("bg-gray-400");
			return;
		}
		if (hovered) {
			setBgColor("bg-black");
			return;
		}
		if ((x + y) % 2 == 0) {
			setBgColor("bg-blue-200");
		} else {
			setBgColor("bg-slate-200");
		}
	}, [center, x, y,hovered]);

	var cnm = `w-full aspect-square ${bgColor}`;
	if (clicked) {
		cnm = `${cnm} border border-2 border-black`;
	}


	return <div data-tron_x={x} data-tron_y={y} className={cnm}>
		{clicked && <Pips pips={pips} />}
	</div>;
}
