function Square({ x, y, center=false }) {
	var cnm = "w-full aspect-square";
	if (center) {
		cnm = `${cnm} bg-gray-400`;
	} else if ((x+y) % 2 == 0) {
		cnm = `${cnm} bg-blue-200`;
	} else {
		cnm = `${cnm} bg-slate-200`;
	}
	return <div className={cnm}> </div>;
}

export default Square;