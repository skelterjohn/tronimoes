"use client";

import { createContext, useContext, useState, useEffect } from 'react';
import { v4 as uuidv4 } from 'uuid';
import clientFor from '../../client/Client';

const GameContext = createContext();

export function GameProvider({ children }) {
	const [gameCode, setGameCode] = useState("");
	const [playerName, setPlayerName] = useState("");
	const [playerKey, setPlayerKey] = useState("");
	const [playerID, setPlayerID] = useState("");
	const [client, setClient] = useState(undefined);


	useEffect(() => {
		setClient(clientFor(playerName, playerID, playerKey));
	}, [playerName, playerKey, playerID]);

	return (
		<GameContext.Provider value={{
			gameCode, setGameCode,
			playerName, setPlayerName,
			playerKey, setPlayerKey,
			playerID, setPlayerID,
			client, setClient,
		}}>
			{children}
		</GameContext.Provider>
	);
}

export function useGameState() {
	return useContext(GameContext);
}