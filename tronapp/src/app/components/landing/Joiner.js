import { useState, useEffect, useRef } from "react";
import { Button, Input } from "antd";
import { useRouter } from "next/navigation";

import { useGameState } from '../../components/GameState';

export default function Joiner({userInfo, loading, setErrorMessage}) {
	const router = useRouter();

	const { setGameCode, playerName, setPlayerName, setPlayerKey, setPlayerID, client } = useGameState();
	const [isRegistered, setIsRegistered] = useState(false);

	const [nameInput, setNameInput] = useState('');

	const inputRef = useRef(null);

	useEffect(() => {
		if (!isRegistered) {
			// Add a small delay to ensure DOM is ready and other focus events have completed
			const timer = setTimeout(() => {
				if (inputRef.current && document.activeElement !== inputRef.current) {
					inputRef.current.focus();
				}
			}, 100);
			
			return () => clearTimeout(timer);
		}
	}, [isRegistered]);

	useEffect(() => {
		if (userInfo === undefined || userInfo === null) {
			setIsRegistered(false);
			setPlayerName('');
		}
	}, [userInfo, setPlayerKey, setPlayerName]);

	useEffect(() => {
		setIsRegistered(playerName !== '')
		setNameInput(playerName);
	}, [playerName]);

	function registerAndJoinCode(code) {
		if (!isRegistered) {
			setPlayerName(nameInput);
			console.log('registering', nameInput);
			client.RegisterPlayerName(nameInput).then(() => {
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
		console.log('joining', nameInput, code);
		setPlayerName(nameInput);
		client.JoinGame(code, nameInput).then((resp) => {
			console.log('joined game', resp);
			setGameCode(resp.code);
			router.push(`/gameboard/${resp.code}`);
		}).catch((error) => {
			console.error('join error', error);
			setPlayerName('');
			setNameInput('');
			setIsRegistered(false);
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
			<div>
				<div className="font-[Roboto_Mono] text-white">designation</div>
				<div className="bg-black rounded-lg h-[4rem] p-4 border-white border text-white">
					<p className="font-[Roboto_Mono] text-xl tracking-wider">{playerName}</p>
				</div>
			</div>
		) : (
			<Input
				ref={inputRef}
				placeholder="enter your designation"
				size="large"
				className="text-lg"
				value={nameInput}
				disabled={isRegistered}
				onChange={(e) => setNameInput(e.target.value)}
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
				disabled={nameInput === ""}
				onChange={registerAndJoinCode}
			/>
		</div>
		<div className="flex justify-end gap-2 items-center">
			<Button
				type="primary"
				size="large"
				disabled={nameInput === ""}
				onClick={joinPickup}
			>
				pick-up game
			</Button>
		</div>
	</div>;
}