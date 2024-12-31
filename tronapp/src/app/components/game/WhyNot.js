import { useCallback } from "react";

export default function WhyNot({message, setMessage, inFlight}) {
	const clickMessage = useCallback(() => {
		setMessage("");
	}, [setMessage]);
	
	if (!message && !inFlight) return null;

	if (message) {
		return (
			<div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 
							p-3 px-5 bg-white border-2 border-red-500 rounded-md text-red-500 
							z-50" onClick={clickMessage}>
				{message}
			</div>
		)
	}
	if (inFlight) {
		return (
			<div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 
							p-3 px-5 bg-white border-2 border-blue-500 rounded-md text-blue-500 
							z-50" onClick={clickMessage}>
				{inFlight}
			</div>
		)
	}
	return null;
}
