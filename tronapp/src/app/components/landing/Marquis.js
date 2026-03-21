"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useGameState } from "@/app/components/GameState";

export const MarquisType = Object.freeze({
	RECENT: "recent",
	ACTIVE: "active",
	PICKUP: "pickup",
});

function fetchScoreboardsByType(client, marquisType, updated) {
	var scoreboards = client.ListScoreboards(updated);
	switch (marquisType) {
	case MarquisType.ACTIVE:
		return scoreboards.Active;
	case MarquisType.PICKUP:
		return scoreboards.Pickup;
	case MarquisType.RECENT:
	default:
		return scoreboards.Recent;
	}
}

export default function Marquis({ title, marquisType = MarquisType.RECENT, refreshCadenceMs = 60000 }) {
	const { client } = useGameState();
	const [summaries, setSummaries] = useState([]);
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

		async function poll() {
			let cursor = 0;
			setLoading(true);
			while (!cancelled) {
				try {
					const scoreboards = await client.ListScoreboards(cursor);
					console.log("scoreboards", scoreboards);
					if (cancelled) {
						return;
					}
					var summaries = [];
					switch (marquisType) {
					case MarquisType.ACTIVE:
						summaries = scoreboards.active;
						break;
					case MarquisType.PICKUP:
						summaries = scoreboards.pickup;
						break;
					case MarquisType.RECENT:
						summaries = scoreboards.recent;
						break;
					}
					cursor = scoreboards.updated;
					setSummaries(summaries);
					setError(null);
				} catch (err) {
					if (cancelled) {
						return;
					}
					const nextError = err?.data?.error || err?.message || "Could not load games";
					setError((prevError) => (prevError === nextError ? prevError : nextError));
					// Avoid tight retry loops when the server returns quickly with errors.
					await new Promise((resolve) => setTimeout(resolve, refreshCadenceMs));
				} finally {
					if (!cancelled) {
						setLoading(false);
					}
				}
			}
		}
		poll();

		return () => {
			cancelled = true;
		};
	}, [client, marquisType, refreshCadenceMs]);

	if (summaries.length === 0) {
		return null;
	}

	return (
		<div className="font-game w-full max-w-sm rounded-lg border border-white bg-black/70 p-4 text-white">
			<div className="mb-2 text-lg tracking-wider">{title}</div>
			{!loading && error && <div className="text-sm text-red-300">{error}</div>}
			{!loading && !error && summaries.length > 0 && (
				<ul className="space-y-1 text-sm">
					{summaries.map((summary) => (
						summary?.code ? (
						<li key={summary.code}>
							<Link className="underline underline-offset-2 hover:opacity-80" href={`/gameboard/${summary.code}`}>
								#{summary.code}
							</Link>
						</li>
						) : null
					))}
				</ul>
			)}
		</div>
	);
}
