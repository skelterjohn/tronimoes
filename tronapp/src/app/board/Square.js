function Square({ x, y }) {

	var cnm = "w-full aspect-square";
	if ((x+y) % 2 == 0) {
		cnm = `${cnm} bg-amber-200`;
	} else {
		cnm = `${cnm} bg-slate-200`;
	}
	return <div className={cnm}> </div>;
}

export default Square;