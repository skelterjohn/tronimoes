import { useState, useEffect } from "react";
import Pips from "./Pips";

export default function Square({
		x, y,
		center = false, clicked = false,
		pips,
		hoveredSquares, setMouseIsOver, dropCallback,
		interactive = true }) {
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
			setBgColor("bg-[#18222b]");
		} else {
			setBgColor("bg-[#34495E]");
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
		onDragOver={interactive ? onDragOver : undefined}
		onDragEnter={interactive ? onDragEnter : undefined}
		onDrop={interactive ? onDrop : undefined}
		data-tron_x={x}
		data-tron_y={y}
		className={`relative w-full aspect-square ${bgColor} ${clicked && "border border-2 border-black"}`}
		style={interactive ? undefined : { pointerEvents: 'none' }}
	>
		{x === 0 && (
			<span className="absolute top-0.5 left-0.5 font-game text-[8px] text-gray-400 select-none">{y}</span>
		)}
		{y === 0 && x !== 0 && (
			<span className="absolute top-0.5 left-0.5 font-game text-[8px] text-gray-400 select-none">{x}</span>
		)}
		{clicked && <Pips pips={pips} />}
	</div>;
}
