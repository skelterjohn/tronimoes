"use client";

import { createContext, useContext, useState, useEffect } from 'react';
import clientFor from '../../client/Client';
import { auth } from "@/config";
import { useAuthState } from 'react-firebase-hooks/auth';

const GameContext = createContext();

export function GameProvider({ children }) {
	const [gameCode, setGameCode] = useState("");
	const [playerName, setPlayerName] = useState("");
	const [playerID, setPlayerID] = useState("");
	const [client, setClient] = useState(undefined);
	const [persistentUser, loading, error] = useAuthState(auth);
	const [userInfo, setUserInfo] = useState(null);
	const [tutorial, setTutorial] = useState(false);
	
	useEffect(() => {
		if (!client?.userInfo || !client?.userid) {
			return;
		}
		client?.GetPlayerName().then((resp) => {
			setPlayerName(resp.name);
			// setIsRegistered(true);
		}).catch((error) => {
			console.error('get player name error', error);
		});
	}, [client, setPlayerName]);

	useEffect(() => {
		setPlayerID(userInfo?.uid);
	}, [userInfo, setPlayerID]);

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
		setClient(clientFor(playerName, playerID, userInfo));
	}, [playerName, playerID, userInfo]);

	return (
		<GameContext.Provider value={{
			gameCode, setGameCode,
			playerName, setPlayerName,
			playerID, setPlayerID,
			client, setClient,
			userInfo, setUserInfo,
			persistentUser, loading, error,
			tutorial, setTutorial,
		}}>
			{children}
		</GameContext.Provider>
	);
}

export function useGameState() {
	return useContext(GameContext);
}