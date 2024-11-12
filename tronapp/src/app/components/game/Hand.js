import Tile from '../board/Tile';
import { Button } from "antd";

function Hand({ player, hidden = false, dead = false, selectedTile, setSelectedTile, playerTurn, drawTile, passTurn, roundInProgress, hintedTiles }) {
	function tileClicked(tile) {
		if (hidden) {
			return;
		}
		setSelectedTile(tile);
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
		console.log('Dragged', sourceTile, 'to', targetTile);
	}

	function handleDragOver(e) {
		e.preventDefault();
	}

	return <div className="h-full flex flex-col items-center p-2">
		<div className="text-center font-bold">
			{player?.name} - ({player?.score}) {player?.chickenFoot && "(footed)"}
		</div>
		<div className="flex items-start justify-center h-full">
			<div className="overflow-y-auto max-h-[calc(100%-2rem)]">
				<div className="flex flex-wrap content-start">
					{player?.hand?.map((t, i) => {
						return (
							<div
								key={i}
								className={hidden ? "w-[1rem]" : "w-[4rem] pr-1 pt-1"}
								draggable={!hidden}
								onClick={() => tileClicked(t)}
								onDragStart={(e) => handleDragStart(t, e)}
								onDrop={(e) => handleDrop(t, e)}
								onDragOver={handleDragOver}
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
			{!hidden && <div className="flex flex-col ml-4 mt-2">
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
			</div>}
		</div>
	</div>;
}

export default Hand;