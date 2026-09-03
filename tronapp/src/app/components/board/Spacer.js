"use client";

function rotationFor(spacer) {
	if (!spacer) {
		return "";
	}
	if (spacer.b.x > spacer.a.x) {
		return "-rotate-90";
	}
	if (spacer.b.x < spacer.a.x) {
		return "rotate-90";
	}
	if (spacer.b.y < spacer.a.y) {
		return "rotate-180";
	}
	return "";
}

export default function Spacer({ spacer }) {
	const rotate = rotationFor(spacer);

	return <div className={`absolute w-full h-full ${rotate}`}>
		<div className="h-[600%] bg-white border-black border-2 rounded-lg flex items-center justify-center">
			<div className="rotate-90 whitespace-nowrap text-black absolute transform origin-center">
				RIGHT-CLICK TO CLEAR
			</div>
		</div>
	</div>
}
