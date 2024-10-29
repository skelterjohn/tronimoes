import Tile from '../board/Tile';

function Hand({ color="white", hidden=false, tiles=[], dead=false }) {
	return <div className="flex items-center justify-center">{tiles.map((t, i) => (
		<div key={i} className="max-w-[10%] ">
			<Tile color={color} pipsa={t.a} pipsb={t.b} back={hidden} dead={dead} />
		</div>
	))}</div>;
}

export default Hand;