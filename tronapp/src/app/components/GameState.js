"use client";

import { createContext, useContext, useState, useEffect } from 'react';
import clientFor from '../../client/Client';
import { auth } from "@/config";
import { useAuthState } from 'react-firebase-hooks/auth';

const GameContext = createContext();

export function GameProvider({ children }) {
	const [gameCode, setGameCode] = useState("");
	const [playerName, setPlayerName] = useState("");
	const [client, setClient] = useState(undefined);
	const [persistentUser, loading, error] = useAuthState(auth);
	const [userInfo, setUserInfo] = useState(null);
	const [tutorial, setTutorial] = useState(false);
	const [config, setConfig] = useState(null);
	
	useEffect(() => {
		if (!client?.userInfo) {
			return;
		}
		client?.GetPlayer().then((resp) => {
			setPlayerName(resp.name);
			setConfig(resp.config);
		}).catch((error) => {
			console.error('get player name error', error);
		});
	}, [client, setPlayerName]);

	useEffect(() => {
		if (!client?.userInfo) {
			return;
		}
		client?.UpdatePlayerConfig(config).catch((error) => {
			console.error('update player error', error);
		});
	}, [config, playerName]);

	useEffect(()=> {
		if (error !== undefined) {
			setErrorMessage(error.message);
			setUserInfo(undefined);
			return;
		}
		if (!loading) {
			setUserInfo(persistentUser);
		}
	}, [persistentUser, loading, error]);

	useEffect(() => {
		setClient(clientFor(playerName, userInfo));
	}, [playerName, userInfo]);

	return (
		<GameContext.Provider value={{
			gameCode, setGameCode,
			playerName, setPlayerName,
			client, setClient,
			userInfo, setUserInfo,
			persistentUser, loading, error,
			tutorial, setTutorial,
			config, setConfig,
		}}>
			{children}
		</GameContext.Provider>
	);
}

export function useGameState() {
	return useContext(GameContext);
}