"use client";

import { useState, useEffect } from "react";
import Tip, { useTipBundle } from "@/app/components/tutorial/Tip";

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
	
	const freeLineBundle = useTipBundle("Now that you've got the spacer placed, you can use a double to start a free line at the end of the spacer. It has to be higher than any other leader on the board.");
	useEffect(() => {
		freeLineBundle.setShow(true);
	}, []);

	return <div className={`absolute w-full h-full ${rotate}`}>
		<Tip bundle={freeLineBundle} />
		<div className="h-[600%] bg-white border-black border-2 rounded-lg flex items-center justify-center">
			<div className="rotate-90 whitespace-nowrap text-black absolute transform origin-center">
				RIGHT-CLICK TO CLEAR
			</div>
		</div>
	</div>
}
