"use client";

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useGameCode } from '../GameState';
import Board from '../board/Board';
import Hand from './Hand';
import client from '../../../client/Client';

const availableColors = [
	"red",
	"blue",
	"green",
	"indigo",
	"orange",
	"fuchsia",
]

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
	const router = useRouter();
	const { gameCode, playerName } = useGameCode();

	// These states come from the server
	const [version, setVersion] = useState(0);
	useEffect(() => {
		console.log("version", version);
	}, [version]);
	const [players, setPlayers] = useState([]);
	useEffect(() => {
		console.log("players", players);
	}, [players]);	

	const [turnIndex, setTurnIndex] = useState(0);
	useEffect(() => {
		console.log("turnIndex", turnIndex);
	}, [turnIndex]);
	const [laidTiles, setLaidTiles] = useState({});
	useEffect(() => {
		console.log("laidTiles", laidTiles);
	}, [laidTiles]);

	// here we query the server
	const [game, setGame] = useState(undefined);
	useEffect(() => {
		if (gameCode === "") {
			setGame(undefined);
			router.push('/');
		}
		console.log("getting game", gameCode, version);
		client.GetGame(gameCode, version).then((resp) => {
			console.log("game", resp);
			setVersion(resp.version);
			setGame(resp);
		}).catch((error) => {
			console.error("error", error);
		});
	}, [gameCode, version]);

	useEffect(() => {
		if (game === undefined) {
			return;
		}
		setVersion(game.version);

		let playerColors = {}

		setPlayers(game.players.map((p, i) => {
			console.log(i, p)
			playerColors[p.name] = availableColors[i];
			return {
				name: p.name,
				color: availableColors[i],
				tiles: p.hand?.map((t) => ({
					a: t.pips_a, 
					b: t.pips_b,
				})),
				dead: false,
			}
		}));

		let allLaidTiles = {}
		if (game.rounds?.length > 0) {
			const lastRound = game.rounds[game.rounds.length-1]
			lastRound?.laid_tiles?.forEach((lt) => {
				allLaidTiles[`${lt.x},${lt.y}`] = {
					a: lt.tile.pips_a,
					b: lt.tile.pips_b,
					orientation: lt.orientation,
					color: playerColors[lt.player_name],
					dead: false,
				}
			});
			setTurnIndex(lastRound?.turn);
		}
		setLaidTiles(allLaidTiles);

	}, [game]);


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

	const playerTurn = players[turnIndex];

	let borderColor = "bg-white";
	let myTurn = false;
	if (playerTurn !== undefined) {
		borderColor = borderColorMap[playerTurn.color];
		myTurn = players.length > 0 && playerTurn.name === playerName;
	}

	return <div className="">
		<div className="flex justify-center items-center">
			<div className="text-center text-5xl font-bold">#{gameCode}</div>
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