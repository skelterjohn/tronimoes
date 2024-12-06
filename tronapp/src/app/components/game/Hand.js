import { useState, useEffect, useCallback } from "react";
import Tile from '../board/Tile';
import ChickenFoot from '../board/ChickenFoot';
import { Button } from "antd";
import Image from "next/image";
import Reaction from "./Reaction";
import Tip, { useTipBundle } from "@/app/components/tutorial/Tip";

function Hand({
		player, players,
		hidden = false, dead = false,
		selectedTile, setSelectedTile,
		playerTurn,
		drawTile, passTurn,
		roundInProgress,
		hintedTiles, hintedSpacer,
		bagCount, turnIndex, playTile,
		setHoveredSquares, mouseIsOver,
		dragOrientation, setDragOrientation, toggleOrientation,
		setShowReactModal
	}) {
	const [handOrder, setHandOrder] = useState([]);
	const [touchStartPos, setTouchStartPos] = useState(null);
	const [draggedTile, setDraggedTile] = useState(null);
	const [spacerAvailable, setSpacerAvailable] = useState(false);
	const [spacerColor, setSpacerColor] = useState("white");
	const [myTurn, setMyTurn] = useState(false);
	const [handBackground, setHandBackground] = useState("bg-white");

	useEffect(() => {
		const colorMap = {
			red: "bg-red-600",
			blue: "bg-blue-600",
			green: "bg-green-600",
			yellow: "bg-yellow-600",
			orange: "bg-orange-600",
			fuchsia: "bg-fuchsia-600",
			white: "bg-white"
		};
		setHandBackground(colorMap[player?.color]);
	}, [player]);

	useEffect(() => {
		setMyTurn(player?.name === players[turnIndex]?.name);
	}, [turnIndex, player, players]);

	useEffect(() => {
		setSpacerAvailable(hintedSpacer && hintedSpacer.length !== 0);
	}, [hintedSpacer]);

	useEffect(() => {
		if (selectedTile?.a == -1 && selectedTile?.b == -1) {
			setSpacerColor("bg-white");
			return;
		}
		if (spacerAvailable) {
			setSpacerColor("bg-gray-200");
		} else {
			setSpacerColor("bg-gray-500");
		}
	}, [spacerAvailable, selectedTile]);

	const [reactURL, setReactURL] = useState(undefined);
	const [showReaction, setShowReaction] = useState(false);
	useEffect(() => {
		setReactURL(player?.reactURL);
	}, [player]);

	useEffect(() => {
		setShowReaction(reactURL !== undefined);
	}, [reactURL]);
	
	const moveTile = useCallback((tile, toTile) => {
		if (tile.a === toTile.a && tile.b === toTile.b) {
			return;
		}
		let newOrder = [];
		let fromEarlier = false;
		handOrder.forEach(t => {
			if (t.a === toTile.a && t.b === toTile.b) {
				if (fromEarlier) {
					newOrder.push(t);
					newOrder.push(tile);
				} else {
					newOrder.push(tile);
					newOrder.push(t);
				}
				return;
			}
			if (t.a === tile.a && t.b === tile.b) {
				fromEarlier = true;
				return;
			}
			newOrder.push(t);
		});
		setHandOrder(newOrder);
	}, [handOrder, setHandOrder]);

	useEffect(() => {
		if (!player?.hand) {
			setHandOrder([]);
			return;
		}
		const oldTileKeys = new Set(handOrder.map(t => `${t.a}:${t.b}`));
		const newTileKeys = new Set(player.hand.map(t => `${t.a}:${t.b}`));

		let newHandOrder = []
		// old tiles in the order they were, if they're in the new hand.
		handOrder.forEach(t => {
			const key = `${t.a}:${t.b}`;
			if (newTileKeys.has(key)) {
				newHandOrder.push(t);
			}
		});
		// new tiles at the end.
		player.hand.forEach(t => {
			const key = `${t.a}:${t.b}`;
			if (!oldTileKeys.has(key)) {
				newHandOrder.push(t);
			}
		});

		setHandOrder(newHandOrder);
	}, [player]);

	const tileClicked = useCallback((tile) => {
		if (hidden) {
			return;
		}
		setHoveredSquares(new Set([]));
		if (selectedTile === tile) {
			toggleOrientation();
		} else {
			setDragOrientation("down");
		}
		setSelectedTile(tile);
	}, [toggleOrientation, setDragOrientation, selectedTile, setSelectedTile]);

	const spacerClicked = useCallback(() => {
		setHoveredSquares(new Set([]));
		if (hidden) {
			return;
		}
		setSelectedTile({ a: -1, b: -1 });
	}, [setSelectedTile, setHoveredSquares, hidden]);

	useEffect(() => {
		if (setDragOrientation) {
			setDragOrientation("down");
		}
	}, [selectedTile, setDragOrientation]);

	const [isDragging, setIsDragging] = useState(false);

	const handleDragStart = useCallback((tile, e) => {
		if (hidden || e.target !== e.currentTarget) return;
		
		// Create a clone of the entire tile container including its children
		const ghost = e.currentTarget.cloneNode(true);
		ghost.style.position = 'absolute';
		ghost.style.top = '-1000px';
		ghost.style.width = '4rem';
		ghost.style.height = '8rem';
		
		let x_offset = 32;
		let y_offset = 32;
		// Ensure the rotation class is applied
		if (tile.a === selectedTile?.a && tile.b === selectedTile?.b) {
			switch (dragOrientation) {
			case "down": 
				ghost.classList.add("rotate-0");
				break;
			case "right": 
				ghost.classList.add("-rotate-90");
				x_offset = 32;
				y_offset = 64;
				break;
			case "up": 
				ghost.classList.add("rotate-180");
				y_offset = 96;
				break;
			case "left": 
				ghost.classList.add("rotate-90");
				x_offset = 96;
				y_offset = 64;
				break;
			}
		}

		document.body.appendChild(ghost);
		setIsDragging(true);
		

		e.dataTransfer.setDragImage(ghost, x_offset, y_offset);
		e.dataTransfer.setData('text/plain', JSON.stringify(tile));
		setSelectedTile(tile);

		requestAnimationFrame(() => {
			document.body.removeChild(ghost);
		});
	}, [dragOrientation, hidden, selectedTile]);

	const handleDrop = useCallback((targetTile, e) => {
		setIsDragging(false);
		setHoveredSquares(new Set([]));
		if (hidden) return;
		e.preventDefault();
		const sourceTile = JSON.parse(e.dataTransfer.getData('text/plain'));
		// Here you can add logic to swap tiles in the hand
		moveTile(sourceTile, targetTile);
	}, [moveTile, hidden, setIsDragging, setHoveredSquares]);

	useEffect(() => {
		if (!isDragging || mouseIsOver === undefined || mouseIsOver[0] === -1 || mouseIsOver[1] === -1) {
			return;
		}
		hoverTile(mouseIsOver[0], mouseIsOver[1]);
	}, [mouseIsOver]);

	function handleDragOver(e) {
		e.preventDefault();
	}


	function orientGhost(ghost, x, y, orientation) {

		const r = ghost.getBoundingClientRect();
		const min_x = r.left;
		const min_y = r.top;
		const max_x = r.right;
		const max_y = r.bottom;
		const width = max_x - min_x;
		const height = max_y - min_y;

		switch (orientation) {
			case "down":
				ghost.style.left = `${x - width / 2}px`;
				ghost.style.top = `${y - height / 4}px`;
				ghost.style.transform = 'rotate(0deg)';
				break;
			case "right":
				ghost.style.left = `${x - width / 8}px`;
				ghost.style.top = `${y - (3*height/4)}px`;
				ghost.style.transform = 'rotate(270deg)';
				break;
			case "up":
				ghost.style.left = `${x - width / 2}px`;
				ghost.style.top = `${y - height * 3 / 4}px`;
				ghost.style.transform = 'rotate(180deg)';
				break;
			case "left":
				ghost.style.left = `${x - width * 5 / 8}px`;
				ghost.style.top = `${y - (3*height/4)}px`;
				ghost.style.transform = 'rotate(90deg)';
				break;
		}
		
	}

	const handleTouchStart = useCallback((tile, e) => {
		if (hidden) return;
		
		// Create ghost element
		const ghost = e.target.cloneNode(true);
		ghost.id = 'touch-ghost';
		ghost.style.position = 'fixed';
		ghost.style.width = '4rem';
		ghost.style.height = '6rem';
		ghost.style.transform = 'scale(1)';
		ghost.style.opacity = '0.8';
		ghost.style.pointerEvents = 'none';
		ghost.style.zIndex = '1000';
		
		const touch = e.touches[0];
		orientGhost(ghost, touch.clientX, touch.clientY, dragOrientation);
		
		document.body.appendChild(ghost);
		setSelectedTile(tile);
		setTouchStartPos({ x: touch.clientX, y: touch.clientY });
		setDraggedTile(tile);
	}, [dragOrientation, hidden, selectedTile, setSelectedTile, setTouchStartPos, setDraggedTile]);

	const handleTouchEnd = useCallback((targetTile, e) => {
		setHoveredSquares(new Set([]));
		if (!draggedTile || !touchStartPos) return;
		
		// Remove the ghost element
		const ghost = document.getElementById('touch-ghost');
		if (ghost) {
			ghost.parentElement.removeChild(ghost);
		}
		
		// Get the element under the touch point
		const touch = e.changedTouches[0];
		const element = document.elementFromPoint(touch.clientX, touch.clientY);
		
		if (element?.dataset?.tron_x && element?.dataset?.tron_y) {
			dropTile(element.dataset.tron_x, element.dataset.tron_y);
		}

		// Find the tile container element
		const tileContainer = element?.closest('[draggable="true"]');
		if (tileContainer) {
			const endTile = JSON.parse(tileContainer.dataset.tile);
			if (draggedTile !== endTile) {
				moveTile(draggedTile, endTile);
			}
		}
		
		setTouchStartPos(null);
		setDraggedTile(null);
	}, [setTouchStartPos, setDraggedTile]);

	const hoverTile = useCallback((x, y) => {
		if (!selectedTile) {
			return;
		}
		if (setHoveredSquares === undefined) {
			return;
		}
		let hs = new Set([`${x},${y}`]);
		switch (dragOrientation) {
		case "down":
			hs.add(`${x},${y + 1}`);
			break;
		case "right":
			hs.add(`${x + 1},${y}`);
			break;
		case "up":
			hs.add(`${x},${y - 1}`);
			break;
		case "left":
			hs.add(`${x - 1},${y}`);
			break;
		}
		setHoveredSquares(hs);
	}, [selectedTile, dragOrientation, setHoveredSquares]);

	const handleTouchMove = useCallback((e) => {
		if (!touchStartPos) return;
		
		// Move the ghost element
		const ghost = document.getElementById('touch-ghost');
		if (!ghost) {
			return;
		}
		const touch = e.touches[0];
		orientGhost(ghost, touch.clientX, touch.clientY, dragOrientation);

		const element = document.elementFromPoint(touch.clientX, touch.clientY);
		hoverTile(parseInt(element?.dataset?.tron_x), parseInt(element?.dataset?.tron_y));
	}, [dragOrientation, hoverTile]);

	const dropTile = useCallback((x, y) => {
		if (!selectedTile || x === undefined || y === undefined) {
			return;
		}
		playTile({
			a: selectedTile.a, b: selectedTile.b,
			coord: {
				x: parseInt(x),
				y: parseInt(y),
			},
			orientation: dragOrientation,
			dead: false,
		});
	}, [selectedTile, dragOrientation]);

	const [selectedTileRotation, setSelectedTileRotation] = useState("rotate-0");
	useEffect(() => {
		switch (dragOrientation) {
			case "down": setSelectedTileRotation("rotate-0"); break;
			case "right": setSelectedTileRotation("-rotate-90"); break;
			case "up": setSelectedTileRotation("rotate-180"); break;
			case "left": setSelectedTileRotation("rotate-90"); break;
		}
	}, [dragOrientation]);

	const [killedPlayers, setKilledPlayers] = useState([]);
	useEffect(() => {
		setKilledPlayers(player?.kills?.map(k =>  players.find(p => p.name === k)));
	}, [player, players]);
	
	const spacerBundle = useTipBundle("You've got a double that can be used to start a free line. Select it, then choose a square next to a playable line, and choose another square 5 spaces away.");
	useEffect(() => {
		if (hintedSpacer !== null && hintedSpacer.length > 0) {
			spacerBundle.setShow(true);
		}
	}, [hintedSpacer]);

	return (
		<div className={`h-full flex flex-col items-center ${myTurn ? "border-2 border-black " + handBackground : ""}`}>
			<div className="w-full text-center font-bold ">
				<div className="flex flex-row items-center justify-center gap-2">
					{killedPlayers?.map(kp => (
						<div key={kp.name} className="relative w-[2rem] h-[2rem] inline-block align-middle">
							<div className="absolute inset-0">
								<ChickenFoot url={kp.chickenFootURL} color={kp.color} />
							</div>
						</div>
					))}
					<span>
						{player?.name} - ({player?.score})
						{player?.chickenFoot && " (footed)"}
						{player?.ready && " (ready)"}
					</span>
					{showReaction && (
						<Reaction 
							url={reactURL}
							setShow={setShowReaction}
						/>
					)}
					{!player?.chickenFoot && !player?.dead &&
						<div className="relative w-[2rem] h-[2rem] inline-block align-middle">
							<div className="absolute inset-0">
								<ChickenFoot url={player.chickenFootURL} color={player.color} />
							</div>
						</div>
					}
					{!hidden && <div className=" flex items-center gap-2">
						<div className="flex flex-row gap-1 justify-center">
							<Button
								size="small"
								className="w-14"
								disabled={!roundInProgress || !playerTurn || player?.just_drew || bagCount == 0}
								onClick={drawTile}
							>
								draw
							</Button>
							<Button
								size="small"
								className="w-14"
								disabled={!roundInProgress || !playerTurn || !(player?.just_drew || bagCount == 0)}
								onClick={passTurn}
							>
								pass
							</Button>
						</div>;
						<div className="flex flex-row items-center ">
							<Image 
								src="/bag.png" 
								alt="bag"
								width={256}
								height={256}
								className="object-contain w-8 h-8"
							/>
							<div className="text-center">
							{`x${bagCount}`} 
							</div>
						</div>
						<Button
							size="small"
							className="w-14"
							onClick={() => setShowReactModal(true)}>
							react
						</Button>
					</div>}
				</div>
			</div>
			{hidden && (
				<div className="flex flex-row items-center gap-1">
					<div className="w-[1rem]">
						<Tile
							color={player?.color}
							pipsa={0}
							pipsb={0}
							back={true}
							dead={dead}
						/>
					</div>
					<div>x{player?.hand?.length}</div>
				</div>
			)}
			{!hidden && (
				<div className="w-full flex flex-col items-center flex-1 border-1 border-t border-black overflow-y-auto">
					<div className="w-full min-h-[10rem] flex flex-col flex-1">
						<div className="w-full flex flex-row justify-center">
							<div className="w-fit flex flex-wrap content-start justify-start">
								{!hidden && (
									<div className="max-h-[120px] aspect-[1/2] p-1">
										<Tip bundle={spacerBundle} />
										<div
											className={`${spacerColor} ${spacerAvailable && "-translate-y-2"} h-full border-black rounded-lg border-2 flex items-center justify-center text-center`}
											onClick={spacerClicked}
										>
											FREE LINE
										</div>
									</div>
								)}
								{!hidden && handOrder.map((t, i) => {
									const isSelected = playerTurn && selectedTile !== undefined && t.a === selectedTile.a && t.b === selectedTile.b;
									return (
										<div
											key={i}
											draggable={true}
											data-tile={JSON.stringify(t)}
											onClick={() => tileClicked(t)}
											onDragStart={(e) => handleDragStart(t, e)}
											onDrop={(e) => handleDrop(t, e)}
											onDragOver={handleDragOver}
											onTouchStart={(e) => handleTouchStart(t, e)}
											onTouchMove={(e) => handleTouchMove(e)}
											onTouchEnd={(e) => handleTouchEnd(t, e)}
										>
											<div
												className={`max-h-[120px] aspect-[1/2] pr-1 pt-1 ${isSelected ? selectedTileRotation : ""}`}
												>
												<div className="pointer-events-none">
													<Tile
														draggable={false}
														color={player?.color}
														pipsa={t.a}
														pipsb={t.b}
														back={false}
														dead={dead}
														hintedTiles={hintedTiles}
														selected={isSelected}
													/>
												</div>
											</div>
										</div>
									);
								})}
							</div>
							{/* This gutter ensures that a touch can land somewhere to scroll without grabbing a tile. */}
							{!hidden && <div className="w-[1rem]"></div>}
						</div>
					</div>
				</div>
			)}
		</div>
	);
}

export default Hand;