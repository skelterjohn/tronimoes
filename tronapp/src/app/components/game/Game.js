"use client";

import {Row, Col} from 'react-bootstrap';
import { useState, useEffect } from 'react';

import Board from '../board/Board';
import Hand from './Hand';

function Game({}) {

	const [opponents, setOpponents] = useState([
		{
			color: "red",
			count: 4,
			dead: true,
		},
		{
			color: "blue",
			count: 3,
			dead: false,
		},
	]);

	const [playerHand, setPlayerHand] = useState([{a:1, b:2}, {a:3, b:12}]);
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

	const [playerTurn, setPlayerTurn] = useState(true);
	const [selectedTile, setSelectedTile] = useState({a:1, b:2});

	function startTurn() {
		setPlayerTurn(true);
		setSelectedTile(undefined);
	}

	function playTile(tile) {
		tile.color = "green";
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
						color={o.color}
						hidden={true}
						dead={o.dead}
						tiles={Array(o.count).fill({a:0, b:0})}
					/>
				</Col>
			))}
		</Row>
		<Row>
			<Board
				width={10} height={11}
				tiles={laidTiles}
				selectedTile={selectedTile}
				playTile={playTile}
			/>
		</Row>
		<Row>
			<Hand
				color="green"
				tiles={playerHand}
				selectedTile={selectedTile}
				setSelectedTile={setSelectedTile}
				playerTurn={playerTurn}
			/>
		</Row>
	</div>;
}

export default Game;