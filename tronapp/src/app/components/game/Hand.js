import { useState, useEffect, useCallback, useRef } from "react";
import Tile from '../board/Tile';
import ChickenFoot from '../board/ChickenFoot';
import Button from "@/app/components/Button";
import Image from "next/image";
import Reaction from "./Reaction";
import Tip, { useTipBundle } from "@/app/components/tutorial/Tip";

function Hand({
		player, players,
		dead = false,
		selectedTile, setSelectedTile,
		playerTurn,
		drawTile, passTurn,
		roundInProgress,
		hintedTiles, hintedSpacer,
		bagCount, turnIndex, playTile,
		setHoveredSquares, mouseIsOver,
		dragOrientation, setDragOrientation, toggleOrientation,
		setShowReactModal,
		boardRef, squareSpan
	}) {
	const [handOrder, setHandOrder] = useState([]);
	const [touchStartPos, setTouchStartPos] = useState(null);
	const [draggedTile, setDraggedTile] = useState(null);
	const [spacerAvailable, setSpacerAvailable] = useState(false);
	const [spacerColor, setSpacerColor] = useState("white");
	const [myTurn, setMyTurn] = useState(false);
	const [handBackground, setHandBackground] = useState("bg-white");
	const scrollContainerRef = useRef(null);

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
		setHoveredSquares(new Set([]));
		setTouchOverBoard(false);
		if (selectedTile === tile) {
			toggleOrientation();
		} else {
			setDragOrientation("down");
		}
		setSelectedTile(tile);
	}, [toggleOrientation, setDragOrientation, selectedTile, setSelectedTile]);

	const spacerClicked = useCallback(() => {
		setHoveredSquares(new Set([]));
		setSelectedTile({ a: -1, b: -1 });
	}, [setSelectedTile, setHoveredSquares]);

	useEffect(() => {
		if (setDragOrientation) {
			setDragOrientation("down");
		}
		setHoveredSquares(new Set([]));
	}, [selectedTile, setDragOrientation]);

	const [isDragging, setIsDragging] = useState(false);


	const handleDragStart = useCallback((tile, e) => {
		if (e.target !== e.currentTarget) return;
		
		// Create a clone of the entire tile container including its children
		const ghost = e.currentTarget.cloneNode(true);
		
		// Reset any existing size-related styles and classes
		ghost.style.position = 'absolute';
		ghost.style.top = '-1000px';
		ghost.style.width = `${squareSpan}px`;
		
		ghost.style.height = `${squareSpan * 2}px`;
		ghost.style.maxHeight = `${squareSpan * 2}px`; // Force max height
		ghost.style.minHeight = `${squareSpan * 2}px`; // Force min height
		ghost.style.aspectRatio = '1/2';
		ghost.style.padding = '0'; // Remove any padding
		ghost.style.margin = '0'; // Remove any margin
		
		// Force the inner tile to match the container size
		const innerTile = ghost.querySelector('div');
		if (innerTile) {
			innerTile.style.width = '100%';
			innerTile.style.height = '100%';
			innerTile.style.maxHeight = '100%';
		}
		
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
	}, [dragOrientation, selectedTile, squareSpan]);

	const handleDrop = useCallback((targetTile, e) => {
		setIsDragging(false);
		setHoveredSquares(new Set([]));
		e.preventDefault();
		const sourceTile = JSON.parse(e.dataTransfer.getData('text/plain'));
		// Here you can add logic to swap tiles in the hand
		moveTile(sourceTile, targetTile);
	}, [moveTile, setIsDragging, setHoveredSquares]);

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

	const [touchOverBoard, setTouchOverBoard] = useState(false);

	const handleTouchStart = useCallback((tile, e) => {
		// Create ghost element
		const ghost = e.target.cloneNode(true);
		ghost.id = 'touch-ghost';
		ghost.style.position = 'fixed';
		ghost.style.width = `${squareSpan}px`;
		ghost.style.height = `${squareSpan * 2}px`;
		ghost.style.maxHeight = `${squareSpan * 2}px`;
		ghost.style.minHeight = `${squareSpan * 2}px`;
		ghost.style.aspectRatio = '1/2';
		ghost.style.padding = '0';
		ghost.style.margin = '0';
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
		setTouchOverBoard(false);
	}, [dragOrientation, selectedTile, setSelectedTile, setTouchStartPos, setDraggedTile, setTouchOverBoard, squareSpan]);


	const handleTouchEnd = useCallback((targetTile, e) => {
		// Remove the ghost element and ensure cleanup
		cleanupGhostElement();
		document.body.style.overflow = ''; // Explicitly restore scrolling

		setHoveredSquares(new Set([]));
		if (!draggedTile || !touchStartPos) return;
		
		// Get the element under the touch point
		const touch = e.changedTouches[0];
		let [x, y] = [touch.clientX, touch.clientY];

		if (touchOverBoard) {
			y -= 100;
		}
		const element = document.elementFromPoint(x, y);
		
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
		setTouchOverBoard(false);
	}, [setTouchStartPos, setDraggedTile, touchOverBoard, setTouchOverBoard]);

	// Update cleanup function to be more aggressive
	const cleanupGhostElement = () => {
		// Find all elements with id starting with 'touch-ghost'
		const ghosts = document.querySelectorAll('[id^="touch-ghost"]');
		ghosts.forEach(ghost => {
			try {
				ghost.remove();
			} catch (e) {
				// If direct removal fails, try removing via parent
				if (ghost.parentElement) {
					ghost.parentElement.removeChild(ghost);
				}
			}
		});
	};

	// Add cleanup to component unmount
	useEffect(() => {
		return () => {
			cleanupGhostElement();
		};
	}, []);

	// Add new touch cancel handler
	const handleTouchCancel = useCallback(() => {
		setHoveredSquares(new Set([]));
		cleanupGhostElement();
		setTouchStartPos(null);
		setDraggedTile(null);
	}, []);

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
		if (!touchStartPos) {
			cleanupGhostElement();
			return;
		}

		const touch = e.touches[0];
		let [x, y] = [touch.clientX, touch.clientY];

		if (touchOverBoard) {
			y -= 100;
		}

		const elementUnderTouch = document.elementFromPoint(x, y);
		
		// Check if the element under touch is within the board
		if (!touchOverBoard && boardRef.current?.contains(elementUnderTouch)) {
			setTouchOverBoard(true);
		}
		const ghost = document.getElementById('touch-ghost');
		if (!ghost) {
			cleanupGhostElement();
			return;
		}
		orientGhost(ghost, x, y, dragOrientation);

		hoverTile(parseInt(elementUnderTouch?.dataset?.tron_x), parseInt(elementUnderTouch?.dataset?.tron_y));
	}, [dragOrientation, hoverTile, touchOverBoard, setTouchOverBoard, boardRef]);

	useEffect(() => {
		if (draggedTile) {
			const preventDefault = (e) => e.preventDefault();
			document.body.style.overflow = 'hidden';
			document.addEventListener('touchmove', preventDefault, { passive: false });
			
			return () => {
				document.body.style.overflow = '';
				document.removeEventListener('touchmove', preventDefault);
			};
		} else {
			// Ensure we restore scrolling when not dragging
			document.body.style.overflow = '';
		}
	}, [draggedTile]);

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
		if (hintedSpacer && hintedSpacer.length > 0) {
			spacerBundle.setShow(true);
		}
	}, [hintedSpacer]);

	const scrollToTop = useCallback(() => {
		setSelectedTile(undefined);
		scrollContainerRef.current?.scrollTo({ top: 0, behavior: 'smooth' });
	}, []);
	
	const scrollToBottom = useCallback(() => {
		setSelectedTile(undefined);
		const container = scrollContainerRef.current;
		container?.scrollTo({ top: container?.scrollHeight, behavior: 'smooth' });
	}, []);
	
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
					<div className=" flex items-center gap-2">
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
						</div>
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
					</div>
					<div className="flex flex-row gap-1">
						<a onClick={scrollToTop} className="cursor-pointer">
							<i className="w-6 aspect-square text-black fa-solid fa-arrow-up"></i>
						</a>
						<a onClick={scrollToBottom} className="cursor-pointer">
							<i className="w-6 aspect-square text-black fa-solid fa-arrow-down"></i>
						</a>
					</div>
				</div>
			</div>
			<div ref={scrollContainerRef} className="w-full flex flex-col items-center flex-1 border-1 border-t border-black overflow-y-auto">
				<div className="w-full min-h-[10rem] flex flex-col flex-1">
					<div className="w-full flex flex-row justify-center">
						<div className="w-[calc(100%-1rem)] flex flex-wrap content-start justify-start">
							<div className="max-h-[15vh] aspect-[1/2] p-1">
								<Tip bundle={spacerBundle} />
								<div
									className={`${spacerColor} ${spacerAvailable && "-translate-y-2"} h-full border-black rounded-lg border-2 flex items-center justify-center text-center`}
									onClick={spacerClicked}
								>
									FREE LINE
								</div>
							</div>
							{handOrder.map((t, i) => {
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
										onTouchCancel={handleTouchCancel}
									>
										<div
											className={`max-h-[15vh] aspect-[1/2] pr-1 pt-1 ${isSelected ? selectedTileRotation : ""}`}
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
					</div>
				</div>
			</div>
		</div>
	);
}

export default Hand;