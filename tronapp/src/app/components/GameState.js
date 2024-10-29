"use client";

import { createContext, useContext, useState } from 'react';

const GameContext = createContext();

export function GameProvider({ children }) {
    const [gameCode, setGameCode] = useState("xyz123");

    return (
        <GameContext.Provider value={{ gameCode, setGameCode }}>
            {children}
        </GameContext.Provider>
    );
}

export function useGameCode() {
    return useContext(GameContext);
}