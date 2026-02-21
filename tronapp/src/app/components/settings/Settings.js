import React, { useState, useEffect } from 'react';
import { Modal, Space, Select, Checkbox } from 'antd';
import Button from "@/app/components/Button";
import { useGameState } from '../GameState';

const Settings = ({ isOpen, onClose, setShowReportIssueModal }) => {
	const { config, setConfig } = useGameState();
	const [tileSet, setTileSet] = useState("classic");
	const [soundEffects, setSoundEffects] = useState(true);

	const [allTileSets, setAllTileSets] = useState([
		"beehive",
		"beehive-mono",
		"classic",
		"classic-mono",
		"numbers",
		"numbers-mono",
	]);

	useEffect(() => {
		if (config?.tileset) {
			setTileSet(config.tileset);
		}
		if (config?.soundEffects !== undefined) {
			setSoundEffects(config.soundEffects);
		}
	}, [config]);
	
	const handleSave = () => {
		setConfig({ ...config, tileset: tileSet, soundEffects });
		onClose();
	};

	const handleReportIssue = () => {
		onClose();
		setShowReportIssueModal(true);
	};

	return (
		<Modal
			title="settings"
			open={isOpen}
			onCancel={onClose}
			className="font-game"
			footer={
				<div style={{ display: 'flex', justifyContent: 'space-between', width: '100%' }}>
					<Button
						key="report"
						className="game-btn"
						onClick={handleReportIssue}
					>
						report issue
					</Button>
					<Button key="save" type="primary" onClick={handleSave}>
						save
					</Button>
				</div>
			}
		>
			<Space orientation="vertical" size="middle" style={{ width: '100%', padding: '20px 0' }}>
				<div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
					<span>tileset</span>
					<Select
						value={tileSet}
						onChange={setTileSet}
						style={{ width: 200 }}
						className="font-game"
						classNames={{ popup: { root: "font-game" } }}
						options={allTileSets.map((set) => ({
							value: set,
							label: set
						}))}
					/>
				</div>
				<div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
					<span>sound effects</span>
					<Checkbox
						checked={soundEffects}
						onChange={(e) => setSoundEffects(e.target.checked)}
						className="font-game"
					/>
				</div>
			</Space>
		</Modal>
	);
};

export default Settings;
