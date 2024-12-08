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
		
		// Add retry logic for getting parent bounds
		const getBounds = (retry = false) => {
			const bounds = tipRef.current?.parentElement?.getBoundingClientRect();
			if (bounds && (bounds.width === 0 && bounds.height === 0)) {
				// Retry after a short delay if bounds are invalid
				if (retry) {
					setTimeout(getBounds, 100);
				}
			} else {
				setParentBounds(bounds);
			}
		};
		
		getBounds(true);
	}, [messageDones, tipRef?.current, setActiveRef]);

	useEffect(() => {
		if (tipRef.current && tooltipRef.current && parentBounds) {
			const tooltipRect = tooltipRef.current.getBoundingClientRect();
			const viewportWidth = window.innerWidth;
			const viewportHeight = document.documentElement.clientHeight;

			// Center horizontally
			let left = parentBounds.left + (parentBounds.width / 2);
			left = Math.min(left, viewportWidth - tooltipRect.width - 10);
			left = Math.max(left, 10);

			// Always position at bottom of viewport
			const mobileBottomPadding = window.innerWidth < 768 ? 60 : 10;
			const top = viewportHeight - tooltipRect.height - mobileBottomPadding;

			setPosition({ left, top });
		}
	}, [parentBounds, tipRef.current, tooltipRef.current]);

	if (messageDones.get(bundle.message)) {
		return null;
	}

	if (activeRef && activeRef !== tipRef?.current) {
		return null;
	}

	console.log(parentBounds);
	return (
		<div ref={tipRef}>
			{parentBounds && parentBounds?.width > 0 && createPortal(
				<>
					<svg 
						style={{
							position: 'fixed',
							top: 0,
							left: 0,
							width: '100%',
							height: '100%',
							pointerEvents: 'none',
							zIndex: 49
						}}
					>
						<line
							x1={parentBounds.left + (parentBounds.width / 2)}
							y1={parentBounds.top + (parentBounds.height / 2)}
							x2={position.left + (tooltipRef.current?.offsetWidth || 0) / 2}
							y2={position.top}
							stroke="black"
							strokeWidth="6"
						/>
						<line
							x1={parentBounds.left + (parentBounds.width / 2)}
							y1={parentBounds.top + (parentBounds.height / 2)}
							x2={position.left + (tooltipRef.current?.offsetWidth || 0) / 2}
							y2={position.top}
							stroke="white"
							strokeWidth="4"
						/>
						<circle 
							cx={parentBounds.left + (parentBounds.width / 2)}
							cy={parentBounds.top + (parentBounds.height / 2)}
							r="8"
							fill="white"
							stroke="black"
							strokeWidth="2"
						/>
					</svg>
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
						}}
					>	
						<div className="p-2 text-black">{bundle.message}</div>
					</div>
				</>,
				document.body
			)}
		</div>
	);
}
