import Image from 'next/image';
import { useGameState } from '@/app/components/GameState';

export default function Pips({pips}) {

	const { config } = useGameState();

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
			<Image 
				src={`/tilesets/${config.tileset || 'beehive'}/${pips}.svg`}
				width={0}
				height={0}
				sizes="100%"
				className="w-[100%] h-[100%]"
				alt="tile"
			/>
		</div>
	)
}
