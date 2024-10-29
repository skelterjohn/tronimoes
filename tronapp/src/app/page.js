"use client";

import Image from 'next/image';
import { Input, Button } from 'antd';

export default function Home() {
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
				/>
				<div className="flex gap-2">
					<Input
						placeholder="code"
						size="large"
						className="text-lg"
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
