export default function React({
	show, setShow, url
}) {
	if (!show) {
		return null;
	}
	return <div
		onClick={() => setShow(false)}
		className="absolute z-50 pointer-events-none"
	>
		<img src={url} alt="react" className="pointer-events-auto origin-top translate-y-[30px]"/>
	</div>;
}

