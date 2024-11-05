import Image from 'next/image';


export default function Pips({pips}) {
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
