import { useState, useEffect } from "react";
import Pips from "./Pips";

export default function Square({
		x, y,
		center = false, clicked = false,
		pips,
		hoveredSquares, setMouseIsOver, dropCallback }) {
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
			setBgColor("bg-white");
			return;
		}
		if ((x + y) % 2 == 0) {
			setBgColor("bg-gray-900");
		} else {
			setBgColor("bg-green-900");
		}
	}, [center, x, y,hovered]);

	function onDragOver(e) {
		e.preventDefault();
		setMouseIsOver([x, y]);
	}

	function onDragEnter(e) {
		e.preventDefault();
	}

	function onDrop(e) {
		e.preventDefault();
		if (dropCallback === undefined) {
			return;
		}
		dropCallback(x, y);
	}

	return <div
		onDragOver={onDragOver}
		onDragEnter={onDragEnter}
		onDrop={onDrop}
		data-tron_x={x}
		data-tron_y={y}
		className={`w-full aspect-square ${bgColor} ${clicked && "border border-2 border-black"}`}
	>
		{clicked && <Pips pips={pips} />}
	</div>;
}
