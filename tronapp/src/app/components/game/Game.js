"use client";

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useGameState } from '../GameState';
import Board from '../board/Board';
import Hand from './Hand';
import History from './History';
import { Button } from 'antd';

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
	const { gameCode, playerName, client } = useGameState();

	// These states come from the server
	const [version, setVersion] = useState(0);
	const [players, setPlayers] = useState([]);

	const [turnIndex, setTurnIndex] = useState(0);
	const [laidTiles, setLaidTiles] = useState({});

	const [roundHistory, setRoundHistory] = useState([]);
	const [gameHistory, setGameHistory] = useState([]);

	// here we query the server
	const [game, setGame] = useState(undefined);
	useEffect(() => {
		if (gameCode === "") {
			setGame(undefined);
			router.push('/');
		}

		const getGame = () => {
			const myCode = gameCode;
			client.GetGame(gameCode, version).then((resp) => {
				if (resp.version === version) {
					// We got back the same one, so let's try again after a bit.
					setTimeout(getGame, 5000);
				}
				setVersion(resp.version);
				// setting the version to something new triggers the next fetch.
				setGame(resp);
			}).catch((error) => {
				if (myCode !== gameCode) {
					// This thread is out of sync.
					return;
				}
				if (error.name !== "AbortError") {
					console.error("error", error);
					setTimeout(getGame, 30000);
					return;
				} 
				const timeoutDuration = new Date() - requestTime;
				if (timeoutDuration < 10000) {
					console.log(`request timed out quickly in ${timeoutDuration}ms`)
					setTimeout(getGame, 5000);
				} 
			})
		}
		getGame();

	}, [gameCode, version]);

	useEffect(() => {
		console.log('game', game);

		if (game === undefined) {
			return;
		}
		setVersion(game.version);

		let playerColors = {}

		setPlayers(game.players.map((p, i) => {
			playerColors[p.name] = availableColors[i];
			return {
				name: p.name,
				color: availableColors[i],
				hand: p.hand?.map((t) => ({
					a: t.pips_a, 
					b: t.pips_b,
				})),
				score: p.score,
				dead: false,
				chickenFoot: p.chicken_foot,
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
			setRoundHistory(lastRound.history || []);
		}
		setGameHistory(game.history || []);
		setTurnIndex(game.turn);
		setLaidTiles(allLaidTiles);
	}, [game]);


	// The remaining states are derived.

	const [opponents, setOpponents] = useState([]);
	const [player, setPlayer] = useState(undefined);
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
		
		if (playerIndex === -1) {
			return;
		}
		setPlayer(players[playerIndex]);
	}, [players, playerName]);

	const [selectedTile, setSelectedTile] = useState(undefined);

	const [roundInProgress, setRoundInProgress] = useState(false);
	useEffect(() => {
		const round = game?.rounds?.[game?.rounds?.length-1];
		if (round === undefined) {
			setRoundInProgress(false);
		} else {
			setRoundInProgress(!round.done);
		}
	}, [game]);

	function startRound() {
		client.StartRound(gameCode).then((resp) => {
			console.log("started round", resp);
		}).catch((error) => {
			console.error("error", error);
		});
	}

	function playTile(tile) {
		setSelectedTile(undefined);
		tile.color = player.color;
		client.LayTile(gameCode, {
			tile:{
				pips_a: tile.a,
				pips_b: tile.b,
			},
			x: tile.x,
			y: tile.y,
			orientation: tile.orientation,
			player_name: player.name,
		}).then((resp) => {
			console.log("laid tile", resp);
		}).catch((error) => {
			console.error("error", error);
		});
	}

	function drawTile() {
		setSelectedTile(undefined);
		client.DrawTile(gameCode).then((resp) => {
			console.log("drew tile", resp);
		}).catch((error) => {
			console.error("error", error);
		});
	}

	const playerTurn = players[turnIndex];

	let borderColor = "bg-white";
	let myTurn = false;
	if (playerTurn !== undefined) {
		borderColor = borderColorMap[playerTurn.color];
		myTurn = players.length > 0 && playerTurn.name === playerName;
	}

	const amFirstPlayer = players.length > 0 && players[0].name === playerName;

	return (
		<div className="h-full">
			<div className="flex justify-between items-center mb-4">
				<span className="text-left text-5xl font-bold">
					#{gameCode}
				</span>
				{amFirstPlayer && !roundInProgress &&
					<Button 
						type="primary"
						size="large"
						onClick={() => startRound()}
					>
						Start Round
					</Button>
				}
				{!amFirstPlayer && !roundInProgress && (players.length > 0) &&
					(<span>waiting for {players[0].name} to start the round...</span>)
				}
			</div>
			
			<div className="flex justify-center items-center gap-4 h-32 max-h-32">
				{opponents.map((o, i) => (
					<div key={i} className="flex-1">
						<Hand
							player={o}
							name={o.name}
							score={o.score}
							color={o.color}
							hidden={true}
							dead={o.dead}
							tiles={o.hand}
						/>
					</div>
				))}
			</div>
			<div className="flex gap-4 justify-center h-[75vh]">
				<span className="w-96 h-full">
					<History history={gameHistory}/>
				</span>
				<div className={`${borderColor} border-8 h-fit`}>
					<Board
						width={10} height={11}
						tiles={laidTiles}
						selectedTile={selectedTile}
						playTile={playTile}
					/>
				</div>
				<span className="w-96 h-full">
					<History history={roundHistory}/>
				</span>
			</div>
			<div className="flex justify-center items-center gap-4 h-32 max-h-32">
				<Hand
					player={player}
					name={playerName}
					hidden={false}
					selectedTile={selectedTile}
					setSelectedTile={setSelectedTile}
					playerTurn={myTurn}
					drawTile={drawTile}
				/>
			</div>
		</div>
	);
}

export default Game;