export default function Hint({ spacerOffset }) {
	const dx = spacerOffset?.x ?? 0;
	const dy = spacerOffset?.y ?? 0;
	const hasOffset = spacerOffset !== undefined;
	const lineLength = Math.hypot(dx, dy) * 100;
	const lineAngle = Math.atan2(dy, dx) * 180 / Math.PI;

	return (
		<div className="w-full h-full flex items-center justify-center relative">
			<div className="w-[90%] h-[90%] border-white border-4 border"></div>
			{hasOffset && (
				<div
					className="absolute left-1/2 top-1/2 h-[4px] bg-white pointer-events-none"
					style={{
						width: `${lineLength}%`,
						transform: `translateY(-50%) rotate(${lineAngle}deg)`,
						transformOrigin: "0 50%",
					}}
				/>
			)}
		</div>
	)
}
