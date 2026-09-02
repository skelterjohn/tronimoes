import { useEffect } from "react";
import Image from "next/image";

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
		{/* eslint-disable-next-line @next/next/no-img-element -- remote Klipy CDN URL, not covered by next/image's remotePatterns config */}
		<img src={url} alt="react" className="pointer-events-auto origin-top translate-y-[30px]"/>
		<Image
			src="/klipy_watermark.png"
			alt="KLIPY"
			width={390}
			height={134}
			className="absolute bottom-0 left-0 right-4 w-20 h-auto translate-y-[30px] p-2 pointer-events-none"
		/>
	</div>;
}

