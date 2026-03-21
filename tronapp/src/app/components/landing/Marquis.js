"use client";

import Link from "next/link";

export default function Marquis({ title, summaries = [] }) {
	if (summaries.length === 0) {
		return null;
	}

	return (
		<div className="font-game w-full max-w-sm rounded-lg border border-white bg-black/70 p-4 text-white">
			<div className="mb-2 text-lg tracking-wider">{title}</div>
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
		</div>
	);
}
