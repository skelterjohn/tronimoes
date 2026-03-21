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

export default function Marquis({
	title,
	summaries = [],
	namesOnly = false,
	showEmpty = false,
	hideTitle = false,
	className = "",
	inlinePlayerInfo = false,
}) {
	const shell = `font-game w-full max-w-sm rounded-lg border border-white bg-black/70 p-4 text-white ${className}`.trim();

	if (summaries.length === 0 && !showEmpty) {
		return null;
	}

	if (summaries.length === 0) {
		return (
			<div className={shell}>
				{!hideTitle && <div className="mb-2 text-lg tracking-wider">{title}</div>}
				<p className="text-sm opacity-70">no games</p>
			</div>
		);
	}

	return (
		<div className={shell}>
			{!hideTitle && <div className="mb-2 text-lg tracking-wider">{title}</div>}
			<ul className="space-y-2 text-sm">
				{summaries.map((summary) => {
					if (!summary?.code) {
						return null;
					}
					const lines = playerScoreLines(summary);
					const link = (
						<Link className="underline underline-offset-2 hover:opacity-80" href={`/gameboard/${summary.code}`}>
							#{summary.code}
						</Link>
					);
					if (inlinePlayerInfo) {
						return (
							<li key={summary.code} className="flex gap-2">
								<span className="shrink-0">{link}</span>
								<div className="min-w-0 flex-1 text-right text-xs leading-snug text-white/90">
									{lines.length === 0 ? (
										<span className="text-white/50">No players</span>
									) : (
										lines.map(({ name, score }) => (
											<div key={name}>
												{namesOnly ? name : `${name}: ${score}`}
											</div>
										))
									)}
								</div>
							</li>
						);
					}
					return (
						<li key={summary.code}>
							<Tooltip
								placement="right"
								title={
									lines.length > 0 ? (
										<ul className="m-0 max-w-xs list-none p-0">
											{lines.map(({ name, score }) => (
												<li key={name}>
													{namesOnly ? name : `${name}: ${score}`}
												</li>
											))}
										</ul>
									) : (
										<span className="opacity-80">No players</span>
									)
								}
							>
								{link}
							</Tooltip>
						</li>
					);
				})}
			</ul>
		</div>
	);
}
