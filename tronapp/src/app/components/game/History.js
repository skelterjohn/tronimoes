import { useEffect, useRef } from 'react';

export default function History({ history }) {
	const scrollRef = useRef();

	useEffect(() => {
		scrollRef.current?.scrollTo(0, scrollRef.current.scrollHeight);
	}, [history]);

	return (
		<div 
			ref={scrollRef}
			className="pl-2 space-y-1 h-full border-green-700 text-green-700 border overflow-y-scroll"
		>
			<ol className="list-decimal list-inside">
				{history.map((h, i) => (
					<li 
						key={i} 
						className="whitespace-normal break-words"
					>
						{h}
					</li>
				))}
			</ol>
		</div>
	);
}
