"use client";

import Link from "next/link";
import { Tooltip } from "antd";

function playerScoreLines(summary) {
	if (!summary?.scoreboard || typeof summary.scoreboard !== "object") {
		return [];
	}
	return Object.entries(summary.scoreboard)
		.sort(([a], [b]) => a.localeCompare(b))
		.map(([name, score]) => ({ name, score }));
}

export default function Marquis({ title, summaries = [] }) {
	if (summaries.length === 0) {
		return null;
	}

	return (
		<div className="font-game w-full max-w-sm rounded-lg border border-white bg-black/70 p-4 text-white">
			<div className="mb-2 text-lg tracking-wider">{title}</div>
			<ul className="space-y-1 text-sm">
				{summaries.map((summary) => {
					if (!summary?.code) {
						return null;
					}
					const lines = playerScoreLines(summary);
					return (
						<li key={summary.code}>
							<Tooltip
								placement="top"
								title={
									lines.length > 0 ? (
										<ul className="m-0 max-w-xs list-none p-0">
											{lines.map(({ name, score }) => (
												<li key={name}>
													{name}: {score}
												</li>
											))}
										</ul>
									) : (
										<span className="opacity-80">No players</span>
									)
								}
							>
								<Link className="underline underline-offset-2 hover:opacity-80" href={`/gameboard/${summary.code}`}>
									#{summary.code}
								</Link>
							</Tooltip>
						</li>
					);
				})}
			</ul>
		</div>
	);
}
