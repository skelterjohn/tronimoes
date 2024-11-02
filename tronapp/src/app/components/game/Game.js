"use client";

import { useState, useEffect, useCallback } from 'react';
import { useGameCode } from '../GameState';
import Board from '../board/Board';
import Hand from './Hand';
import client from '../../../client/Client';

const borderColorMap = {
	red: "border-red-300",
	blue: "border-blue-300",
	green: "border-green-300",
	indigo: "border-indigo-300",
	orange: "border-orange-300",
	fuchsia: "border-fuchsia-300",
	white: "border-white"
};

function Game() {
	const { gameCode, playerName } = useGameCode();

	// These states come from the server
	const [version, setVersion] = useState(0);
	const [players, setPlayers] = useState([]);

	useEffect(() => {
		setPlayers([
			{
				name: "Cool Symbiote",
				color: "red",
				tiles: [{},{},{},{}],
				dead: true,
			},
			{
				name: "Hot Xenophage",
				color: "blue",
				tiles: [{},{},{}],
				dead: false,
			},
			{
				name: playerName,
				color: "green",
				tiles: [{a:1, b:2}, {a:3, b:12}],
				dead: false,
			},
		])
	}, [playerName]);

	const [turnIndex, setTurnIndex] = useState(2);

	const [laidTiles, setLaidTiles] = useState({
		"4,5": {a:12, b:12, orientation:"right", color:"white", dead:false},

		"6,5": {a:12, b:3, orientation:"down", color:"red", dead:true},

		"5,6": {a:12, b:8, orientation:"down", color:"green", dead:false},
		"6,7": {a:8, b:10, orientation:"right", color:"green", dead:false},
		"7,6": {a:10, b:2, orientation:"up", color:"green", dead:false},
		
		"3,5": {a:12, b:13, orientation:"left", color:"blue", dead:false},
		"2,4": {a:13, b:7, orientation:"right", color:"blue", dead:false},
		"4,4": {a:7, b:3, orientation:"right", color:"blue", dead:false},
		"6,4": {a:3, b:15, orientation:"right", color:"blue", dead:false},
	});

	// The remaining states are derived.

	const [playerColor, setPlayerColor] = useState("green");

	const [opponents, setOpponents] = useState([]);
	useEffect(() => {
		var playerIndex = players.findIndex(p => p.name === playerName);
		if (playerIndex === -1) {
			return;
		}
		var oppList = [];
		for (let offset=1; offset<players.length; offset++) {
			const opp = players[(playerIndex+offset)%players.length];
			oppList.push(opp);
		}
		setOpponents(oppList);
		setPlayerHand(players[playerIndex].tiles);
		setPlayerColor(players[playerIndex].color);
	}, [players, playerName]);

	const [playerHand, setPlayerHand] = useState([]);
	const [selectedTile, setSelectedTile] = useState(undefined);

	// here we query the server
	useEffect(() => {
		console.log("getting game", gameCode, version);
		client.GetGame(gameCode, version).then((resp) => {
			console.log("game", resp);
			setVersion(resp.version);
		}).catch((error) => {
			console.error("error", error);
		});
	}, [gameCode, version]);


	function startTurn() {
		setTurnIndex(2);
		setSelectedTile(undefined);
	}

	function playTile(tile) {
		tile.color = playerColor;
		setLaidTiles({...laidTiles, [`${tile.x},${tile.y}`]: tile});
		startTurn();
	}

	useEffect(() => {
		startTurn();
	}, []);

	let borderColor = "bg-white";
	if (players.length > 0) {
		borderColor = borderColorMap[players[turnIndex].color];
	}
	let myTurn = players.length > 0 && players[turnIndex].name === playerName;

	return <div className="">
		<div className="flex justify-center items-center">
			<div className="text-center text-5xl font-bold">{gameCode}</div>
		</div>
		<div className="flex justify-center items-center gap-4">
			{opponents.map((o, i) => (
				<div key={i} className="flex-1">
					<Hand
						name={o.name}
						color={o.color}
						hidden={true}
						dead={o.dead}
						tiles={o.tiles}
						/>
				</div>
			))}
		</div>
		<div>
			<div className={`${borderColor} border-8`}>
				<Board
					width={10} height={11}
						tiles={laidTiles}
						selectedTile={selectedTile}
						playTile={playTile}
				/>
			</div>
		</div>
		<div>
			<Hand
				name={playerName}
				color={playerColor}
				tiles={playerHand}
				selectedTile={selectedTile}
				setSelectedTile={setSelectedTile}
				playerTurn={myTurn}
			/>
		</div>
	</div>;
}

export default Game;