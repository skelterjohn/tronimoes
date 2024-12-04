
import Tip, { makeTipBundle } from "@/app/components/tutorial/Tip";
import { useEffect } from "react";

export default function Hint() {
	const hintBundle = makeTipBundle("This white square indicates the tile you have selected can be played here.");

	useEffect(() => {
		hintBundle.setShow(true);
	}, [hintBundle]);

	return (
        <div ref={hintBundle.parentRef} className="w-full h-full flex items-center justify-center">
			<Tip bundle={hintBundle} />
            <div className={`w-[90%] h-[90%] border-white border-4 border`}></div>
        </div>
    )
}
