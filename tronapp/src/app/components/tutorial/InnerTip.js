"use client";

import { useState, useRef, useContext, useEffect, createContext } from "react";
import { createPortal } from "react-dom";


export const TipContext = createContext();

export function TipProvider({ children }) {
	const [activeRef, setActiveRef] = useState(null);
	const [messageDones, setMessageDones] = useState(new Map());

	return (
		<TipContext.Provider value={{ activeRef, setActiveRef, messageDones, setMessageDones }}>
			{children}
		</TipContext.Provider>
	);
}

export default function InnerTip({ bundle }) {
	const tipRef = useRef(null);
	const { activeRef, setActiveRef, messageDones, setMessageDones } = useContext(TipContext);
	const [parentBounds, setParentBounds] = useState(null);

	useEffect(() => {
		const bundleDone = messageDones.get(bundle.message);

		if (bundleDone) {
			setActiveRef((prev) => {
				if (prev === tipRef?.current) {
					return null;
				}
				return prev;
			});
			return;
		}

		setActiveRef((prev) => prev === null ? tipRef?.current : prev);
		setParentBounds(tipRef.current?.parentElement?.getBoundingClientRect());
	}, [messageDones, tipRef?.current, setActiveRef]);
	
	if (messageDones.get(bundle.message)) {
		return null;
	}

	if (activeRef && activeRef !== tipRef?.current) {
		return null;
	}

	return (
		<div ref={tipRef}>
			{createPortal(
				<div
					onClick={(event) => {
						event.stopPropagation();
						setMessageDones((prev) => {
							const newMap = new Map(prev);
							newMap.set(bundle.message, true);
							return newMap;
						});
						setActiveRef(null);
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
