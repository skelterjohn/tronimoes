import Tile from '../board/Tile';
import { Button } from "antd";
import { useEffect } from "react";

function Hand({ player, hidden=false, dead=false, selectedTile, setSelectedTile, playerTurn, drawTile }) {
	function tileClicked(tile) {
		if (hidden) {
			return;
		}
		setSelectedTile(tile);
	}
	console.log(player);
	return <div className="h-full flex flex-col items-center">
		<div className="text-center font-bold">{player?.name} - ({player?.score}) {player?.chickenFoot && "(footed)"}</div>
		<div className="flex items-center justify-center overflow-hidden">
			{player?.hand?.map((t, i) => (
				<div key={i} className="w-[4rem]" onClick={()=>tileClicked(t)}>
					<Tile
						color={player?.color}
						pipsa={t.a}
						pipsb={t.b}
						back={hidden}
						dead={dead}
						selected={playerTurn && selectedTile!==undefined && t.a===selectedTile.a && t.b===selectedTile.b}
					/>
				</div>
			))}
			{!hidden && <span>
				<Button
					disabled={!playerTurn}
					onClick={drawTile}
				>
					Draw
				</Button>
			</span>}
		</div>
	</div>;
}

export default Hand;