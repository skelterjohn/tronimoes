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
	appleWebApp: {
		capable: true,
		statusBarStyle: "black-translucent",
		title: "tronimoes",
	},
};

export const viewport = {
	themeColor: "#0a0a0a",
	width: "device-width",
	initialScale: 1,
	maximumScale: 1,
	userScalable: false,
	viewportFit: "cover",
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
