"use client";

import { useState, useEffect } from "react";

export default function Spacer({ spacer }) {
	const [rotate, setRotate] = useState("");
	useEffect(() => {
		if (!spacer) {
			return;
		}
		if (spacer.b.x > spacer.a.x) {
			setRotate("-rotate-90");
		}
		if (spacer.b.x < spacer.a.x) {
			setRotate("rotate-90");
		}
		if (spacer.b.y > spacer.a.y) {
			setRotate("");
		}
		if (spacer.b.y < spacer.a.y) {
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
