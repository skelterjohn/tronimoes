import { useEffect, useRef } from 'react';

export default function History({ history, cursor = false }) {
	const scrollRef = useRef();

	useEffect(() => {
		scrollRef.current?.scrollTo(0, scrollRef.current.scrollHeight);
	}, [history]);

	return (
		<div 
			ref={scrollRef}
			className="font-game pl-2 space-y-1 h-full border-green-700 text-green-700 border overflow-y-scroll"
		>
			<ol className="list-decimal list-inside">
				{history.map((h, i) => (
					<li 
						key={i} 
						className="whitespace-normal wrap-break-word"
					>
						{h}
					</li>
				))}
				{cursor && <li>
					<span
						aria-hidden
						className="inline-block ml-0.5"
						style={{ animation: 'history-cursor-blink 1s step-end infinite' }}
					>
						█
					</span>
				</li>}
			</ol>
		</div>
	);
}
