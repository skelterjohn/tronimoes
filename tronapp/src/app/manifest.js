export default function manifest() {
	return {
		name: "tronimoes",
		short_name: "tronimoes",
		description: "tronimoes – domino board game",
		start_url: "/",
		display: "standalone",
		background_color: "#0a0a0a",
		theme_color: "#0a0a0a",
		orientation: "any",
		icons: [
			{
				src: "/icon-192.png",
				sizes: "192x192",
				type: "image/png",
				purpose: "any",
			},
			{
				src: "/icon-512.png",
				sizes: "512x512",
				type: "image/png",
				purpose: "any",
			},
			{
				src: "/icon-512.png",
				sizes: "512x512",
				type: "image/png",
				purpose: "maskable",
			},
		],
	};
}
