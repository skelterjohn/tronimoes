export default function History({ history }) {
	return (
		<div className="space-y-1 h-[40rem] overflow-y-auto">
			<ol className="list-decimal list-inside">
				{history.map((h, i) => (
					<li 
						key={i} 
						className="whitespace-nowrap overflow-hidden text-ellipsis"
					>
						{h}
					</li>
				))}
			</ol>
		</div>
	);
}
