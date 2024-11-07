import Tile from '../board/Tile';
import { Button } from "antd";

function Hand({ player, hidden=false, dead=false, selectedTile, setSelectedTile, playerTurn, drawTile, passTurn, roundInProgress, hintedTiles }) {
	function tileClicked(tile) {
		if (hidden) {
			return;
		}
		setSelectedTile(tile);
	}
	return <div className="h-full flex flex-col items-center p-2">
		<div className="text-center font-bold">{player?.name} - ({player?.score}) {player?.chickenFoot && "(footed)"}</div>
		<div className="flex items-center justify-center overflow-hidden p-2">
			{player?.hand?.map((t, i) => {
				return (
					<div key={i} className="w-[4rem]" onClick={()=>tileClicked(t)}>
						<Tile
							color={player?.color}
							pipsa={t.a}
							pipsb={t.b}
							back={hidden}
							dead={dead}
							hintedTiles={hintedTiles}
							selected={playerTurn && selectedTile!==undefined && t.a===selectedTile.a && t.b===selectedTile.b}
						/>
					</div>
				);
			})}
			{!hidden && <div className="flex flex-col gap-2">
				<Button
					type="primary"
					size="large"
					disabled={!roundInProgress || !playerTurn || player?.just_drew}
					onClick={drawTile}
				>
					Draw
				</Button>
				<Button
					type="primary"
					size="large"
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