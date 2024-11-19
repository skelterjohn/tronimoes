import { useEffect, useRef } from 'react';

export default function History({ history }) {
	const scrollRef = useRef();

	useEffect(() => {
		scrollRef.current?.scrollTo(0, scrollRef.current.scrollHeight);
	}, [history]);

	return (
		<div 
			ref={scrollRef}
			className="space-y-1 h-full overflow-y-scroll"
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
