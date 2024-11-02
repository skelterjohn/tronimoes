import Tile from '../board/Tile';
import { Button } from "antd";
import { useEffect } from "react";

function Hand({ name, color="white", hidden=false, tiles=[], dead=false, selectedTile, setSelectedTile, playerTurn, drawTile, score }) {
	function tileClicked(tile) {
		if (hidden) {
			return;
		}
		setSelectedTile(tile);
	}
	return <div className="h-full flex flex-col items-center">
		<div className="text-center font-bold">{name} - ({score})</div>
		<div className="flex items-center justify-center overflow-hidden">
			{tiles.map((t, i) => (
				<div key={i} className="w-[4rem]" onClick={()=>tileClicked(t)}>
					<Tile
						color={color}
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