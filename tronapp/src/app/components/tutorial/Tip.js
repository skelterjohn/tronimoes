"use client";

import { useState, useRef, useEffect, createContext, useContext } from "react";
import { createPortal } from "react-dom";
import { useGameState } from "@/app/components/GameState";

const TipContext = createContext();

export function TipProvider({ children }) {
	const [messageRefs] = useState(new Map());
	const [activeTip, setActiveTip] = useState(null);

	return (
		<TipContext.Provider value={{ messageRefs, activeTip, setActiveTip }}>
			{children}
		</TipContext.Provider>
	);
}

export default function Tip({ bundle }) {
	const { tutorial } = useGameState();
	const { messageRefs, activeTip, setActiveTip } = useContext(TipContext);

	const tipRef = useRef(null);

	const [parentBounds, setParentBounds] = useState(null);


	useEffect(() => {
		if (!tutorial || !bundle.show) {
			return;
		}
		
		if (bundle.done && activeTip === bundle.message) {
			setActiveTip((prev) => prev === bundle.message ? null : prev);
		}

		if (!bundle.done && !activeTip) {
			setActiveTip((prev) => prev || bundle.message);
		}
	}, [tutorial, bundle.show, bundle.done, bundle.message, activeTip, setActiveTip]);

	useEffect(() => {
		if (!tutorial) {
			return;
		}

		// First parent to show this message wins.
		if (!tipRef?.current) {
			return;
		}
		if (!messageRefs.has(bundle.message)) {
			messageRefs.set(bundle.message, tipRef);
		}
		setParentBounds(tipRef.current?.parentElement?.getBoundingClientRect());
	}, [bundle, tipRef, messageRefs, tutorial]);
	
	if (!tutorial) {
		return null;
	}

	if (!bundle.show || bundle.done) {
		return null;
	}

	if (activeTip && activeTip !== bundle.message) {
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
						setActiveTip(null);
					}}
					className="absolute border border-black bg-white rounded-lg p-2 shadow-lg z-50 text-black"
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
