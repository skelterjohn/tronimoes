import { useEffect } from "react";

export default function Reaction({
	setShow, url
}) {
	useEffect(() => {
		setTimeout(() => {
			setShow(false);
		}, 10000);
	}, [setShow]);

	return <div
		onClick={() => setShow(false)}
		className="absolute z-50 pointer-events-none"
	>
		<img src={url} alt="react" className="pointer-events-auto origin-top translate-y-[30px]"/>
		<img 
	        src="/klipy_watermark.png" 
	        alt="KLIPY" 
	        className="absolute bottom-0 left-0 right-4 w-20 h-auto translate-y-[30px] p-2 pointer-events-none"
	    />
	</div>;
}

