
import Tip, { useTipBundle } from "@/app/components/tutorial/Tip";
import { useEffect } from "react";

export default function Hint() {
	const hintBundle = useTipBundle("This white square indicates the tile you have selected can be played here.");
	useEffect(() => {
		hintBundle.setShow(true);
	}, []);
	const clickBundle = useTipBundle("You can also click the squares. First click where the top of the tile goes, then the bottom.");
	useEffect(() => {
		if (hintBundle.done()) {
			clickBundle.setShow(true);
		}
	}, [hintBundle.done]);

	return (
        <div className="w-full h-full flex items-center justify-center">
			<Tip bundle={hintBundle} />
			<Tip bundle={clickBundle} />
            <div className={`w-[90%] h-[90%] border-white border-4 border`}></div>
        </div>
    )
}
