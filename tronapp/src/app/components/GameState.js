"use client";

import { createContext, useContext, useState, useEffect } from 'react';
import { v4 as uuidv4 } from 'uuid';
import clientFor from '../../client/Client';

const GameContext = createContext();

export function GameProvider({ children }) {
    const [gameCode, setGameCode] = useState("");
	const [playerName, setPlayerName] = useState(undefined);
	const [playerKey, setPlayerKey] = useState(uuidv4());
	const [client, setClient] = useState(undefined);


	useEffect(() => {
		console.log("playerName", playerName);
		setClient(clientFor(playerName, playerKey));
	}, [playerName, playerKey]);

    return (
        <GameContext.Provider value={{
			gameCode, setGameCode,
			playerName, setPlayerName,
			playerKey, setPlayerKey,
			client, setClient,
		}}>
            {children}
        </GameContext.Provider>
    );
}

export function useGameState() {
    return useContext(GameContext);
}