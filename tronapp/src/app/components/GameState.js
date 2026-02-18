"use client";

import { createContext, useContext, useState, useEffect } from 'react';
import clientFor from '../../client/Client';
import { auth } from "@/config";
import { useAuthState } from 'react-firebase-hooks/auth';

export const GameContext = createContext();

const defaultConfig = {
	tileset: "classic",
	soundEffects: true,
};

export function GameProvider({ children }) {
	const [gameCode, setGameCode] = useState("");
	const [playerName, setPlayerName] = useState("");
	const [client, setClient] = useState(undefined);
	const [persistentUser, loading, error] = useAuthState(auth);
	const [userInfo, setUserInfo] = useState(null);
	const [config, setConfig] = useState(defaultConfig);
	
	useEffect(() => {
		if (!client?.userInfo) {
			return;
		}
		client?.GetPlayer().then((resp) => {
			setPlayerName(resp.name);
			console.log("got player config", resp);
			setConfig(resp.config);
		}).catch((error) => {
			console.error('get player name error', error);
		});
	}, [client, setPlayerName]);

	useEffect(() => {
		if (!client?.userInfo) {
			return;
		}
		if (!config) {
			return;
		}
		console.log("updating player config", config);
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
			config, setConfig,
		}}>
			{children}
		</GameContext.Provider>
	);
}

export function useGameState() {
	return useContext(GameContext);
}