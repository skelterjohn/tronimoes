import Image from 'next/image';
import { useState, useEffect } from 'react';
import { useGameState } from '@/app/components/GameState';

export default function Pips({pips, parentRotation = 0}) {

	const { config } = useGameState();
	const [isNumbers, setIsNumbers] = useState(false);

	useEffect(() => {
		setIsNumbers(config?.tileset === 'numbers');
	}, [config]);

	if (pips === -1) {
		return (
			<div className="w-full h-full flex justify-center items-center flex-col border border-black rounded-lg bg-white">
				<div>FREE LINE</div>
				<div>SPACER</div>
				<div>START</div>
			</div>
		);
	}
	if (pips === undefined) {
		return (
			<div className="w-full h-full flex justify-center items-center">
				<div className="w-[90%] h-[90%] flex flex-col justify-center items-center border border-black rounded-lg bg-gray-400">
					<div>choose</div>
					<div>tile</div>
				</div>
			</div>
		);
	}

	return (
		<div className="w-full flex justify-center items-center">
			{isNumbers ? (
				<div style={{ transform: `rotate(${-parentRotation}deg)`, transformOrigin: 'center' }}>
					<Image 
						src={`/tilesets/numbers/${pips}.svg`}
						width={0}
						height={0}
						sizes="100%"
						className="w-full h-full"
						alt="tile"
					/>
				</div>
			) : (
				<Image 
					src={`/tilesets/${config?.tileset || 'classic-mono'}/${pips}.svg`}
					width={0}
					height={0}
					sizes="100%"
					className="w-full h-full"
					alt="tile"
				/>
			)}
		</div>
	)
}
