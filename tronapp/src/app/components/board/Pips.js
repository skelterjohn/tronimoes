function Pip({color="red"}) {
	const colorMap = {
        red: "bg-red-500",
        blue: "bg-blue-500",
        green: "bg-green-500",
        indigo: "bg-indigo-500",
        orange: "bg-orange-500",
        yellow: "bg-yellow-500",
        purple: "bg-purple-500",
        white: "bg-white"
    };
	return <td className="w-[15%] aspect-square">
		<div className={`p-1 aspect-square rounded-full ${colorMap[color]}`}></div>
	</td>;
}

function Pips({pips}) {
	const defaultcolor = "white";
	var c00 = defaultcolor;
	var c10 = defaultcolor;
	var c20 = defaultcolor;
	var c30 = defaultcolor;
	var c40 = defaultcolor;
	var c01 = defaultcolor;
	var c11 = defaultcolor;
	var c21 = defaultcolor;
	var c31 = defaultcolor;
	var c41 = defaultcolor;
	var c02 = defaultcolor;
	var c12 = defaultcolor;
	var c22 = defaultcolor;
	var c32 = defaultcolor;
	var c42 = defaultcolor;
	var c03 = defaultcolor;
	var c13 = defaultcolor;
	var c23 = defaultcolor;
	var c33 = defaultcolor;
	var c43 = defaultcolor;
	var c04 = defaultcolor;
	var c14 = defaultcolor;
	var c24 = defaultcolor;
	var c34 = defaultcolor;
	var c44 = defaultcolor;
	if (pips == 1) {
		const c = "blue";
		c22 = c;
	}
	if (pips == 2) {
		const c = "red";
		c12 = c;
		c32 = c;
	}
	if (pips == 3) {
		const c = "green";
		c11 = c;
		c22 = c;
		c33 = c;
	}
	if (pips == 4) {
		const c = "indigo";
		c11 = c;
		c13 = c;
		c31 = c;
		c33 = c;
	}
	if (pips == 5) {
		const c = "orange";
		c11 = c;
		c13 = c;
		c22 = c;
		c31 = c;
		c33 = c;
	}
	if (pips == 6) {
		const c = "yellow";
		c11 = c;
		c12 = c;
		c13 = c;
		c31 = c;
		c32 = c;
		c33 = c;
	}
	if (pips == 7) {
		const c = "purple";
		c11 = c;
		c12 = c;
		c13 = c;
		c22 = c;
		c31 = c;
		c32 = c;
		c33 = c;
	}
	if (pips == 8) {
		const c = "blue";
		c11 = c;
		c12 = c;
		c13 = c;
		c21 = c;
		c23 = c;
		c31 = c;
		c32 = c;
		c33 = c;
	}
	if (pips == 9) {
		const c = "red";
		c11 = c;
		c12 = c;
		c13 = c;
		c21 = c;
		c22 = c;
		c23 = c;
		c31 = c;
		c32 = c;
		c33 = c;
	}
	if (pips == 10) {
		const c = "green";
		c10 = c;
		c11 = c;
		c12 = c;
		c13 = c;
		c14 = c;
		c30 = c;
		c31 = c;
		c32 = c;
		c33 = c;
		c34 = c;
	}
	if (pips == 11) {
		const c = "indigo";
		c10 = c;
		c11 = c;
		c12 = c;
		c13 = c;
		c14 = c;
		c22 = c;
		c30 = c;
		c31 = c;
		c32 = c;
		c33 = c;
		c34 = c;
	}
	if (pips == 12) {
		const c = "orange";
		c10 = c;
		c11 = c;
		c12 = c;
		c13 = c;
		c14 = c;
		c21 = c;
		c23 = c;
		c30 = c;
		c31 = c;
		c32 = c;
		c33 = c;
		c34 = c;
	}
	if (pips == 13) {
		const c = "yellow";
		c10 = c;
		c11 = c;
		c12 = c;
		c13 = c;
		c14 = c;
		c20 = c;
		c22 = c;
		c24 = c;
		c30 = c;
		c31 = c;
		c32 = c;
		c33 = c;
		c34 = c;
	}
	if (pips == 14) {
		const c = "purple";
		c10 = c;
		c11 = c;
		c12 = c;
		c13 = c;
		c14 = c;
		c02 = c;
		c21 = c;
		c23 = c;
		c42 = c;
		c30 = c;
		c31 = c;
		c32 = c;
		c33 = c;
		c34 = c;
	}
	if (pips == 15) {
		const c = "blue";
		c10 = c;
		c11 = c;
		c12 = c;
		c13 = c;
		c14 = c;
		c02 = c;
		c21 = c;
		c22 = c;
		c23 = c;
		c42 = c;
		c30 = c;
		c31 = c;
		c32 = c;
		c33 = c;
		c34 = c;
	}
	if (pips == 16) {
		const c = "red";
		c20 = c;
		c01 = c;
		c10 = c;
		c11 = c;
		c03 = c;
		c24 = c;
		c13 = c;
		c14 = c;
		c30 = c;
		c31 = c;
		c02 = c;
		c41 = c;
		c33 = c;
		c34 = c;
		c43 = c;
		c42 = c;
	}
	return (
		<table className="w-full aspect-square table-fixed">
			<tbody className="h-full w-full">
				<tr className="w-full">
					<Pip color={c00}/>
					<Pip color={c10}/>
					<Pip color={c20}/>
					<Pip color={c30}/>
					<Pip color={c40}/>
				</tr>
				<tr className="w-full">
					<Pip color={c01}/>
					<Pip color={c11}/>
					<Pip color={c21}/>
					<Pip color={c31}/>
					<Pip color={c41}/>
				</tr>
				<tr className="w-full">
					<Pip color={c02}/>
					<Pip color={c12}/>
					<Pip color={c22}/>
					<Pip color={c32}/>
					<Pip color={c42}/>
				</tr>
				<tr className="w-full">
					<Pip color={c03}/>
					<Pip color={c13}/>
					<Pip color={c23}/>
					<Pip color={c33}/>
					<Pip color={c43}/>
				</tr>
				<tr className="w-full">
					<Pip color={c04}/>
					<Pip color={c14}/>
					<Pip color={c24}/>
					<Pip color={c34}/>
					<Pip color={c44}/>
				</tr>
			</tbody>
		</table>
	);
}

export default Pips;
export {safelist}