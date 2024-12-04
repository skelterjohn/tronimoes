"use client";

import { useState, useRef } from "react";
import { createPortal } from "react-dom";
import { useGameState } from "@/app/components/GameState";

// Track which tips have already been given.
const messageRefs = new Map();
// Add a variable to track the currently active tip
let activeMessage = undefined;

export default function Tip({ bundle }) {
	const { tutorial } = useGameState();

	if (activeMessage === undefined) {
		activeMessage = bundle.message;
	}
	if (activeMessage !== bundle.message) {
		return null;
	}

	if (bundle.done || !tutorial || !bundle.show || !bundle.parentRef?.current) {
		return null;
	}

	const tipRef = useRef(null);
	// First parent to show this message wins.
	if (!messageRefs.has(bundle.message)) {
		messageRefs.set(bundle.message, tipRef);
	}
	if (messageRefs.get(bundle.message) !== tipRef) {
		return null;
	}
	const parentBounds = bundle.parentRef.current.getBoundingClientRect();
	
	return createPortal(
		<div
			onClick={() => {
				bundle.setDone(true)
				activeMessage = undefined;
			}}
			className="absolute bg-white rounded-lg p-2 shadow-lg z-50 text-black"
			style={{
				position: 'fixed',
				left: parentBounds.left + (parentBounds.width / 2),
				top: parentBounds.top + (parentBounds.height / 2),
				transform: 'translate(-50%, -50%)'
			}}
		>
			{bundle.message}
		</div>,
		document.body
	);
}

export function makeTipBundle(message) {
	const [done, setDone] = useState(false);
	const [show, setShow] = useState(false);
	const parentRef = useRef(null);
	
	return { done, setDone, show, setShow, message, parentRef };
}
