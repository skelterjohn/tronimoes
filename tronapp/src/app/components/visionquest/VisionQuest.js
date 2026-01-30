import { Modal, Input, Spin } from 'antd';
import { useState, useEffect, useRef } from 'react';

const KLIPY_API_KEY = 'lJodCdaswwEpTPg5Lhix66aIcaFsXBTrKKGlAzX1rPSQvagQHDNczRi42lNJ6x56';
const KLIPY_CLIENT_KEY = 'tronimoes-js';

function VisionQuest({ title = "Vision Quest", isOpen, onClose, setURL }) {
	const [path, setPath] = useState("");
	const [gifs, setGifs] = useState([]);
	const [loading, setLoading] = useState(false);

	const inputRef = useRef(null);

	useEffect(() => {
		// Add a small delay to ensure DOM is ready and other focus events have completed
		const timer = setTimeout(() => {
			if (inputRef.current && document.activeElement !== inputRef.current) {
				inputRef.current.focus();
			}
		}, 100);
		
		return () => clearTimeout(timer);
	}, []);

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
				`https://api.klipy.com/v2/search?q=${encodeURIComponent(searchTerm)}&key=${KLIPY_API_KEY}&client_key=${KLIPY_CLIENT_KEY}&limit=9&searchfilter=sticker`
			);
			const data = await response.json();
			setGifs(data.results);
		} catch (error) {
			console.error('Error fetching stickers:', error);
		} finally {
			setLoading(false);
		}
	};

	return (
		<Modal
			title={title}
			
			open={isOpen}
			onCancel={onClose}
			footer={null}
			centered
			width={800}
			className="vision-quest-modal"
			styles={{
				header: {
					backgroundColor: '#f5f5f5', // Or use your specific grey, e.g., '#f5f5f5'
					marginBottom: 0,
					paddingBottom: '16px'
				},
				content: {
					backgroundColor: '#f5f5f5', // Ensures the body matches
				}
			}}
		>
			<div className="flex items-center justify-between gap-2 mb-4">
			    <Input
			        inputRef={inputRef}
			        placeholder="Search KLIPY"
			        value={path}
			        onChange={(e) => setPath(e.target.value)}
			        className="w-4/5" // Keep your 80% width
			    />
			    <div className="flex flex-1 justify-end">
			        <img 
			            src="/klipy_powered.png"
			            alt="Powered by KLIPY"
			            className="max-h-10 w-auto object-contain" 
			        />
			    </div>
			    
			    {loading && (
			        <div className="absolute right-0 flex items-center pr-2">
			            <Spin />
			        </div>
			    )}
			</div>

			<div className="grid grid-cols-3 gap-4">
				{gifs.map((gif) => (
					<div 
						key={gif.id} 
						className="cursor-pointer hover:opacity-80 transition-opacity border-2 border "
						onClick={() => {
							setURL(gif.media_formats.tinygif_transparent.url);
							onClose();
						}}
					>
						<img 
							src={gif.media_formats.tinygif_transparent.url}
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
