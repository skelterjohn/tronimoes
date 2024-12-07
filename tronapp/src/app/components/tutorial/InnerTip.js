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
	const tooltipRef = useRef(null);
	const { activeRef, setActiveRef, messageDones, setMessageDones } = useContext(TipContext);
	const [parentBounds, setParentBounds] = useState(null);
	const [position, setPosition] = useState({ left: -1000, top: -1000 });

	useEffect(() => {
		const bundleDone = messageDones.get(bundle.message);

		if (bundleDone) {
			setActiveRef((prev) => prev === tipRef?.current ? null : prev);
			return;
		}

		setActiveRef((prev) => prev === null ? tipRef?.current : prev);
		setParentBounds(tipRef.current?.parentElement?.getBoundingClientRect());
	}, [messageDones, tipRef?.current, setActiveRef]);

	useEffect(() => {
		if (tipRef.current && tooltipRef.current && parentBounds) {
			const tooltipRect = tooltipRef.current.getBoundingClientRect();
			const viewportWidth = window.innerWidth;
			const viewportHeight = window.innerHeight;

			let left = parentBounds.left + (parentBounds.width / 2);
			let top = parentBounds.top + (parentBounds.height / 2);

			left = Math.min(left, viewportWidth - (tooltipRect.width / 2) - 10);
			top = Math.min(top, viewportHeight - (tooltipRect.height / 2) - 10);

			left = Math.max(left, tooltipRect.width / 2 + 10);
			top = Math.max(top, tooltipRect.height / 2 + 10);

			setPosition({ left, top });
		}
	}, [parentBounds, tipRef.current, tooltipRef.current]);

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
					ref={tooltipRef}
					onClick={(event) => {
						event.stopPropagation();
						setMessageDones((prev) => {
							const newMap = new Map(prev);
							newMap.set(bundle.message, true);
							return newMap;
						});
						setActiveRef(null);
					}}
					className="absolute z-50 border border-black bg-white rounded-lg shadow-lg"
					style={{
						position: 'fixed',
						left: position.left,
						top: position.top,
						//transform: 'translate(-50%, -50%)'
					}}
				>	
					<i className="rotate-45 fas fa-arrow-left"></i>
					<div className="p-2 text-black">{bundle.message}</div>
				</div>,
				document.body
			)}
		</div>
	);
}
