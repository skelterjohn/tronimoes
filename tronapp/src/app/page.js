"use client";

import Image from 'next/image';
import { useState } from 'react';
import Joiner from './components/landing/Joiner';
import SignIn from './components/landing/SignIn';
export default function Home() {
	const [userInfo, setUserInfo] = useState(undefined);
	const [showSignIn, setShowSignIn] = useState(false);

	return (
		<main className="relative min-h-screen w-screen bg-slate-800">
			<Image 
				src="/trondude.png"
				alt="Background"
				fill
				className="object-cover z-0"
				priority
			/>
			<div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-fit space-y-4">
				<Joiner/>
			</div>
			<div className="absolute top-4 right-4 w-fit text-white">
				{userInfo === undefined && (
					<div onClick={() => setShowSignIn(true)} className="cursor-pointer">
						sign in
					</div>
				)}
				{userInfo !== undefined && (
					<div onClick={() => setUserInfo(undefined)} className="cursor-pointer">
						sign out
					</div>
				)}
			</div>
			<SignIn userInfo={userInfo} setUserInfo={setUserInfo} isOpen={showSignIn} onClose={() => setShowSignIn(false)} />
		</main>
	);
}
