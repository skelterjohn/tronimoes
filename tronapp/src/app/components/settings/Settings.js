import React, { useState, useEffect } from 'react';
import { Modal, Space, Select } from 'antd';
import Button from "@/app/components/Button";
import { useGameState } from '../GameState';

const Settings = ({ isOpen, onClose }) => {
	const { config, setConfig } = useGameState();
	const [tileSet, setTileSet] = useState("beehive");

	const [allTileSets, setAllTileSets] = useState(["beehive", "numbers"]);

	useEffect(() => {
		if (config?.tileset) {
			setTileSet(config.tileset);
		}
	}, [config, setTileSet]);
	
	const handleSave = () => {
		setConfig({ ...config, tileset: tileSet });
		onClose();
	};

	return (
		<Modal
			title="Settings"
			open={isOpen}
			onCancel={onClose}
			footer={[
				<Button key="save" type="primary" onClick={handleSave}>
					Save
				</Button>
			]}
		>
			<Space direction="vertical" size="middle" style={{ width: '100%', padding: '20px 0' }}>
				<div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
					<span>tileset</span>
					<Select
						value={tileSet}
						onChange={setTileSet}
						style={{ width: 120 }}
						options={allTileSets.map((set) => ({
							value: set,
							label: set
						}))}
					/>
				</div>
			</Space>
		</Modal>
	);
};

export default Settings;
