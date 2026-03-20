"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useGameState } from "@/app/components/GameState";

export const MarquisType = Object.freeze({
	RECENT: "recent",
	ACTIVE: "active",
	PICKUP: "pickup",
});

function fetchScoreboardsByType(client, marquisType) {
	switch (marquisType) {
	case MarquisType.ACTIVE:
		return client.GetActiveScoreboards();
	case MarquisType.PICKUP:
		return client.GetPickupScoreboards();
	case MarquisType.RECENT:
	default:
		return client.GetRecentScoreboards();
	}
}

export default function Marquis({ title, marquisType = MarquisType.RECENT, refreshCadenceMs = 3000 }) {
	const { client } = useGameState();
	const [codes, setCodes] = useState([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState(null);

	useEffect(() => {
		let cancelled = false;

		if (!client) {
			setLoading(false);
			return () => {
				cancelled = true;
			};
		}

		const refresh = (isInitialLoad = false) => {
			if (isInitialLoad) {
				setLoading(true);
			}
			fetchScoreboardsByType(client, marquisType).then((games) => {
				if (cancelled) {
					return;
				}
				const nextCodes = Array.isArray(games) ? games.map((game) => game.code).filter(Boolean) : [];
				setCodes(nextCodes);
				setError(null);
			}).catch((err) => {
				if (cancelled) {
					return;
				}
				const nextError = err?.data?.error || err?.message || "Could not load games";
				setError((prevError) => (prevError === nextError ? prevError : nextError));
			}).finally(() => {
				if (!cancelled && isInitialLoad) {
					setLoading(false);
				}
			});
		};

		refresh(true);
		const interval = setInterval(() => {
			refresh(false);
		}, refreshCadenceMs);

		return () => {
			cancelled = true;
			clearInterval(interval);
		};
	}, [client, marquisType, refreshCadenceMs]);

	if (codes.length === 0) {
		return null;
	}

	return (
		<div className="font-game w-full max-w-sm rounded-lg border border-white bg-black/70 p-4 text-white">
			<div className="mb-2 text-lg tracking-wider">{title}</div>
			{!loading && error && <div className="text-sm text-red-300">{error}</div>}
			{!loading && !error && codes.length > 0 && (
				<ul className="space-y-1 text-sm">
					{codes.map((code) => (
						<li key={code}>
							<Link className="underline underline-offset-2 hover:opacity-80" href={`/gameboard/${code}`}>
								#{code}
							</Link>
						</li>
					))}
				</ul>
			)}
		</div>
	);
}
