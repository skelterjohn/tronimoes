"use client";

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useGameState } from '../GameState';
import Board from '../board/Board';
import Hand from './Hand';
import History from './History';
import { Button, Modal } from 'antd';
import WhyNot from './WhyNot';
import VisionQuest from '../visionquest/VisionQuest';

const availableColors = [
	"red",
	"blue",
	"green",
	"indigo",
	"orange",
	"fuchsia",
]

function Game({ code }) {
	const router = useRouter();
	const { playerName, client } = useGameState();

	// These states come from the server
	const [version, setVersion] = useState(-1);
	const [players, setPlayers] = useState([]);
	const [boardWidth, setBoardWidth] = useState(1);
	const [boardHeight, setBoardHeight] = useState(1);

	const [turnIndex, setTurnIndex] = useState(0);
	const [laidTiles, setLaidTiles] = useState({});
	const [lineHeads, setLineHeads] = useState({});
	const [spacer, setSpacer] = useState(undefined);

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
			let requestTime = new Date();
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
					// setGame(undefined);
					// router.push('/');
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
	}, [code, version, client, router]);

	useEffect(() => {
		if (!playerName) {
			// this is definitely lower than the version on the server,
			// so we get this loop started.
			setVersion(-10);
		}
	}, [playerName]);

	const [bagCount, setBagCount] = useState(0);

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
				hints: p.hints,
				spacer_hints: p.spacer_hints,
				score: p.score,
				dead: p.dead,
				chickenFoot: p.chicken_foot,
				chickenFootX: p.chicken_foot_coord.x,
				chickenFootY: p.chicken_foot_coord.y,
				chickenFootURL: p.chicken_foot_url || undefined,
				just_drew: p.just_drew,
				kills: p.kills,
			}
		}));

		let allLaidTiles = {}
		if (game.rounds?.length > 0) {
			const lastRound = game.rounds[game.rounds.length - 1]
			lastRound?.laid_tiles?.forEach((lt) => {
				allLaidTiles[`${lt.x},${lt.y}`] = {
					a: lt.tile.pips_a,
					b: lt.tile.pips_b,
					orientation: lt.orientation,
					color: playerColors[lt.player_name],
					dead: lt.dead,
				}
			});
			setLineHeads(Object.values(lastRound?.player_lines).map((line) => {
				return line[line.length - 1];
			}))
			setRoundHistory(lastRound.history || []);
			setSpacer(lastRound.spacer);
		}
		setGameHistory(game.history || []);
		setTurnIndex(game.turn);
		setLaidTiles(allLaidTiles);
		setBoardWidth(game.board_width);
		setBoardHeight(game.board_height);
		setBagCount(game.bag?.length || 0);
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
		for (let offset = 1; offset < players.length; offset++) {
			const opp = players[(playerIndex + offset) % players.length];
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
		const round = game?.rounds?.[game?.rounds?.length - 1];
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
	const [chickenFeetURLs, setChickenFeetURLs] = useState({});
	useEffect(() => {
		let allFeet = {};
		let allURLs = {};
		players.forEach((p) => {
			if (!p.chickenFoot) {
				return;
			}
			allFeet[`${p.chickenFootX},${p.chickenFootY}`] = p.color;
			allURLs[`${p.chickenFootX},${p.chickenFootY}`] = p.chickenFootURL;
			console.log(p.chickenFootURL);
		});
		setChickenFeet(allFeet);
		setChickenFeetURLs(allURLs);
	}, [players]);

	const [indicated, setIndicated] = useState(undefined);

	const [hints, setHints] = useState({});
	useEffect(() => {
		if (!selectedTile) {
			setHints({});
			return;
		}
		if (selectedTile.a == -1 && selectedTile.b == -1) {
			setHints({});
			return;
		}
		if (player?.hints === null || player?.hints === undefined) {
			setHints({});
			return;
		}
		player.hand.forEach((t, i) => {
			if (t?.a !== selectedTile?.a || t?.b != selectedTile?.b) {
				return;
			}
			if (player.hints[i] === null) {
				setHints({});
				return;
			}
			let hintSet = {};
			player.hints[i].forEach((coord) => {
				hintSet[coord] = true;
			})
			setHints(hintSet);
		})
	}, [selectedTile, player]);

	const [hintedTiles, setHintedTiles] = useState([]);
	useEffect(() => {
		if (player?.hints === null || player?.hints === undefined) {
			setHintedTiles([]);
			return;
		}
		let ht = [];
		player.hand.forEach((t, i) => {
			if (player.hints[i] !== null && player.hints[i].length > 0) {
				ht.push(t);
			}
		})
		setHintedTiles(ht);
	}, [game, player]);

	function startRound() {
		client.StartRound(code).then((resp) => {
			console.log("started round", resp);
		}).catch((error) => {
			console.error("error", error);
		});
	}

	const [playErrorMessage, setPlayErrorMessage] = useState("");
	const [playA, setPlayA] = useState(undefined);

	function playTile(tile) {
		tile.color = player.color;
		client.LayTile(code, {
			tile: {
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
			},
		}).then((resp) => {
			setSelectedTile(undefined);
			setIndicated(undefined);
			setHints({});
			console.log("laid tile", resp);
		}).catch((error) => {
			console.error("error", error);
			setPlayErrorMessage(error.data.error);
		});
	}

	function playSpacer(spacer) {
		client.LaySpacer(code, spacer).then((resp) => {
			setSelectedTile(undefined);
			setIndicated(undefined);
			setHints({});
			console.log("laid spacer", resp);
		}).catch((error) => {
			console.error("error", error);
			setPlayErrorMessage(error.data.error);
		});
	}

	function clearSpacer() {
		client.LaySpacer(code, {}).then((resp) => {
			setSelectedTile(undefined);
			setIndicated(undefined);
			setHints({});
			console.log("cleared spacer", resp);
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
			setPlayErrorMessage(error.data.error);
		});
	}

	const [showVisionQuestModal, setShowVisionQuestModal] = useState(false);
	const [chickenFootURL, setChickenFootURL] = useState(undefined);

	useEffect(() => {
		setChickenFootURL(player?.chickenFootURL);
	}, [player]);

	useEffect(() => {
		if (chickenFootURL === undefined) {
			return;
		}
		client.SetChickenFoot(code, chickenFootURL);
	}, [chickenFootURL, client, code]);

	function passTurn() {
		setSelectedTile(undefined);
		if (chickenFootURL === undefined) {
			setShowVisionQuestModal(true);
		}
		
		client.Pass(code, {
			selected_x: playA !== undefined ? playA.x : -1,
			selected_y: playA !== undefined ? playA.y : -1,
		}).then((resp) => {
			console.log("passed");
		}).catch((error) => {
			console.error("error", error);
			setPlayErrorMessage(error.data.error);
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
		<div className="h-full " onClick={() => setPlayErrorMessage("")}>
			<div className="flex justify-end items-center mb-4 pl-3 pr-3">
				<span className="hidden md:block text-left text-5xl font-bold mr-auto">
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

			<div className="flex justify-center items-center gap-4 max-h-32">
				{opponents.map((o, i) => (
					<div key={i} className="flex-1 overflow-x-auto">
						<Hand
							player={o}
							players={players}
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
			<div className="flex gap-4 justify-center max-h-[75vh]">
				<span className="w-96 hidden lg:block">
					<History history={gameHistory} />
				</span>
				<div className="border-black border-8 flex-1 overflow-auto">
					<Board
						width={boardWidth} height={boardHeight}
						tiles={laidTiles}
						spacer={spacer}
						lineHeads={lineHeads}
						selectedTile={selectedTile}
						playTile={playTile}
						playSpacer={playSpacer}
						clearSpacer={clearSpacer}
						chickenFeet={chickenFeet}
						chickenFeetURLs={chickenFeetURLs}
						indicated={indicated}
						setIndicated={setIndicated}
						playerTurn={myTurn}
						activePlayer={roundInProgress && players[turnIndex]}
						hints={hints}
						playA={playA}
						setPlayA={setPlayA}
						spacerHints={player?.spacer_hints}
					/>
					<WhyNot message={playErrorMessage} />
				</div>
				<span className="w-40 hidden lg:block">
					<History history={roundHistory} />
				</span>
			</div>
			{player &&
				<div className="flex justify-center items-center gap-4">
					<div className="overflow-x-auto w-full">
						<Hand
							player={player}
							players={players}
							name={playerName}
							hidden={false}
							selectedTile={selectedTile}
							setSelectedTile={setSelectedTile}
							playerTurn={myTurn}
							drawTile={drawTile}
							passTurn={passTurn}
							roundInProgress={roundInProgress}
							hintedTiles={hintedTiles}
							hintedSpacer={player.spacer_hints}
							bagCount={bagCount}
						/>
					</div>
				</div>
			}
			{showVisionQuestModal && (
				<VisionQuest
					onClose={() => setShowVisionQuestModal(false)}
					isOpen={showVisionQuestModal}
					setChickenFootURL={setChickenFootURL}
				/>
			)}
		</div>
	);
}

export default Game;