import { useState, useEffect } from "react";
import Tile from '../board/Tile';
import { Button } from "antd";

function Hand({ player, hidden = false, dead = false, selectedTile, setSelectedTile, playerTurn, drawTile, passTurn, roundInProgress, hintedTiles, hintedSpacer, bagCount }) {
	const [handOrder, setHandOrder] = useState([]);
	const [touchStartPos, setTouchStartPos] = useState(null);
	const [draggedTile, setDraggedTile] = useState(null);
	const [spacerAvailable, setSpacerAvailable] = useState(false);
	const [spacerColor, setSpacerColor] = useState("white");

	useEffect(() => {
		setSpacerAvailable(hintedSpacer);
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
		const oldTileKeys = new Set(handOrder);
		const newTileKeys = new Set(player.hand);

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
		
		// Position the ghost at the touch point
		const touch = e.touches[0];
		ghost.style.left = `${touch.clientX - 32}px`;
		ghost.style.top = `${touch.clientY - 48}px`;
		
		document.body.appendChild(ghost);
		setTouchStartPos({ x: touch.clientX, y: touch.clientY });
		setDraggedTile(tile);
	}

	function handleTouchMove(e) {
		if (!touchStartPos) return;
		e.preventDefault();
		
		// Move the ghost element
		const ghost = document.getElementById('touch-ghost');
		if (ghost) {
			const touch = e.touches[0];
			ghost.style.left = `${touch.clientX - 32}px`;
			ghost.style.top = `${touch.clientY - 48}px`;
		}
	}

	function handleTouchEnd(targetTile, e) {
		if (!draggedTile || !touchStartPos) return;
		
		// Remove the ghost element
		const ghost = document.getElementById('touch-ghost');
		if (ghost) {
			ghost.remove();
		}
		
		// Get the element under the touch point
		const touch = e.changedTouches[0];
		const element = document.elementFromPoint(touch.clientX, touch.clientY);
		
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
	}

	// Add useEffect for touch event setup
	useEffect(() => {
		// Get all draggable tile elements
		const tileElements = document.querySelectorAll('[draggable="true"]');
		
		// Add non-passive touch move listeners
		tileElements.forEach(element => {
			element.addEventListener('touchmove', handleTouchMove, { passive: false });
		});

		// Cleanup
		return () => {
			tileElements.forEach(element => {
				element.removeEventListener('touchmove', handleTouchMove);
			});
		};
	}, [touchStartPos]); // Re-run when touchStartPos changes

	return (
		<div className="h-full flex flex-col items-center p-2">
			<div className="text-center font-bold">
				{player?.name} - ({player?.score}) {player?.chickenFoot && "(footed)"}
			</div>
			<div className="flex flex-col items-center h-full">
				<div className="overflow-y-auto max-h-[calc(100%-14rem)]">
					{!hidden && (
						<div className="flex justify-center">
							<div
								className={`${spacerColor} w-[24rem] h-[4rem] border-black rounded-lg border-2 flex items-center justify-center text-center`}
								onClick={spacerClicked}
							>
								FREE LINE SPACER
							</div>
							
				{!hidden && (
					<div>
						<div className="text-center">
							{`${bagCount} tile${bagCount === 1 ? "" : "s"} in the bag`} 
						</div>
						<div className="flex flex-row gap-1 ml-4 mt-2">
							<Button
								type="primary"
								size="large"
								className="w-14"
								disabled={!roundInProgress || !playerTurn || player?.just_drew}
								onClick={drawTile}
							>
								Draw
							</Button>
							<Button
								type="primary"
								size="large"
								className="w-14"
								disabled={!roundInProgress || !playerTurn || !player?.just_drew}
								onClick={passTurn}
							>
								Pass
							</Button>
						</div>
					</div>
				)}
						</div>
					)}
					<div className="flex flex-wrap content-start">
						{handOrder.map((t, i) => {
							return (
								<div
									key={i}
									className={hidden ? "w-[1rem]" : "w-[4rem] pr-1 pt-1"}
									draggable={!hidden}
									data-tile={JSON.stringify(t)}
									onClick={() => tileClicked(t)}
									onDragStart={(e) => handleDragStart(t, e)}
									onDrop={(e) => handleDrop(t, e)}
									onDragOver={handleDragOver}
									onTouchStart={(e) => handleTouchStart(t, e)}
									onTouchEnd={(e) => handleTouchEnd(t, e)}
								>
									<div className="pointer-events-none">
										<Tile
											draggable={false}
											color={player?.color}
											pipsa={t.a}
											pipsb={t.b}
											back={hidden}
											dead={dead}
											hintedTiles={hintedTiles}
											selected={playerTurn && selectedTile !== undefined && t.a === selectedTile.a && t.b === selectedTile.b}
										/>
									</div>
								</div>
							);
						})}
					</div>
				</div>
			</div>
		</div>
	);
}

export default Hand;