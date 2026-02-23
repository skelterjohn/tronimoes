import React, { useState, useEffect } from 'react';
import { Modal, Space, Select, Checkbox } from 'antd';
import Button from '@/app/components/Button';
import { useGameState } from '../GameState';
	
const Options = ({ isOpen, onClose }) => {
	const [roodle, setRoodle] = useState(false);
	const { options, setOptions } = useGameState();

	useEffect(() => {
		setRoodle(options.roodle);
	}, [options]);

	useEffect(() => {
		setOptions({ ...options, roodle });
	}, [roodle]);

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
						<span>roodle crush</span>
						<Checkbox
							checked={roodle}
							onChange={(e) => setRoodle(e.target.checked)}
							className="font-game"
						/>
					</div>
					<div style={{ fontSize: 12, color: 'rgba(0,0,0,0.45)', marginTop: 4 }}>
						Add the roodle crush bot to your game. (not implemented yet)
					</div>
				</div>
			</Space>
		</Modal>
	);
};

export default Options;
