import { useState, useEffect, useRef } from "react";
import { Button, Input } from "antd";
import { useRouter } from "next/navigation";

import { useGameState } from '../../components/GameState';

export default function Joiner({userInfo, loading, setErrorMessage}) {
	const router = useRouter();

	const { setGameCode, playerName, setPlayerName, setPlayerKey, setPlayerID, client } = useGameState();
	const [isRegistered, setIsRegistered] = useState(false);

	const inputRef = useRef(null);

	useEffect(() => {
		// Focus the input when component mounts
		inputRef.current?.focus();
	}, []);

	useEffect(() => {
		setPlayerKey(userInfo?.accessToken);
		setPlayerID(userInfo?.uid);
		if (userInfo === undefined || userInfo === null) {
			setIsRegistered(false);
			setPlayerName('');
		}
	}, [userInfo, setPlayerKey, setPlayerID, setPlayerName]);

	useEffect(() => {
		if (!client?.key || !client?.userid) {
			return;
		}
		client?.GetPlayerName().then((resp) => {
			setPlayerName(resp.name);
			setIsRegistered(true);
		}).catch((error) => {
			console.error('get player name error', error);
		});
	}, [client, setPlayerName]);

	function registerAndJoinCode(code) {
		if (!isRegistered) {
			console.log('registering', playerName);
			client.RegisterPlayerName(playerName).then(() => {
				setIsRegistered(true);
				joinCode(code);
			}).catch((error) => {
				console.error('register error', error);
				setErrorMessage(error.data.error);
			});
		} else {
			joinCode(code);
		}
	}
	function joinCode(code) {
		console.log('joining', playerName, code);
		client.JoinGame(code, playerName).then((resp) => {
			console.log('joined game', resp);
			setGameCode(resp.code);
			router.push(`/gameboard/${resp.code}`);
		}).catch((error) => {
			console.error('join error', error);
			setPlayerName('');
			setErrorMessage(error.data.error);
		});
	}

	function joinPickup() {
		registerAndJoinCode("PICKUP");
	}

	if (loading) {
		return (
			<div className="bg-black rounded-lg p-4 border-white border text-white">
				<p className="font-[Roboto_Mono] text-xl tracking-wider">checking connection...</p>
			</div>
		);
	}
	return <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-fit min-w-[20rem] space-y-8">
		{isRegistered ? (
			<div className="bg-black rounded-lg p-4 border-white border text-white">
				<p className="font-[Roboto_Mono] text-xl tracking-wider">designation: {playerName}</p>
			</div>
		) : (
			<Input
				ref={inputRef}
				placeholder="enter your designation"
				size="large"
				className="text-lg"
				value={playerName}
				disabled={isRegistered}
				onChange={(e) => setPlayerName(e.target.value)}
				onPressEnter={joinPickup}
			/>
		)}
		<div className="flex gap-2 text-white">
			<span className="text-3xl text-white font-bold">#</span>
			<Input.OTP
				placeholder="enter code or leave blank"
				size="large"
				className="text-lg"
				formatter={(str) => str.toUpperCase()}
				disabled={playerName === ""}
				onChange={registerAndJoinCode}
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