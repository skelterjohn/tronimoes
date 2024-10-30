import localFont from "next/font/local";
import "./globals.css";
import 'antd/dist/reset.css';
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
