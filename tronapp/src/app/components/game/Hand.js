import Tile from '../board/Tile';

function Hand({}) {
	var tiles = [
		{a:1, b:2},
		{a:3, b:12},
	]
	return <div className="flex items-center justify-center">{tiles.map((t) => (
		<div key={`${t.a}-${t.b}`} className="max-w-[10%] aspect-square">
			<Tile pipsa={t.a} pipsb={t.b} />
		</div>
	))}</div>;
}

export default Hand;