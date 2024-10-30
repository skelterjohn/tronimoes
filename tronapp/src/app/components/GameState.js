"use client";

import { createContext, useContext, useState } from 'react';

const GameContext = createContext();

export function GameProvider({ children }) {
    const [gameCode, setGameCode] = useState("xyz123");
	const [playerName, setPlayerName] = useState("Rad Bicycle");

    return (
        <GameContext.Provider value={{ gameCode, setGameCode, playerName, setPlayerName }}>
            {children}
        </GameContext.Provider>
    );
}

export function useGameCode() {
    return useContext(GameContext);
}