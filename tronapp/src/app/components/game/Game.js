"use client";

import {Row, Col} from 'react-bootstrap';
import { useState, useEffect } from 'react';

import Board from '../board/Board';
import Hand from './Hand';

const borderColorMap = {
	red: "border-red-300",
	blue: "border-blue-300",
	green: "border-green-300",
	indigo: "border-indigo-300",
	orange: "border-orange-300",
	fuchsia: "border-fuchsia-300",
	white: "border-white"
};

function Game({}) {
	// These states come from the server
	const [playerName, setPlayerName] = useState("Rad Bicycle");
	const [players, setPlayers] = useState([
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
			name: "Rad Bicycle",
			color: "green",
			tiles: [{a:1, b:2}, {a:3, b:12}],
			dead: false,
		},
	]);

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
		var oppList = [];
		for (let offset=1; offset<players.length; offset++) {
			const opp = players[(playerIndex+offset)%players.length];
			oppList.push(opp);
		}
		setOpponents(oppList);

		setPlayerHand(players[playerIndex].tiles);
		setPlayerColor(players[playerIndex].color);
	}, [players]);

	const [playerHand, setPlayerHand] = useState([]);
	const [selectedTile, setSelectedTile] = useState(undefined);

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

	return <div className="">
		<Row className="flex justify-center items-center">
			{opponents.map((o, i) => (
				<Col key={i}>
					<Hand
						name={o.name}
						color={o.color}
						hidden={true}
						dead={o.dead}
						tiles={o.tiles}
					/>
				</Col>
			))}
		</Row>
		<Row>
			<div className={`${borderColorMap[players[turnIndex].color]} border-8`}>
				<Board
					width={10} height={11}
					tiles={laidTiles}
					selectedTile={selectedTile}
					playTile={playTile}
				/>
			</div>
		</Row>
		<Row>
			<Hand
				name={playerName}
				color={playerColor}
				tiles={playerHand}
				selectedTile={selectedTile}
				setSelectedTile={setSelectedTile}
				playerTurn={players[turnIndex].name === playerName}
			/>
		</Row>
	</div>;
}

export default Game;