import { useState, useEffect, useRef } from "react";
import { Button, Input } from "antd";
import { useRouter } from "next/navigation";

import { useGameState } from '../../components/GameState';

export default function Joiner() {
	const router = useRouter();

	const { setGameCode, playerName, setPlayerName, client } = useGameState();

	const inputRef = useRef(null);

	useEffect(() => {
		// Focus the input when component mounts
		inputRef.current?.focus();
	}, []);

	function joinCode(code) {
		console.log('joining', playerName, code);
		client.JoinGame(code, playerName).then((resp) => {
			console.log('joined game', resp);
			setGameCode(resp.code);
			router.push(`/gameboard/${resp.code}`);
		}).catch((error) => {
			console.error('join error', error);
			setGameCode('');
			//window.location.reload();
		});
	}

	function joinPickup() {
		joinCode("<>");
	}

	return <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-fit min-w-[20rem] space-y-8">
		<Input
			ref={inputRef}
			placeholder="enter your designation"
			size="large"
			className="text-lg"
			value={playerName}
			onChange={(e) => setPlayerName(e.target.value)}
			onPressEnter={joinPickup}
		/>
		<div className="flex gap-2 text-white">
			<span className="text-3xl text-white font-bold">#</span>
			<Input.OTP
				placeholder="enter code or leave blank"
				size="large"
				className="text-lg"
				formatter={(str) => str.toUpperCase()}
				disabled={playerName === ""}
				onChange={joinCode}
			/>
		</div>
		<div className="flex justify-end gap-2 items-center">
			<Button
				type="primary"
				size="large"
				disabled={playerName === ""}
				onClick={joinPickup}
			>
				pick-up game
			</Button>
		</div>
	</div>;
}