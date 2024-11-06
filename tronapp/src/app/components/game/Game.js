"use client";

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useGameState } from '../GameState';
import Board from '../board/Board';
import Hand from './Hand';
import History from './History';
import { Button, Modal } from 'antd';
import WhyNot from './WhyNot';

const availableColors = [
	"red",
	"blue",
	"green",
	"indigo",
	"orange",
	"fuchsia",
]

function Game({code}) {
	const router = useRouter();
	const { playerName, client } = useGameState();

	// These states come from the server
	const [version, setVersion] = useState(-1);
	const [players, setPlayers] = useState([]);
	const [boardWidth, setBoardWidth] = useState(1);
	const [boardHeight, setBoardHeight] = useState(1);

	const [turnIndex, setTurnIndex] = useState(0);
	const [laidTiles, setLaidTiles] = useState({});
	const [lineHeads, setLineHeads] = useState({})

	const [roundHistory, setRoundHistory] = useState([]);
	const [gameHistory, setGameHistory] = useState([]);

	// here we query the server
	const [game, setGame] = useState(undefined);
	useEffect(() => {
		if (code === "") {
			setGame(undefined);
			router.push('/');
		}

		// Add cleanup flag
		let isActive = true;

		const getGame = () => {
			const myCode = code;
			if (client === undefined) {
				isActive = false;
				return;
			}
			client.GetGame(code, version).then((resp) => {
				// Only update state if component is still mounted
				if (!isActive) return;

				if (resp.version === version) {
					// We got back the same one, so let's try again after a bit.
					setTimeout(getGame, 5000);
				}
				setVersion(resp.version);
				setGame(resp);
			}).catch((error) => {
				if (error?.status === 404) {
					isActive = false;
					setGame(undefined);
					router.push('/');
					return;
				}
				if (!isActive) return;
				if (myCode !== code) {
					isActive = false;
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
			});
		};

		getGame();

		// Cleanup function
		return () => {
			isActive = false;
		};
	}, [code, version]);

	useEffect(()=> {
		if (!playerName) {
			// this is definitely lower than the version on the server,
			// so we get this loop started.
			setVersion(-10);
		}
	}, [playerName]);

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
				dead: p.dead,
				chickenFoot: p.chicken_foot,
				chickenFootX: p.chicken_foot_x,
				chickenFootY: p.chicken_foot_y,
				just_drew: p.just_drew,
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
			setLineHeads(Object.values(lastRound?.player_lines).map((line) => {
				return line[line.length-1];
			}))
			setRoundHistory(lastRound.history || []);
		}
		setGameHistory(game.history || []);
		setTurnIndex(game.turn);
		setLaidTiles(allLaidTiles);
		setBoardWidth(game.board_width);
		setBoardHeight(game.board_height);
	}, [game]);


	// The remaining states are derived.

	const [opponents, setOpponents] = useState([]);
	const [player, setPlayer] = useState(undefined);
	useEffect(() => {
		var playerIndex = players.findIndex(p => p.name === playerName);
		if (playerIndex === -1) {
			setOpponents(players);
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

	const [gameInProgress, setGameInProgress] = useState(false);
	useEffect(() => {
		setGameInProgress(game?.rounds !== undefined && game?.rounds?.length > 0);
	}, [game]);

	const [chickenFeet, setChickenFeet] = useState({});
	useEffect(() => {
		let allFeet = {};
		players.forEach((p) => {
			if (!p.chickenFoot) {
				return;
			}
			allFeet[`${p.chickenFootX},${p.chickenFootY}`] = p.color;
		});
		setChickenFeet(allFeet);
	}, [players]);

	const [indicated, setIndicated] = useState(undefined);

	function startRound() {
		client.StartRound(code).then((resp) => {
			console.log("started round", resp);
		}).catch((error) => {
			console.error("error", error);
		});
	}

	const [playErrorMessage, setPlayErrorMessage] = useState("");

	function playTile(tile) {
		tile.color = player.color;
		client.LayTile(code, {
			tile:{
				pips_a: tile.a,
				pips_b: tile.b,
			},
			x: tile.x,
			y: tile.y,
			orientation: tile.orientation,
			player_name: player.name,
			indicated: {
				pips_a: indicated !== undefined ? indicated.a : -1,
				pips_b: indicated !== undefined ? indicated.b : -1,
			}
		}).then((resp) => {
			setSelectedTile(undefined);
			setIndicated(undefined);
			console.log("laid tile", resp);
		}).catch((error) => {
			console.error("error", error);
			setPlayErrorMessage(error.data.error);
		});
	}

	function drawTile() {
		setSelectedTile(undefined);
		client.DrawTile(code).then((resp) => {
			console.log("drew tile", resp);
		}).catch((error) => {
			console.error("error", error);
		});
	}


	function passTurn() {
		setSelectedTile(undefined);
		client.Pass(code).then((resp) => {
			console.log("passed");
		}).catch((error) => {
			console.error("error", error);
		});
	}

	function leaveOrQuit() {
		client.LeaveOrQuit(code).catch((error) => {
			console.error("error", error);
			router.push("/");
		}).finally(() => {
			router.push("/");
		});
	}

	const playerTurn = players[turnIndex];

	let myTurn = false;
	if (playerTurn !== undefined) {
		myTurn = players.length > 0 && playerTurn.name === playerName;
	}

	const amFirstPlayer = players.length > 0 && players[0].name === playerName;

	return (
		<div className="h-full " onClick={()=>setPlayErrorMessage("")}>
			<div className="flex justify-between items-center mb-4">
				<span className="text-left text-5xl font-bold">
					#{code} {game?.done && "(done)"}
				</span>
				<div className="flex flex-col items-end gap-2">
					<div className="flex gap-2">
						<Button 
							type="primary"
							size="large"
							className="w-28"
							disabled={!amFirstPlayer || roundInProgress || game?.done}
							onClick={() => startRound()}
						>
							Start Round
						</Button>
						<Button 
							type="primary"
							size="large"
							className="w-28"
							onClick={() => leaveOrQuit()}
						>
							{(gameInProgress && !game?.done) && (<div>Quit</div>) || (<div>Leave</div>)}
						</Button>
					</div>
					{!amFirstPlayer && !roundInProgress && (players.length > 0) &&
						(<span>waiting for {players[0].name} to start the round...</span>)
					}
				</div>
			</div>
			
			<div className="flex justify-center items-center gap-4 h-32 max-h-32">
				{opponents.map((o, i) => (
					<div key={i} className="flex-1 overflow-x-auto">
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
			<div className="flex gap-4 justify-center h-[75vh] overflow-hidden">
				<span className="w-96 hidden landscape:block">
					<History history={gameHistory}/>
				</span>
				<div className="border-black border-8 min-h-0 min-w-0 flex-1 relative">
					<Board
						width={boardWidth} height={boardHeight}
						tiles={laidTiles}
						lineHeads={lineHeads}
						selectedTile={selectedTile}
						playTile={playTile}
						chickenFeet={chickenFeet}
						indicated={indicated}
						setIndicated={setIndicated}
						playerTurn={myTurn}
						activePlayer={roundInProgress && players[turnIndex]}
					/>
					<WhyNot message={playErrorMessage} />
				</div>
				<span className="w-96 hidden landscape:block">
					<History history={roundHistory}/>
				</span>
			</div>
			{player && 
				<div className="flex justify-center items-center gap-4 h-32 max-h-32">
					<div className="overflow-x-auto w-full">
						<Hand
							player={player}
							name={playerName}
							hidden={false}
							selectedTile={selectedTile}
							setSelectedTile={setSelectedTile}
								playerTurn={myTurn}
								drawTile={drawTile}
								passTurn={passTurn}
								roundInProgress={roundInProgress}
						/>
					</div>
				</div>
			}
		</div>
	);
}

export default Game;