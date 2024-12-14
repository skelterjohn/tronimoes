import { useState, useEffect } from "react";
import Tile from '../board/Tile';
import ChickenFoot from '../board/ChickenFoot';
import Reaction from "./Reaction";

function Opponent({
		player, players,
		dead = false, turnIndex
	}) {
	const [myTurn, setMyTurn] = useState(false);
	const [handBackground, setHandBackground] = useState("bg-white");

	useEffect(() => {
		const colorMap = {
			red: "bg-red-600",
			blue: "bg-blue-600",
			green: "bg-green-600",
			yellow: "bg-yellow-600",
			orange: "bg-orange-600",
			fuchsia: "bg-fuchsia-600",
			white: "bg-white"
		};
		setHandBackground(colorMap[player?.color]);
	}, [player]);

	useEffect(() => {
		setMyTurn(player?.name === players[turnIndex]?.name);
	}, [turnIndex, player, players]);

	const [reactURL, setReactURL] = useState(undefined);
	const [showReaction, setShowReaction] = useState(false);
	useEffect(() => {
		setReactURL(player?.reactURL);
	}, [player]);

	useEffect(() => {
		setShowReaction(reactURL !== undefined);
	}, [reactURL]);
	
	const [killedPlayers, setKilledPlayers] = useState([]);
	useEffect(() => {
		setKilledPlayers(player?.kills?.map(k =>  players.find(p => p.name === k)));
	}, [player, players]);
	

	return (
		<div className={`h-full flex flex-col items-center ${myTurn ? "border-2 border-black " + handBackground : ""}`}>
			<div className="w-full text-center font-bold ">
				<div className="flex flex-row items-center justify-center gap-2">
					{killedPlayers?.map(kp => (
						<div key={kp.name} className="relative w-[2rem] h-[2rem] inline-block align-middle">
							<div className="absolute inset-0">
								<ChickenFoot url={kp.chickenFootURL} color={kp.color} />
							</div>
						</div>
					))}
					<span>
						{player?.name} - ({player?.score})
						{player?.chickenFoot && " (footed)"}
						{player?.ready && " (ready)"}
					</span>
					{showReaction && (
						<Reaction 
							url={reactURL}
							setShow={setShowReaction}
						/>
					)}
					{!player?.chickenFoot && !player?.dead &&
						<div className="relative w-[2rem] h-[2rem] inline-block align-middle">
							<div className="absolute inset-0">
								<ChickenFoot url={player.chickenFootURL} color={player.color} />
							</div>
						</div>
					}
				</div>
			</div>
			<div className="flex flex-row items-center gap-1">
				<div className="w-[1rem]">
					<Tile
						color={player?.color}
						pipsa={0}
						pipsb={0}
						back={true}
						dead={dead}
					/>
				</div>
				<div>x{player?.hand?.length}</div>
			</div>
		</div>
	);
}

export default Opponent;