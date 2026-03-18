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

const defaultOptions = {
	agent_round_out: 0,
};

export function GameProvider({ children }) {
	const [gameCode, setGameCode] = useState("");
	const [playerName, setPlayerName] = useState("");
	const [client, setClient] = useState(undefined);
	const [persistentUser, loading, error] = useAuthState(auth);
	const [userInfo, setUserInfo] = useState(null);
	const [config, setConfig] = useState(defaultConfig);
	const [options, setOptions] = useState(defaultOptions);
	
	useEffect(() => {
		if (!client?.userInfo) {
			return;
		}
		client?.GetPlayer().then((resp) => {
			setPlayerName(resp.name);
			console.log("got player config", resp);
			setConfig(resp.config);
		}).catch((error) => {
			// This can legitimately happen the first time a logged-in user opens the app
			// before they registered their player name in this game server.
			const serverError = error?.data?.error;
			if (error?.status === 404 && serverError === 'no registered player') {
				return;
			}
			console.error('GetPlayer failed', {
				status: error?.status,
				message: error?.message,
				data: error?.data,
			});
		});
	}, [client]);

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
	}, [config, client]);

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

	// Avoid recreating the client (and re-triggering effects) while the authenticated
	// user is typing/changing `playerName`. When `userInfo` exists, the server uses
	// the Firebase token + X-Player-Id rather than X-Player-Name.
	useEffect(() => {
		if (!userInfo) return;
		setClient(clientFor("", userInfo));
	}, [userInfo]);

	useEffect(() => {
		if (userInfo) return;
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
			options, setOptions,
		}}>
			{children}
		</GameContext.Provider>
	);
}

export function useGameState() {
	return useContext(GameContext);
}