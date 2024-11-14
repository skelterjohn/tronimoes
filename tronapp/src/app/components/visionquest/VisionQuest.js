import { Modal, Input, Spin } from 'antd';
import { useState, useEffect } from 'react';

// Replace with your actual Tenor API key
const TENOR_API_KEY = 'AIzaSyBPpZRb23wy2zTKQ2j5eJHS8YVPtjIvcGQ';
const TENOR_CLIENT_KEY = 'tronimoes'; // Replace with your app name

function VisionQuest({ isOpen, onClose, setChickenFootID }) {
	const [path, setPath] = useState("");
	const [gifs, setGifs] = useState([]);
	const [loading, setLoading] = useState(false);

	// Debounced search function
	useEffect(() => {
		const searchTimer = setTimeout(() => {
			if (path.trim()) {
				searchGifs(path);
			}
		}, 500);

		return () => clearTimeout(searchTimer);
	}, [path]);

	const searchGifs = async (searchTerm) => {
		setLoading(true);
		try {
			const response = await fetch(
				`https://tenor.googleapis.com/v2/search?q=${encodeURIComponent(searchTerm)}&key=${TENOR_API_KEY}&client_key=${TENOR_CLIENT_KEY}&limit=9`
			);
			const data = await response.json();
			console.log("data", data);
			setGifs(data.results);
		} catch (error) {
			console.error('Error fetching GIFs:', error);
		} finally {
			setLoading(false);
		}
	};

	return (
		<Modal
			title="Vision Quest"
			open={isOpen}
			onCancel={onClose}
			footer={null}
			centered
			width={800}
			className="vision-quest-modal"
		>
			<Input
				placeholder="enter your path"
				value={path}
				onChange={(e) => setPath(e.target.value)}
				className="mb-4"
			/>
			
			{loading && (
				<div className="flex justify-center my-4">
					<Spin />
				</div>
			)}

			<div className="grid grid-cols-3 gap-4">
				{gifs && gifs.map((gif) => (
					<div 
						key={gif.id} 
						className="cursor-pointer hover:opacity-80 transition-opacity"
						onClick={() => {
							setChickenFootID(gif.id);
							onClose();
						}}
					>
						<img 
							src={gif.media_formats.gif.url} 
							alt={gif.content_description}
							className="w-full h-auto rounded"
						/>
					</div>
				))}
			</div>
		</Modal>
	);
}

export default VisionQuest;
