import { useState, useEffect } from "react";
import { Button, Input } from "antd";
import { useRouter } from "next/navigation";

import { useGameCode } from '../../components/GameState';

import client from '../../../client/Client';

export default function Joiner() {
	const router = useRouter();

	const [name, setName] = useState("");

	const { gameCode, setGameCode, setPlayerName } = useGameCode();

	function joinCode(code) {
		console.log('joining', name, code);
		client.join(name, code).then((resp) => {
			console.log('join response', resp);
			setGameCode(code);
			setPlayerName(name);
			router.push('/gameboard');
		}).catch((error) => {
			console.error('join error', error);
			setGameCode('');
		});
		// router.push('/gameboard');
	}

	function joinPickup() {
		setGameCode(null);
		setPlayerName(name);
		router.push('/gameboard');
	}

	useEffect(() => {
		if (name === "") {
			setGameCode('');
		}
	}, [name]);

	return <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-fit min-w-[20rem] space-y-8">
		<Input
			placeholder="enter your designation"
			size="large"
			className="text-lg"
			value={name}
			onChange={(e) => setName(e.target.value)}
		/>
		<div className="flex gap-2 text-white">
			<span className="text-3xl text-white font-bold">#</span>
			<Input.OTP
				placeholder="enter code or leave blank"
				size="large"
				className="text-lg"
				formatter={(str) => str.toUpperCase()}
				disabled={name === ""}
				value={gameCode}
				onChange={joinCode}
			/>
		</div>
		<div className="flex justify-end gap-2 items-center">
			<Button
				type="primary"
				size="large"
				disabled={name === ""}
				onClick={joinPickup}
			>
				pick-up game
			</Button>
		</div>
	</div>;
}