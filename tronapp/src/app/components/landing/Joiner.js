import { useState, useEffect, useRef } from "react";
import { Input } from "antd";
import Button from '@/app/components/Button';
import { useRouter } from "next/navigation";

import { useGameState } from '../../components/GameState';
import Options from './Options';

export default function Joiner({userInfo, loading, setErrorMessage}) {
	const router = useRouter();

	const { setGameCode, playerName, setPlayerName, setPlayerKey, client, options } = useGameState();
	const [isRegistered, setIsRegistered] = useState(false);
	const [showOptionsModal, setShowOptionsModal] = useState(false);

	const [nameInput, setNameInput] = useState('');

	const inputRef = useRef(null);

	useEffect(() => {
		if (!isRegistered) {
			// Add a small delay to ensure DOM is ready and other focus events have completed
			const timer = setTimeout(() => {
				inputRef.current?.focus?.();
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

	useEffect(() => {
		console.log('unsetting game code');
		setGameCode(undefined);
	}, []);

	function register() {
		if (isRegistered) return Promise.resolve();
		if (!nameInput.trim()) return Promise.reject(new Error('empty name'));
		setPlayerName(nameInput);
		console.log('registering', nameInput);
		return client.RegisterPlayerName(nameInput)
			.then(() => setIsRegistered(true))
			.catch((error) => {
				console.error('register error', error);
				setPlayerName('');
				setNameInput('');
				setIsRegistered(false);
				setErrorMessage(error.data.error);
				throw error;
			});
	}


	function joinCode(code) {
		console.log('joining', nameInput, code);
		setPlayerName(nameInput);
		client.JoinGame(code, nameInput, options).then((resp) => {
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
		joinCode("PICKUP");
	}

	if (loading) {
		return (
			<div className="font-game bg-black rounded-lg h-16 min-w-[20rem] p-4 border-white border text-white">
				<p className="text-xl tracking-wider">checking connection...</p>
			</div>
		);
	}
	return <div className="font-game absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-fit min-w-[20rem] space-y-8">
		<div>
			<div className="text-white">designation</div>
			<Input
				placeholder="enter your username"
				size="large"
				className={`font-game text-lg ${isRegistered ? '!bg-black !text-white border-white [&.ant-input-readonly]:!bg-black [&.ant-input-readonly]:!text-white [&.ant-input-readonly]:!opacity-100 [&.ant-input-readonly]:cursor-default' : ''}`}
				style={isRegistered ? { backgroundColor: 'black', color: 'white' } : undefined}
				styles={isRegistered ? { input: { backgroundColor: 'black', color: 'white' } } : undefined}
				value={isRegistered ? playerName : nameInput}
				readOnly={isRegistered}
				onChange={(e) => !isRegistered && setNameInput(e.target.value)}
				onPressEnter={() => (isRegistered ? joinPickup() : register())}
			/>
		</div>
		<div className="flex gap-2 text-white">
			<span className="text-3xl text-white font-bold">#</span>
			<Input.OTP
				placeholder="enter code or leave blank"
				size="large"
				className="font-game text-lg"
				formatter={(str) => str.toUpperCase()}
				disabled={playerName === ""}
				onChange={joinCode}
			/>
		</div>
		<div className="flex justify-between gap-2 items-center">
			<Button
				size="small"
				className="game-btn"
				onClick={() => setShowOptionsModal(true)}
			>
				game options
			</Button>
			<Button
				type="primary"
				size="large"
				className="game-btn"
				disabled={playerName === ""}
				onClick={joinPickup}
			>
				pick-up game
			</Button>
		</div>
		<Options
			isOpen={showOptionsModal}
			onClose={() => setShowOptionsModal(false)}
		/>
	</div>;
}