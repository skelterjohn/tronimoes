"use client";

import { useState, useContext, useCallback } from "react";
import InnerTip, { TipContext } from "./InnerTip";
import { useGameState } from "@/app/components/GameState";

export default function Tip({ bundle }) {
	const { tutorial } = useGameState();
	
	if (!tutorial || !bundle.show) {
		return null;
	}

	return (
		<InnerTip bundle={bundle} />
	);
}

export function useTipBundle(message) {

	const { messageDones } = useContext(TipContext);

	const [show, setShow] = useState(false);
	
	const done = useCallback(() => {
		return messageDones.get(message);
	}, [messageDones, message]);

	return { show, setShow, message, done };
}
