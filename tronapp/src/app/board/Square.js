function Square({ x, y, center=false }) {

	var intensity = 200;
	if (center) {
		intensity = 400;
	}

	var cnm = "w-full aspect-square";
	if ((x+y) % 2 == 0) {
		cnm = `${cnm} bg-blue-${intensity}`;
	} else {
		cnm = `${cnm} bg-slate-${intensity}`;
	}
	return <div className={cnm}> </div>;
}

export default Square;