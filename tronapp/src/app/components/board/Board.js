"use client";

import { useState, useEffect, useRef } from 'react';
import Tip, { useTipBundle } from "@/app/components/tutorial/Tip";

import Spacer from "./Spacer";
import Square from "./Square";
import Tile from "./Tile";
import ChickenFoot from './ChickenFoot';
import Hint from './Hint';

const bgColorMap = {
	red: "bg-red-500",
	blue: "bg-blue-500",
	green: "bg-green-500",
	yellow: "bg-yellow-500",
	orange: "bg-orange-500",
	fuchsia: "bg-fuchsia-500",
	white: "bg-white"
};

export default function Board({ 
		width = 10, height = 11,
		tiles, spacer,
		lineHeads,
		roundLeader, freeLeaders,
		selectedTile,
		playTile, playSpacer,
		chickenFeet, chickenFeetURLs,
		indicated, setIndicated,
		activePlayer, hints,
		playA, setPlayA,
		spacerHints, clearSpacer,
		hoveredSquares, setMouseIsOver,
		dropCallback,
		setSquareSpan
	}) {
	const [zoom, setZoom] = useState(1);
	const [position, setPosition] = useState({ x: 0, y: 0 });
	const boardContainerRef = useRef(null);
	const [touchStartDistance, setTouchStartDistance] = useState(null);
	const [touchStartZoom, setTouchStartZoom] = useState(1);
	const [touchStartPosition, setTouchStartPosition] = useState(null);
	const [lastTouchPosition, setLastTouchPosition] = useState(null);
	const [isDragging, setIsDragging] = useState(false);
	const [dragStart, setDragStart] = useState(null);

	function handleWheel(evt) {
		// evt.preventDefault();
		
		const container = boardContainerRef.current;
		if (!container) return;

		const rect = container.getBoundingClientRect();
		const mouseX = evt.clientX - rect.left;
		const mouseY = evt.clientY - rect.top;

		const mouseXPercent = 100 * mouseX / rect.width;
		const mouseYPercent = 100 * mouseY / rect.height;

		const boardX = mouseXPercent - position.x;

		const delta = evt.deltaY * -0.001;
		const newZoom = Math.min(Math.max(zoom * (1 + delta), 1), 3);
		
		// Calculate new position to keep mouse point fixed
		const newPosition = {
			x: mouseXPercent - (mouseXPercent - position.x) * (newZoom / zoom),
			y: mouseYPercent - (mouseYPercent - position.y) * (newZoom / zoom)
		};


		setZoom(newZoom);
		setPosition(newZoom === 1 ? {x: 0, y: 0} : newPosition);
	}

	function rightClick(evt) {
		evt.preventDefault();
		if (!isDragging) {
			setPlayA(undefined);
			setIndicated(undefined);
			clearSpacer();
		}
	}

	function handleMouseDown(evt) {
		if (evt.button === 2 && zoom > 1) { // Right click
			evt.preventDefault();
			const container = boardContainerRef.current;
			if (!container) return;

			const rect = container.getBoundingClientRect();
			setIsDragging(true);
			setDragStart({
				x: evt.clientX - rect.left,
				y: evt.clientY - rect.top
			});
		}
	}

	function handleMouseMove(evt) {
		if (isDragging) {
			evt.preventDefault();
			const container = boardContainerRef.current;
			if (!container) return;

			const rect = container.getBoundingClientRect();
			const currentPos = {
				x: evt.clientX - rect.left,
				y: evt.clientY - rect.top
			};

			// Calculate the movement as a percentage of container size
			const deltaX = (currentPos.x - dragStart.x) / rect.width * 100;
			const deltaY = (currentPos.y - dragStart.y) / rect.height * 100;

			setPosition({
				x: position.x + deltaX,
				y: position.y + deltaY
			});

			setDragStart(currentPos);
		}
	}

	function handleMouseUp() {
		setIsDragging(false);
		setDragStart(null);
	}

	useEffect(() => {
		setPlayA(undefined);
	}, [selectedTile, setPlayA]);

	const [gutterColor, setGutterColor] = useState("bg-gray-900")
	useEffect(() => {
		if (activePlayer !== undefined) {
			setGutterColor(bgColorMap[activePlayer.color]);
		} else {
			setGutterColor("bg-green-900");
		}
	}, [activePlayer])

	const [cellSpan, setCellSpan] = useState("");
	const [gutterSpan, setGutterSpan] = useState("");
	useEffect(() => {
		setCellSpan(`${100 / width + 1}%`);
		setGutterSpan(`${50 / width + 1}%`);
	}, [width]);

	function clickSquare(x, y) {
		// useful for choosing your chicken-foot
		if (selectedTile === undefined) {
			setPlayA({ x: x, y: y });
			return;
		}
		if (playA === undefined) {
			setPlayA({ x: x, y: y });
			return;
		}

		if (selectedTile.a == -1 && selectedTile.b == -1) {
			clickForSpacer(x, y);
			return;
		}

		var orientation = undefined;
		if (x === playA.x + 1 && y === playA.y) {
			orientation = "right";
		} else if (x === playA.x - 1 && y === playA.y) {
			orientation = "left";
		} else if (x === playA.x && y === playA.y + 1) {
			orientation = "down";
		} else if (x === playA.x && y === playA.y - 1) {
			orientation = "up";
		} else {
			setPlayA({ x: x, y: y });
			return;
		}
		playTile({
			a: selectedTile.a, b: selectedTile.b,
			coord: {
				x: playA.x,
				y: playA.y,
			},
			orientation: orientation,
			dead: false,
		});
		setPlayA(undefined);
	}

	const [spacerHintPrefix, setSpacerHintPrefix] = useState({});
	useEffect(() => {
		let prefix = {};
		if (selectedTile?.a == -1 && selectedTile?.b == -1 && spacerHints) {
			spacerHints.forEach(hint => {
				const [first, second] = hint.split("-");
			prefix[first] = second;
			prefix[second] = first;
			});
		}
		setSpacerHintPrefix(prefix);
	}, [spacerHints, selectedTile, hints]);

	function clickForSpacer(x, y) {
		if (!(`${x},${y}` in spacerHintPrefix)) {
			setPlayA(undefined);
			return;
		}
		playSpacer({
			a: playA,
			b: {x: x, y: y},
		});
		setPlayA(undefined);
	}

	const [spacerA, setSpacerA] = useState(undefined);
	useEffect(() => {
		setSpacerA(`${spacer?.a.x},${spacer?.a.y}`);
	}, [spacer]);

	const chickenFootBundle = useTipBundle("This is a chicken foot. The player whose foot this is can only play on their own line until the foot is gone. Other players who are not footed can also play on this line.");
	useEffect(() => {
		if (Object.keys(chickenFeet).length > 0) {
		chickenFootBundle.setShow(true);
		}
	}, [chickenFeet]);

	const lineBundle = useTipBundle("This is a line. Build your line, matching pips, to wall-in your opponents can come out on top.");
	useEffect(() => {
		if (Object.keys(tiles).length > 1) {
			lineBundle.setShow(true);
		}
	}, [tiles]);

	const playableBoardRef = useRef(null);
	useEffect(() => {
		const updateSpan = () => {
			if (playableBoardRef.current) {
				setSquareSpan(playableBoardRef.current.clientHeight / height);
			}
		};
	
		// Initial measurement
		updateSpan();

		const spanObserver = new ResizeObserver(updateSpan);
		if (playableBoardRef.current) {
			spanObserver.observe(playableBoardRef.current);
		}
	
		// Cleanup
		return () => {
			spanObserver.disconnect();
		}
	}, [playableBoardRef, setSquareSpan, height]);

	function getDistance(touch1, touch2) {
		return Math.hypot(
			touch1.clientX - touch2.clientX,
			touch1.clientY - touch2.clientY
		);
	}

	function handleTouchStart(evt) {
		if (evt.touches.length === 2) {
			evt.preventDefault();
			const distance = getDistance(evt.touches[0], evt.touches[1]);
			setTouchStartDistance(distance);
			setTouchStartZoom(zoom);
		} else if (evt.touches.length === 1 && zoom > 1) {
			// Only prevent default for panning when zoomed in
			evt.preventDefault();
			const touch = evt.touches[0];
			const container = boardContainerRef.current;
			if (!container) return;

			const rect = container.getBoundingClientRect();
			setTouchStartPosition({
				x: touch.clientX - rect.left,
				y: touch.clientY - rect.top
			});
			setLastTouchPosition({
				x: touch.clientX - rect.left,
				y: touch.clientY - rect.top
			});
		}
		// Single taps when not zoomed will propagate normally
	}

	function handleTouchMove(evt) {
		if (evt.touches.length === 2) {
			evt.preventDefault();
			const container = boardContainerRef.current;
			if (!container || !touchStartDistance) return;

			const rect = container.getBoundingClientRect();
			const touch1 = evt.touches[0];
			const touch2 = evt.touches[1];
			
			// Calculate center point of the two touches
			const centerX = (touch1.clientX + touch2.clientX) / 2 - rect.left;
			const centerY = (touch1.clientY + touch2.clientY) / 2 - rect.top;
			
			const centerXPercent = 100 * centerX / rect.width;
			const centerYPercent = 100 * centerY / rect.height;

			const currentDistance = getDistance(touch1, touch2);
			const scale = currentDistance / touchStartDistance;
			const newZoom = Math.min(Math.max(touchStartZoom * scale, 1), 3);

			// Calculate new position to keep pinch center point fixed
			const newPosition = {
				x: centerXPercent - (centerXPercent - position.x) * (newZoom / zoom),
				y: centerYPercent - (centerYPercent - position.y) * (newZoom / zoom)
			};

			setZoom(newZoom);
			setPosition(newZoom === 1 ? {x: 0, y: 0} : newPosition);
		} else if (evt.touches.length === 1 && zoom > 1 && touchStartPosition) {
			evt.preventDefault();
			const touch = evt.touches[0];
			const container = boardContainerRef.current;
			if (!container) return;

			const rect = container.getBoundingClientRect();
			const currentTouch = {
				x: touch.clientX - rect.left,
				y: touch.clientY - rect.top
			};

			// Calculate the movement as a percentage of container size
			const deltaX = (currentTouch.x - lastTouchPosition.x) / rect.width * 100;
			const deltaY = (currentTouch.y - lastTouchPosition.y) / rect.height * 100;

			setPosition({
				x: position.x + deltaX,
				y: position.y + deltaY
			});

			setLastTouchPosition(currentTouch);
		}
	}

	function handleTouchEnd(evt) {
		if (evt.touches.length !== 0) {
			evt.preventDefault(); // Only prevent default for multi-touch gestures
		}
		setTouchStartDistance(null);
		setTouchStartPosition(null);
		setLastTouchPosition(null);
	}

	return (
		<div 
			ref={boardContainerRef}
			onWheel={handleWheel}
			onContextMenu={rightClick}
			onMouseDown={handleMouseDown}
			onMouseMove={handleMouseMove}
			onMouseUp={handleMouseUp}
			onMouseLeave={handleMouseUp}
			onTouchStart={handleTouchStart}
			onTouchMove={handleTouchMove}
			onTouchEnd={handleTouchEnd}
			style={{
				touchAction: 'none',
				cursor: isDragging ? 'grabbing' : (zoom > 1 ? 'grab' : 'default')
			}}
			className={`aspect-square w-full h-full border-8 border-gray-500 flex items-center justify-center overflow-hidden ${gutterColor}`}
			>
			<div 
				className="aspect-square"
				style={{ 
					transform: `translate(${position.x}%, ${position.y}%) scale(${zoom})`,
					transformOrigin: '0 0',
					transition: 'none',
					width: '100%',
					height: '100%'
				}}
			>
				<div className="aspect-square w-full h-full">
					<table className="w-full h-full table-fixed" ref={playableBoardRef}>
						<tbody>
							{Array.from({ length: height }, (_, y) => (
								<tr key={y}>
									<td className={`p-0 border-0 ${gutterColor}`} style={{ height: cellSpan, width: gutterSpan }}>
									</td>
									{Array.from({ length: width }, (_, x) => (
										<td key={y * width + x} className="p-0 border-0 bg-green-900" style={{ height: cellSpan, width: cellSpan }}>
											<div className="w-full pb-[100%] relative">
												{hints[`${x},${y}`] && (
													<div className="w-full h-full z-20 absolute pointer-events-none">
														<Hint />
													</div>
												)}
												{spacerHintPrefix[`${x},${y}`] && (
													<div className="w-full h-full z-20 absolute pointer-events-none">
														<Hint />
													</div>
												)}
												{spacerA === `${x},${y}` && (
													<div className="w-full h-full z-20 absolute">
														<Spacer spacer={spacer} />
													</div>
												)}
												{tiles[`${x},${y}`] && (
													<div className="w-full h-full z-20 absolute">
														{tiles[`${x},${y}`].color !== "white" && <Tip bundle={lineBundle} /> }
														<Tile
															pipsa={tiles[`${x},${y}`].a}
															pipsb={tiles[`${x},${y}`].b}
															orientation={tiles[`${x},${y}`].orientation}
															color={tiles[`${x},${y}`].color}
															dead={tiles[`${x},${y}`].dead}
															lineHeads={lineHeads}
															roundLeader={roundLeader}
															freeLeaders={freeLeaders}
															indicated={indicated}
															setIndicated={setIndicated} />
													</div>
												)}
												{chickenFeet[`${x},${y}`] && (
													<div className="w-full h-full z-30 absolute pointer-events-none">
														<Tip bundle={chickenFootBundle} />
														<ChickenFoot
															url={chickenFeetURLs[`${x},${y}`]}
															color={chickenFeet[`${x},${y}`]} />
													</div>
												)}
												<div
													className="z-10 absolute inset-0"
													onClick={() => clickSquare(x, y)}
												>
													<Square
														x={x} y={y}
														hoveredSquares={hoveredSquares}
														setMouseIsOver={setMouseIsOver}
														dropCallback={dropCallback}
														center={y == (height - 1) / 2 && (x == (width / 2) - 1 || x == (width / 2))}
														clicked={playA !== undefined && playA.x == x && playA.y == y}
														pips={selectedTile?.a}
													/>
												</div>
											</div>
										</td>
									))}
									<td className={`p-0 border-0 ${gutterColor}`} style={{ height: cellSpan, width: gutterSpan }}>
									</td>
								</tr>
							))}
						</tbody>
					</table>
				</div>
			</div>
		</div>
	);
}

