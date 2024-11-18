import 'antd/dist/reset.css';
import localFont from "next/font/local";
import "./globals.css";
import { GameProvider } from "./components/GameState";

const geistSans = localFont({
	src: "./fonts/GeistVF.woff",
	variable: "--font-geist-sans",
	weight: "100 900",
});
const geistMono = localFont({
	src: "./fonts/GeistMonoVF.woff",
	variable: "--font-geist-mono",
	weight: "100 900",
});

export const metadata = {
	title: "tronimoes",
	description: "tronimoes",
	viewport: 'width=device-width, initial-scale=1, viewport-fit=cover, user-scalable=no',
	'apple-mobile-web-app-capable': 'yes',
	'mobile-web-app-capable': 'yes',
};

export default function RootLayout({ children }) {
	return (
		<html lang="en">
			<body
			className={`${geistSans.variable} ${geistMono.variable} antialiased`}
			>
				<GameProvider>
					{children}
				</GameProvider>
			</body>
		</html>
	);
}
