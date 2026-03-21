"use client";

import { Fragment } from "react";
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

	const validSummaries = summaries.filter((s) => s?.code);

	return (
		<div className={shell}>
			{!hideTitle && <div className="mb-2 text-lg tracking-wider">{title}</div>}
			<table className="w-full border-collapse text-sm">
				<tbody>
					{validSummaries.map((summary, idx) => {
						const lines = playerScoreLines(summary);
						const link = (
							<Link className="underline underline-offset-2 hover:opacity-80" href={`/gameboard/${summary.code}`}>
								#{summary.code}
							</Link>
						);

						if (lines.length === 0) {
							return (
								<tr key={summary.code} className={idx > 0 ? "border-t border-white/15" : undefined}>
									<td className="align-top py-1 pr-3">{link}</td>
									<td className="py-1 text-xs text-white/50" colSpan={namesOnly ? 1 : 3}>
										No players
									</td>
								</tr>
							);
						}

						if (namesOnly) {
							return (
								<Fragment key={summary.code}>
									{lines.map(({ name }, i) => (
										<tr
											key={name}
											className={
												i === 0 && idx > 0 ? "border-t border-white/15" : undefined
											}
										>
											{i === 0 && (
												<td className="align-top py-1 pr-3" rowSpan={lines.length}>
													{link}
												</td>
											)}
											<td className="py-1 text-right text-xs leading-snug text-white/90">{name}</td>
										</tr>
									))}
								</Fragment>
							);
						}

						return (
							<Fragment key={summary.code}>
								{lines.map(({ name, score }, i) => (
									<tr
										key={name}
										className={
											i === 0 && idx > 0 ? "border-t border-white/15" : undefined
										}
									>
										{i === 0 && (
											<td className="align-top py-1 pr-3 whitespace-nowrap" rowSpan={lines.length}>
												{link}
											</td>
										)}
										<td className="py-1 text-right text-xs leading-snug text-white/90">{name}</td>
										<td className="px-1 py-1 text-center text-xs leading-snug text-white/90" aria-hidden="true">
											:
										</td>
										<td className="py-1 text-left text-xs leading-snug text-white/90">{score}</td>
									</tr>
								))}
							</Fragment>
						);
					})}
				</tbody>
			</table>
		</div>
	);
}
