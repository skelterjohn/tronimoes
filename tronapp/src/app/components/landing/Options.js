import React, { useState, useEffect } from 'react';
import { Modal, Space, Select, Checkbox } from 'antd';
import Button from '@/app/components/Button';
import { useGameState } from '../GameState';
	
const Options = ({ isOpen, onClose }) => {
	const [randomChoice, setRandomChoice] = useState(false);
	const [gibbsPlanner, setGibbsPlanner] = useState(false);
	const { options, setOptions } = useGameState();

	useEffect(() => {
		if (isOpen) {
			setRandomChoice(!!options.randomChoice);
			setGibbsPlanner(!!options.gibbsPlanner);
		}
	}, [isOpen, options.randomChoice, options.gibbsPlanner]);

	const handleRandomChoiceChange = (e) => {
		const value = e.target.checked;
		setRandomChoice(value);
		setOptions((prev) => ({ ...prev, randomChoice: value }));
	};

	const handleGibbsPlannerChange = (e) => {
		const value = e.target.checked;
		setGibbsPlanner(value);
		setOptions((prev) => ({ ...prev, gibbsPlanner: value }));
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
						<Checkbox
							checked={gibbsPlanner}
							onChange={handleGibbsPlannerChange}
							className="font-game"
						/>
					</div>
					<div style={{ fontSize: 12, color: 'rgba(0,0,0,0.45)', marginTop: 4 }}>
						Round out a 4-player game with bots.
					</div>
				</div>
			</Space>
		</Modal>
	);
};

export default Options;
