import Image from 'next/image';


export default function Pips({pips}) {
	if (pips === -1) {
		return (
			<div className="w-full h-full flex justify-center items-center flex-col border border-black rounded-lg bg-white">
				<div>FREE LINE</div>
				<div>SPACER</div>
				<div>START</div>
			</div>
		);
	}
	return (
		<div className="w-full flex justify-center items-center">
			<Image 
				src={`/tilesets/beehive/${pips}.png`}
				width={0}
				height={0}
				sizes="100%"
				className="w-[100%] h-[100%]"
				alt="tile"
			/>
		</div>
	)
}
