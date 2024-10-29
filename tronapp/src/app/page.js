"use client";

import { useState } from 'react';
import Image from 'next/image';
import { Input, Button } from 'antd';

export default function Home() {

	const [name, setName] = useState("");
	const [code, setCode] = useState("");

	return (
		<div className="relative min-h-screen flex items-center justify-center">
			<Image 
			src="/trondude.png"
			alt="Background"
			fill
			className="object-cover z-0"
			priority
			/>

			<div className="relative z-10 flex flex-col gap-4 text-white">
				<Input
					placeholder="name"
					size="large"
					className="text-lg"
					value={name}
					onChange={(e) => setName(e.target.value)}
				/>
				<div className="flex gap-2">
					<Input
						placeholder="code"
						size="large"
						className="text-lg"
						value={code}
						onChange={(e) => setCode(e.target.value)}
					/>
					<Button
						type="primary"
						size="large"
					>
						join
					</Button>
				</div>
				<div className="flex w-full gap-2">
					<Button
						className="w-full"
						type="primary"
						size="large"
					>
						random
					</Button>
				</div>
			</div>
		</div>
	);
}
