"use client";

import Link from "next/link";

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
}) {
	const shell = `font-game w-full md:w-max md:min-w-[15vw] rounded-lg border border-white bg-black/70 p-4 text-white ${className}`.trim();

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
					return (
						<li key={summary.code} className="flex items-start justify-between gap-2">
							<span className="shrink-0">
								<Link className="underline underline-offset-2 hover:opacity-80" href={`/gameboard/${summary.code}`}>
									#{summary.code}
								</Link>
							</span>
							<div className="shrink-0 text-right text-xs leading-snug text-white/90">
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
				})}
			</ul>
		</div>
	);
}
