import { useEffect } from "react";

export default function React({
	show, setShow, url
}) {
	if (!show) {
		return null;
	}

	useEffect(() => {
		if (!show) {
			return;
		}
		setTimeout(() => {
			setShow(false);
		}, 10000);
	}, [show]);

	return <div
		onClick={() => setShow(false)}
		className="absolute z-50 pointer-events-none"
	>
		<img src={url} alt="react" className="pointer-events-auto origin-top translate-y-[30px]"/>
	</div>;
}

