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

export default function Home() {
	const [errorMessage, setErrorMessage] = useState(null);
	const [showSignIn, setShowSignIn] = useState(false);
	const { userInfo, setUserInfo, loading, error, client } = useGameState();

	const [recentSummaries, setRecentSummaries] = useState([]);
	const [pickupSummaries, setPickupSummaries] = useState([]);
	const [activeSummaries, setActiveSummaries] = useState([]);

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

		async function poll() {
			let cursor = 0;
			while (!cancelled) {
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
		};
	}, [client]);

	return (
		<main onClick={() => setErrorMessage(null)} className="font-game relative min-h-screen w-screen bg-slate-800">
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
			<div className="absolute left-3/4 top-1/4 -translate-x-1/2 w-fit">
				<Marquis title="recent games" summaries={recentSummaries} />
			</div>
			<div className="absolute left-1/4 top-1/4 -translate-x-1/2 w-fit space-y-4">
				<Marquis title="pickup games" summaries={pickupSummaries} />
				<Marquis title="active games" summaries={activeSummaries} />
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
