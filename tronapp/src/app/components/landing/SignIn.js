import { Modal } from "antd";
import { signInAnonymously } from "firebase/auth";
import { Button } from "antd";

import { initializeApp } from 'firebase/app';
import { getAuth } from 'firebase/auth';

const firebaseConfig = {
	apiKey: "AIzaSyBQnYBVhGiSeoJaHFBdKyF5P6syZ7LFPM0",
	authDomain: "tronimoes.firebaseapp.com",
	projectId: "tronimoes",
	storageBucket: "tronimoes.firebasestorage.app",
	messagingSenderId: "1010961884428",
	appId: "1:1010961884428:web:fbe232e4ea189e9ac41529",
	measurementId: "G-75HZE3JJC7"
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);

// Initialize Firebase Authentication and get a reference to the service
export const auth = getAuth(app);

export default function SignIn({userInfo, setUserInfo, isOpen, onClose}) {
	const signInAsGuest = async () => {
		try {
			const result = await signInAnonymously(auth);
			console.log(result);
			setUserInfo(result.user);
			onClose();
		} catch (error) {
			console.error("Error signing in anonymously:", error);
		}
	};

	return (
		<Modal open={isOpen} title="sign in to tronimoes" onCancel={onClose} footer={null} centered width={800}>
			<div className="flex flex-col items-center gap-4 p-4">
				<Button onClick={signInAsGuest} size="large">
					Browser Session
				</Button>
			</div>
		</Modal>
	);
}
