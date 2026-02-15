import { Modal, Input, Spin } from 'antd';
import { useState, useEffect, useRef } from 'react';

const KLIPY_API_KEY = 'lJodCdaswwEpTPg5Lhix66aIcaFsXBTrKKGlAzX1rPSQvagQHDNczRi42lNJ6x56';
const KLIPY_CLIENT_KEY = 'tronimoes-js';

function VisionQuest({ title = "Vision Quest", isOpen, onClose, setURL }) {
	const [path, setPath] = useState("");
	const [gifs, setGifs] = useState([]);
	const [loading, setLoading] = useState(false);

	const inputRef = useRef(null);

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
			afterOpenChange={(open) => {
				if (open) {
					inputRef.current?.focus();
				}
			}}
			onCancel={onClose}
			footer={null}
			centered
			width={800}
			className="vision-quest-modal"
			styles={{
				header: {
					backgroundColor: 'transparent',
					marginBottom: 0,
					paddingBottom: '16px'
				},
				content: {
					backgroundColor: '#f5f5f5'
				}
			}}
		>
			<div className="flex items-center gap-3 mb-4">
			    <Input
			        ref={inputRef}
			        placeholder="Search KLIPY"
			        value={path}
			        onChange={(e) => setPath(e.target.value)}
			        className="w-[80%] min-w-0"
			    />
			    <div className="flex-shrink-0 relative h-8 flex items-center">
			        <img 
			            src="/klipy_powered.png"
			            alt="Powered by KLIPY"
			            className="max-h-full w-auto object-contain object-center block" 
			        />
			        {loading && (
			            <span className="absolute -top-1 -right-1">
			                <Spin />
			            </span>
			        )}
			    </div>
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
							className="w-full h-auto rounded-sm"
						/>
					</div>
				))}
			</div>
		</Modal>
	);
}

export default VisionQuest;
