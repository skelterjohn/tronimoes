import Pips from "./Pips";

export default function Square({ x, y, center=false, clicked=false, pips=0 }) {
	var cnm = "w-full aspect-square";
	if (center) {
		cnm = `${cnm} bg-gray-400`;
	} else if ((x+y) % 2 == 0) {
		cnm = `${cnm} bg-blue-200`;
	} else {
		cnm = `${cnm} bg-slate-200`;
	}

	if (clicked) {
		cnm = `${cnm} border border-2 border-black`;
	}

	return <div className={cnm}>
		{clicked && <Pips pips={pips} />}
	</div>;
}
