import React, { useState } from 'react';
import { Modal, Space, Input } from 'antd';
import Button from "@/app/components/Button";
import { useGameState } from '../GameState';

const ReportIssue = ({ isOpen, onClose, code, playErrorMessage }) => {
	const { client } = useGameState();
	const [summary, setSummary] = useState('');
	const [whatHappened, setWhatHappened] = useState('');
	const [whatShouldHappen, setWhatShouldHappen] = useState('');
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState(null);

	const handleReport = async () => {
		if (!client || !code) return;
		setSubmitting(true);
		setError(null);
		try {
			await client.ReportIssue(code, summary, whatHappened, whatShouldHappen, playErrorMessage);
			setSummary('');
			setWhatHappened('');
			setWhatShouldHappen('');
			setError(null);
			onClose();
		} catch (err) {
			setError(err?.message || err?.data?.error || 'Failed to report');
		} finally {
			setSubmitting(false);
		}
	};

	return (
		<Modal
			title={`report issue for ${code}`}
			open={isOpen}
			onCancel={onClose}
			className="font-game"
			footer={
				<div style={{ display: 'flex', justifyContent: 'flex-end', width: '100%' }}>
					{error && <span style={{ color: 'var(--ant-color-error)', marginRight: 12 }}>{error}</span>}
					<Button key="save" type="primary" onClick={handleReport} loading={submitting} disabled={submitting}>
						report
					</Button>
				</div>
			}
		>
			<Space orientation="vertical" size="middle" style={{ width: '100%', padding: '20px 0' }}>
				<div>
					<div style={{ marginBottom: 8 }}>current error message</div>
					<div className="text-red-500">{playErrorMessage || "<none>"}</div>
				</div>
				<div>
					<div style={{ marginBottom: 8 }}>summary</div>
					<Input
						value={summary}
						onChange={(e) => setSummary(e.target.value)}
						className="font-game"
						style={{ width: '100%' }}
					/>
				</div>
				<div>
					<div style={{ marginBottom: 8 }}>what happened?</div>
					<Input.TextArea
						value={whatHappened}
						onChange={(e) => setWhatHappened(e.target.value)}
						className="font-game"
						style={{ width: '100%' }}
						rows={4}
						autoSize={{ minRows: 4 }}
					/>
				</div>
				<div>
					<div style={{ marginBottom: 8 }}>what should happen instead?</div>
					<Input.TextArea
						value={whatShouldHappen}
						onChange={(e) => setWhatShouldHappen(e.target.value)}
						className="font-game"
						style={{ width: '100%' }}
						rows={4}
						autoSize={{ minRows: 4 }}
					/>
				</div>
			</Space>
		</Modal>
	);
};

export default ReportIssue;
