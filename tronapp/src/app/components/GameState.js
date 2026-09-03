"use client";

import { createContext, useContext, useState, useEffect, useMemo } from 'react';
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
	const [persistentUser, loading, error] = useAuthState(auth);
	const [userInfo, setUserInfo] = useState(null);
	const [config, setConfig] = useState(defaultConfig);
	const [options, setOptions] = useState(defaultOptions);

	// userInfo is also set directly by sign-in/sign-out flows (SignIn.js, page.js),
	// so this is a sync-then-diverge pattern: adjust it during render (rather than
	// in an effect) when the underlying auth state actually changes.
	const [prevAuthState, setPrevAuthState] = useState({ persistentUser, loading, error });
	if (persistentUser !== prevAuthState.persistentUser || loading !== prevAuthState.loading || error !== prevAuthState.error) {
		setPrevAuthState({ persistentUser, loading, error });
		if (error !== undefined) {
			setErrorMessage(error.message);
			setUserInfo(undefined);
		} else if (!loading) {
			setUserInfo(persistentUser);
		}
	}

	// client is a pure derivation of playerName/userInfo. When userInfo exists, the
	// server uses the Firebase token + X-Player-Id rather than X-Player-Name, so we
	// avoid recreating the client (and re-triggering effects below) while the
	// authenticated user is typing/changing `playerName`. clientFor() reads
	// window.location, so it can only run client-side.
	const client = useMemo(() => {
		if (typeof window === 'undefined') {
			return undefined;
		}
		return userInfo ? clientFor("", userInfo) : clientFor(playerName, userInfo);
	}, [playerName, userInfo]);

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

	return (
		<GameContext.Provider value={{
			gameCode, setGameCode,
			playerName, setPlayerName,
			client,
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