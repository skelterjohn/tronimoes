"use client";

import { useState } from 'react';
import { useRouter } from 'next/navigation'; 
import Image from 'next/image';
import { Input, Button, Row, Col } from 'antd';

import Joiner from './components/landing/Joiner';

export default function Home() {
	const router = useRouter();


	return (
		<main className="relative min-h-screen w-screen bg-slate-800">
			<Image 
				src="/trondude.png"
				alt="Background"
				fill
				className="object-cover z-0"
				priority
			/>
			<div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-80 space-y-4">
				<Joiner/>
			</div>
		</main>
	);
}
