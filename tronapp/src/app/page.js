"use client";

import { useState } from 'react';
import { useRouter } from 'next/navigation'; 
import Image from 'next/image';
import { Input, Button, Row, Col } from 'antd';

import { useGameCode } from './components/GameState';

export default function Home() {
	const router = useRouter();

	const [name, setName] = useState("");
	const [code, setCode] = useState("");

	const { setGameCode, setPlayerName } = useGameCode();

	function joinCode() {
		setGameCode(code);
		setPlayerName(name);
		router.push('/gameboard');
	}

	function joinRandom() {
		setGameCode("abc123");
		setPlayerName(name);
		router.push('/gameboard');
	}

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
				<Input
					placeholder="name"
					size="large"
					className="text-lg"
					value={name}
					onChange={(e) => setName(e.target.value)}
				/>
				<div className="flex gap-2">
					<Input
						placeholder="code"
						size="large"
						className="text-lg"
						value={code}
						onChange={(e) => setCode(e.target.value)}
					/>
					<Button
						type="primary"
						size="large"
						onClick={joinCode}
					>
						join
					</Button>
				</div>
				<Button
					className="w-full"
					type="primary"
					size="large"
					onClick={joinRandom}
				>
					random
				</Button>
			</div>
		</main>
	);
}
