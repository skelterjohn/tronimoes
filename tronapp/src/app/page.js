"use client";

import Image from 'next/image';

import Joiner from './components/landing/Joiner';

export default function Home() {
	return (
		<main className="relative min-h-screen w-screen bg-slate-800">
			<Image 
				src="/trondude.png"
				alt="Background"
				fill
				className="object-cover z-0"
				priority
			/>
			<div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-fit space-y-4">
				<Joiner/>
			</div>
		</main>
	);
}
