import { useState, useEffect } from "react";
import Tile from '../board/Tile';
import ChickenFoot from '../board/ChickenFoot';
import { Button } from "antd";

function Hand({ player, players, hidden = false, dead = false, selectedTile, setSelectedTile, playerTurn, drawTile, passTurn, roundInProgress, hintedTiles, hintedSpacer, bagCount, turnIndex, playTile, setHoveredSquares }) {
	const [handOrder, setHandOrder] = useState([]);
	const [touchStartPos, setTouchStartPos] = useState(null);
	const [draggedTile, setDraggedTile] = useState(null);
	const [spacerAvailable, setSpacerAvailable] = useState(false);
	const [spacerColor, setSpacerColor] = useState("white");
	const [myTurn, setMyTurn] = useState(false);
	const [handBackground, setHandBackground] = useState("bg-white");
	const [dragOrientation, setDragOrientation] = useState("down");

	function toggleOrientation() {
		switch (dragOrientation) {
			case "down": setDragOrientation("right"); break;
			case "right": setDragOrientation("up"); break;
			case "up": setDragOrientation("left"); break;
			case "left": setDragOrientation("down"); break;
		}
	}

	useEffect(() => {
		console.log(dragOrientation);
	}, [dragOrientation]);

	useEffect(() => {
		const colorMap = {
			red: "bg-red-100",
			blue: "bg-blue-100",
			green: "bg-green-100",
			yellow: "bg-yellow-100",
			orange: "bg-orange-100",
			fuchsia: "bg-fuchsia-100",
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

	function moveTile(tile, toTile) {
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
	}

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

	function tileClicked(tile) {
		if (hidden) {
			return;
		}
		if (selectedTile === tile) {
			toggleOrientation();
		} else {
			setDragOrientation("down");
		}
		setSelectedTile(tile);
	}

	function spacerClicked() {
		if (hidden) {
			return;
		}
		setSelectedTile({ a: -1, b: -1 });
	}

	function handleDragStart(tile, e) {
		if (hidden || e.target !== e.currentTarget) return;
		
		// Create a clone of the tile being dragged
		const dragImage = e.target.cloneNode(true);
		dragImage.style.position = 'absolute';
		dragImage.style.top = '-1000px';
		dragImage.style.width = '4rem';
		dragImage.style.height = '6rem';
		dragImage.style.transform = 'scale(1)';
		document.body.appendChild(dragImage);
		
		// Set the drag image
		e.dataTransfer.setDragImage(dragImage, 32, 48);
		e.dataTransfer.setData('text/plain', JSON.stringify(tile));
		setSelectedTile(tile);

		// Remove the clone after the drag starts
		requestAnimationFrame(() => {
			document.body.removeChild(dragImage);
		});
	}

	function handleDrop(targetTile, e) {
		if (hidden) return;
		e.preventDefault();
		const sourceTile = JSON.parse(e.dataTransfer.getData('text/plain'));
		// Here you can add logic to swap tiles in the hand
		moveTile(sourceTile, targetTile);
	}

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
				ghost.style.left = `${x - width / 4}px`;
				ghost.style.top = `${y - height / 2}px`;
				ghost.style.transform = 'rotate(270deg)';
				break;
			case "up":
				ghost.style.left = `${x - width / 2}px`;
				ghost.style.top = `${y - height * 3 / 4}px`;
				ghost.style.transform = 'rotate(180deg)';
				break;
			case "left":
				ghost.style.left = `${x - width * 3 / 4}px`;
				ghost.style.top = `${y - height / 2}px`;
				ghost.style.transform = 'rotate(90deg)';
				break;
		}
		
	}

	function handleTouchStart(tile, e) {
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
		setTouchStartPos({ x: touch.clientX, y: touch.clientY });
		setDraggedTile(tile);
	}

	function handleTouchEnd(targetTile, e) {
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
			dropTile(element.dataset.tron_x, element.dataset.tron_y, dragOrientation);
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
		setHoveredSquares(new Set([]));
	}

	function handleTouchMove(e) {
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
	}

	function hoverTile(x, y) {
		if (!selectedTile) {
			return;
		}
		setHoveredSquares(`${x},${y}`);
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
	}

	function dropTile(x, y, orientation) {
		if (!selectedTile) {
			return;
		}
		console.log("drop", x, y, orientation, selectedTile);
		playTile({
			a: selectedTile.a, b: selectedTile.b,
			coord: {
				x: parseInt(x),
				y: parseInt(y),
			},
			orientation: orientation,
			dead: false,
		});
	}

	const [killedPlayers, setKilledPlayers] = useState([]);
	useEffect(() => {
		setKilledPlayers(player?.kills?.map(k =>  players.find(p => p.name === k)));
	}, [player, players]);

	function DrawPassButtons() {
		return <div className="flex flex-row gap-1 justify-center">
			<Button
				type="primary"
				size="small"
				className="w-14"
				disabled={!roundInProgress || !playerTurn || player?.just_drew || bagCount == 0}
				onClick={drawTile}
			>
				Draw
			</Button>
			<Button
				type="primary"
				size="small"
				className="w-14"
				disabled={!roundInProgress || !playerTurn || !(player?.just_drew || bagCount == 0)}
				onClick={passTurn}
			>
				Pass
			</Button>
		</div>;
	}

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
						{player?.name} - ({player?.score}) {player?.chickenFoot && "(footed)"}
					</span>
					{!player?.chickenFoot && !player?.dead &&
						<div className="relative w-[2rem] h-[2rem] inline-block align-middle">
							<div className="absolute inset-0">
								<ChickenFoot url={player.chickenFootURL} color={player.color} />
							</div>
						</div>
					}
					{!hidden && <div className=" flex items-center gap-2">
						<DrawPassButtons/>
						<div className="text-center">
							{`${bagCount} tile${bagCount === 1 ? "" : "s"} in the bag`} 
						</div>
					</div>}
				</div>
			</div>
			<div className="w-full flex flex-col items-center flex-1 min-h-0 border-1 border-t border-black">
				<div className="w-full flex flex-col flex-1 overflow-y-auto">
					<div className="w-full flex flex-row justify-center">
						<div className="w-fit flex flex-wrap content-start justify-start">
							{!hidden && (
								<div className="w-[4rem] h-[8rem] p-1">
									<div
										className={`${spacerColor} ${spacerAvailable && "-translate-y-2"} h-full border-black rounded-lg border-2 flex items-center justify-center text-center`}
										onClick={spacerClicked}
									>
										FREE LINE
									</div>
								</div>
							)}
							{!hidden && handOrder.map((t, i) => {
								return (
									<div
										key={i}
										className="w-[4rem] h-[8rem] pr-1 pt-1"
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
										<div className="pointer-events-none">
											<Tile
												draggable={false}
												color={player?.color}
												pipsa={t.a}
												pipsb={t.b}
												back={false}
												dead={dead}
												hintedTiles={hintedTiles}
												selected={playerTurn && selectedTile !== undefined && t.a === selectedTile.a && t.b === selectedTile.b}
											/>
										</div>
									</div>
								);
							})}
						</div>
						{/* This gutter ensures that a touch can land somewhere to scroll without grabbing a tile. */}
						{!hidden && <div className="w-[1rem]"></div>}
						{hidden && (
							<div className="w-[1rem]">
								<Tile
									color={player?.color}
									pipsa={0}
									pipsb={0}
									back={true}
									dead={dead}
								/>
							</div>
						)}
						{hidden && (
							<div>x{player?.hand?.length}</div>
						)}
					</div>
				</div>
			</div>
		</div>
	);
}

export default Hand;