"use client";

import Image from 'next/image';
import { useState, useEffect } from 'react';
import Joiner from './components/landing/Joiner';
import Marquis from './components/landing/Marquis';
import SignIn from './components/landing/SignIn';
import { auth } from "@/config";
import { signOut } from "firebase/auth";
import Error from './components/landing/Error';
import { useGameState } from './components/GameState';
import Link from 'next/link';

const SCOREBOARD_POLL_ERROR_BACKOFF_MS = 60000;
const SCOREBOARD_IDLE_HALT_MS = 5 * 60 * 1000;

export default function Home() {
	const [errorMessage, setErrorMessage] = useState(null);
	const [showSignIn, setShowSignIn] = useState(false);
	const { userInfo, setUserInfo, loading, error, client } = useGameState();

	const [recentSummaries, setRecentSummaries] = useState([]);
	const [pickupSummaries, setPickupSummaries] = useState([]);
	const [activeSummaries, setActiveSummaries] = useState([]);
	const [mobilePanel, setMobilePanel] = useState(null);

	useEffect(()=> {
		setErrorMessage(error?.message);
	}, [error]);

	useEffect(() => {
		let cancelled = false;

		if (!client) {
			return () => {
				cancelled = true;
			};
		}

		const lastActivityRef = { current: Date.now() };
		let idleWake = null;
		let visibleWake = null;

		function noteActivity() {
			lastActivityRef.current = Date.now();
			if (idleWake) {
				const wake = idleWake;
				idleWake = null;
				wake();
			}
		}

		const activityOpts = { passive: true };
		const activityEvents = ["pointerdown", "keydown", "touchstart", "scroll", "wheel"];
		for (const type of activityEvents) {
			window.addEventListener(type, noteActivity, activityOpts);
		}
		function onVisibility() {
			if (!document.hidden) {
				noteActivity();
				if (visibleWake) {
					const wake = visibleWake;
					visibleWake = null;
					wake();
				}
			}
		}
		document.addEventListener("visibilitychange", onVisibility);

		async function poll() {
			let cursor = 0;
			while (!cancelled) {
				if (document.hidden) {
					await new Promise((resolve) => {
						if (cancelled) {
							resolve();
							return;
						}
						visibleWake = resolve;
					});
					if (cancelled) {
						return;
					}
					continue;
				}
				if (Date.now() - lastActivityRef.current >= SCOREBOARD_IDLE_HALT_MS) {
					await new Promise((resolve) => {
						if (cancelled) {
							resolve();
							return;
						}
						idleWake = resolve;
					});
					if (cancelled) {
						return;
					}
					continue;
				}
				try {
					const scoreboards = await client.ListScoreboards(cursor);
					if (cancelled) {
						return;
					}
					setRecentSummaries(Array.isArray(scoreboards.recent) ? scoreboards.recent : []);
					setPickupSummaries(Array.isArray(scoreboards.pickup) ? scoreboards.pickup : []);
					setActiveSummaries(Array.isArray(scoreboards.active) ? scoreboards.active : []);
					cursor = scoreboards.updated ?? 0;
				} catch (err) {
					if (cancelled) {
						return;
					}
					await new Promise((resolve) => setTimeout(resolve, SCOREBOARD_POLL_ERROR_BACKOFF_MS));
				}
			}
		}
		poll();

		return () => {
			cancelled = true;
			if (idleWake) {
				const wake = idleWake;
				idleWake = null;
				wake();
			}
			if (visibleWake) {
				const wake = visibleWake;
				visibleWake = null;
				wake();
			}
			for (const type of activityEvents) {
				window.removeEventListener(type, noteActivity, activityOpts);
			}
			document.removeEventListener("visibilitychange", onVisibility);
		};
	}, [client]);

	return (
		<main
			onClick={() => setErrorMessage(null)}
			className={`font-game relative min-h-screen w-screen bg-slate-800 ${
				mobilePanel === null
					? "pb-14 md:pb-0"
					: "pb-[calc(100dvh-max(50vh,12rem))] md:pb-0"
			}`}
		>
			<Image 
				src="/trondude.png"
				alt="Background"
				fill
				className="object-cover z-0"
				priority
			/>
			<div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-fit space-y-4">
				<Joiner userInfo={userInfo} loading={loading} setErrorMessage={setErrorMessage} />
			</div>
			<div className="absolute left-3/4 top-1/4 hidden w-fit -translate-x-1/2 md:block">
				<Marquis title="recent games" summaries={recentSummaries} />
			</div>
			<div className="absolute left-1/4 top-1/4 hidden w-fit -translate-x-1/2 space-y-4 md:block">
				<Marquis title="pickup games" summaries={pickupSummaries} namesOnly />
				<Marquis title="active games" summaries={activeSummaries} />
			</div>

			<div
				className={`fixed bottom-0 left-0 right-0 z-40 box-border flex min-h-0 flex-col bg-transparent pb-[env(safe-area-inset-bottom)] md:hidden ${
					mobilePanel === null ? "h-auto" : "max-h-[calc(100dvh-max(50vh,12rem))]"
				}`}
			>
				{mobilePanel !== null && (
					<div className="min-h-0 flex-1 overflow-y-auto bg-transparent px-2 pt-2">
						{mobilePanel === "recent" && (
							<Marquis title="recent games" summaries={recentSummaries} showEmpty hideTitle className="max-w-none" />
						)}
						{mobilePanel === "pickup" && (
							<Marquis title="pickup games" summaries={pickupSummaries} namesOnly showEmpty hideTitle className="max-w-none" />
						)}
						{mobilePanel === "active" && (
							<Marquis title="active games" summaries={activeSummaries} showEmpty hideTitle className="max-w-none" />
						)}
					</div>
				)}
				<nav
					className="flex shrink-0 items-stretch justify-around gap-1 border-t border-white/30 bg-black/90 px-1 py-1.5 font-game text-xs text-white backdrop-blur-sm"
					aria-label="Game lists"
				>
					<button
						type="button"
						aria-pressed={mobilePanel === "pickup"}
						className={`flex flex-1 flex-col items-center rounded px-2 py-1.5 transition-colors ${
							mobilePanel === "pickup" ? "bg-white/25 font-semibold text-white ring-1 ring-inset ring-white/50" : "hover:bg-white/10"
						}`}
						onClick={(e) => {
							e.stopPropagation();
							setMobilePanel((p) => (p === "pickup" ? null : "pickup"));
						}}
					>
						pickup
					</button>
					<button
						type="button"
						aria-pressed={mobilePanel === "active"}
						className={`flex flex-1 flex-col items-center rounded px-2 py-1.5 transition-colors ${
							mobilePanel === "active" ? "bg-white/25 font-semibold text-white ring-1 ring-inset ring-white/50" : "hover:bg-white/10"
						}`}
						onClick={(e) => {
							e.stopPropagation();
							setMobilePanel((p) => (p === "active" ? null : "active"));
						}}
					>
						active
					</button>
					<button
						type="button"
						aria-pressed={mobilePanel === "recent"}
						className={`flex flex-1 flex-col items-center rounded px-2 py-1.5 transition-colors ${
							mobilePanel === "recent" ? "bg-white/25 font-semibold text-white ring-1 ring-inset ring-white/50" : "hover:bg-white/10"
						}`}
						onClick={(e) => {
							e.stopPropagation();
							setMobilePanel((p) => (p === "recent" ? null : "recent"));
						}}
					>
						recent
					</button>
				</nav>
			</div>
			<div className="absolute top-4 left-4 w-fit text-white">
				<Link href="/rules" target="_blank" rel="noopener noreferrer" className="cursor-pointer underline underline-offset-2">
					rules
				</Link>
			</div>
			<div className="absolute top-4 right-4 w-fit text-white">
				{!loading && userInfo === null && (
					<div onClick={() => setShowSignIn(true)} className="cursor-pointer">
						sign in
					</div>
				)}
				{!loading && userInfo !== null && (
					<div onClick={() => {
						signOut(auth).then(() => {
							setUserInfo(null);
						}).catch((error) => {
							console.error("Sign out error:", error);
						});
					}} className="cursor-pointer">
						sign out
					</div>
				)}
			</div>
			<SignIn userInfo={userInfo} setErrorMessage={setErrorMessage} setUserInfo={setUserInfo} isOpen={showSignIn} onClose={() => setShowSignIn(false)} />
			<Error message={errorMessage} />
		</main>
	);
}
