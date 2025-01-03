"use client";

import Image from 'next/image';
import { useState, useEffect } from 'react';
import Joiner from './components/landing/Joiner';
import SignIn from './components/landing/SignIn';
import { auth } from "@/config";
import { signOut } from "firebase/auth";
import Error from './components/landing/Error';
import { useGameState } from './components/GameState';

export default function Home() {
	const [errorMessage, setErrorMessage] = useState(null);
	const [showSignIn, setShowSignIn] = useState(false);
	const { userInfo, setUserInfo, loading, error } = useGameState();

	useEffect(()=> {
		setErrorMessage(error?.message);
	}, [error]);

	return (
		<main onClick={() => setErrorMessage(null)} className="relative min-h-screen w-screen bg-slate-800">
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
