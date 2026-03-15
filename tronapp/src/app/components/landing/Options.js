import React, { useState, useEffect } from 'react';
import { Modal, Space, Select } from 'antd';
import Button from '@/app/components/Button';
import { useGameState } from '../GameState';

const AGENT_OPTIONS = [0, 2, 3, 4, 5, 6];

const Options = ({ isOpen, onClose }) => {
	const [roundOut, setRoundOut] = useState(0);
	const { options, setOptions } = useGameState();

	useEffect(() => {
		if (isOpen) {
			setRoundOut(options.agent_round_out ?? 0);
		}
	}, [isOpen, options.agent_round_out]);

	const handleAgentRoundOutChange = (value) => {
		setRoundOut(value);
		setOptions((prev) => ({ ...prev, agent_round_out: value }));
	};

	return (
		<Modal
			title="create game options"
			open={isOpen}
			onCancel={onClose}
			className="font-game"
			footer={
				<div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: 16, width: '100%' }}>
					<span style={{ display: 'block', flex: '1 1 0', fontSize: 12, color: 'rgba(0,0,0,0.45)', minWidth: 0, textAlign: 'left' }}>
						These options are only used if you create the game.
					</span>
					<span style={{ flexShrink: 0 }}>
						<Button key="ok" type="primary" className="game-btn" onClick={onClose}>
							Ok
						</Button>
					</span>
				</div>
			}
		>
			<Space orientation="vertical" size="middle" style={{ width: '100%', padding: '20px 0' }}>
				<div>
					<div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
						<span>install agents</span>
						<Select
							value={roundOut}
							onChange={handleAgentRoundOutChange}
							options={AGENT_OPTIONS.map((n) => ({ value: n, label: String(n) }))}
							className="font-game"
							style={{ width: 72 }}
						/>
					</div>
					<div style={{ fontSize: 12, color: 'rgba(0,0,0,0.45)', marginTop: 4 }}>
						Round out a game with bots to reach this many players.
					</div>
				</div>
			</Space>
		</Modal>
	);
};

export default Options;
