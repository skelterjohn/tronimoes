"use client";

import { useState, useEffect } from "react";

export default function Spacer({ spacer }) {
	const [rotate, setRotate] = useState("");
	useEffect(() => {
		if (!spacer) {
			return;
		}
		if (spacer.x2 > spacer.x1) {
			setRotate("-rotate-90");
		}
		if (spacer.x2 < spacer.x1) {
			setRotate("rotate-90");
		}
		if (spacer.y2 > spacer.y1) {
			setRotate("");
		}
		if (spacer.y2 < spacer.y1) {
			setRotate("rotate-180");
		}
	}, [spacer]);

	return <div className={`absolute w-full h-full ${rotate}`}>
		<div className="h-[600%] bg-white border-black border-2 rounded-lg flex items-center justify-center">
			<div className="rotate-90 whitespace-nowrap absolute transform origin-center">
				RIGHT-CLICK TO CLEAR
			</div>
		</div>
	</div>
}
