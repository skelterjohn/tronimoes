"use client";

import { useState, useRef, useEffect, createContext, useContext } from "react";
import { createPortal } from "react-dom";
import { useGameState } from "@/app/components/GameState";

const TipContext = createContext();

export function TipProvider({ children }) {
	const [messageRefs] = useState(new Map());

	return (
		<TipContext.Provider value={{ messageRefs }}>
			{children}
		</TipContext.Provider>
	);
}

export default function Tip({ bundle }) {

	const { tutorial } = useGameState();

	const { messageRefs } = useContext(TipContext);

	const tipRef = useRef(null);

	const [parentBounds, setParentBounds] = useState(null);

	useEffect(() => {
		// First parent to show this message wins.
		if (!tipRef?.current) {
			return;
		}
		if (!messageRefs.has(bundle.message)) {
			messageRefs.set(bundle.message, tipRef);
		}
		setParentBounds(tipRef.current?.parentElement?.getBoundingClientRect());
	}, [bundle, tipRef, messageRefs]);

	if (!tutorial) {
		return null;
	}

	if (!bundle.show || bundle.done) {
		return null;
	}

	if (tipRef?.current && messageRefs.get(bundle.message) !== tipRef) {
		return null;
	}
	
	return (
		<div ref={tipRef}>
			{parentBounds && createPortal(
				<div
					onClick={() => {
						bundle.setDone(true);
					}}
					className="absolute bg-white rounded-lg p-2 shadow-lg z-50 text-black"
					style={{
						position: 'fixed',
						left: parentBounds?.left + (parentBounds?.width / 2) || -1000,
						top: parentBounds?.top + (parentBounds?.height / 2) || -1000,
						transform: 'translate(-50%, -50%)'
					}}
				>
					{bundle.message}
				</div>,
				document.body
			)}
		</div>
	);
}

export function makeTipBundle(message) {
	const [done, setDone] = useState(false);
	const [show, setShow] = useState(false);
	
	return { done, setDone, show, setShow, message };
}
