import { useState } from "react";
import { Button, Input } from "antd";

import { useGameCode } from '../../components/GameState';

export default function Joiner() {
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

	return <div>
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
	</div>;
}