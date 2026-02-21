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
			title="game options"
			open={isOpen}
			onCancel={onClose}
			className="font-game"
			footer={[
				<Button key="ok" type="primary" onClick={onClose}>
					Ok
				</Button>,
			]}
		>
			<Space orientation="vertical" size="middle" style={{ width: '100%', padding: '20px 0' }}>
				<div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
					<span>roodle crush</span>
					<Checkbox
						checked={roodle}
						onChange={(e) => setRoodle(e.target.checked)}
						className="font-game"
					/>
				</div>
			</Space>
		</Modal>
	);
};

export default Options;
