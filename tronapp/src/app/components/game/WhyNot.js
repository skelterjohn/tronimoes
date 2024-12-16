import { useCallback } from "react";

export default function WhyNot({message, setMessage}) {
	if (!message) return null;

	const clickMessage = useCallback(() => {
		setMessage("");
	}, [setMessage]);

	return (
		<div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 
						p-3 px-5 bg-white border-2 border-red-500 rounded-md text-red-500 
						z-50" onClick={clickMessage}>
			{message}
		</div>
	)
}
