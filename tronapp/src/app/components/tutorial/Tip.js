"use client";

import { useState } from "react";
import InnerTip from "./InnerTip";
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
	const [show, setShow] = useState(false);
	
	return { show, setShow, message };
}
