import { Modal } from "antd";
import { signInAnonymously, GoogleAuthProvider, signInWithPopup } from "firebase/auth";
import { Button } from "antd";
import { auth } from "@/config";


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

	const signInWithGoogle = async () => {
		try {
			const provider = new GoogleAuthProvider();
			const result = await signInWithPopup(auth, provider);
			setUserInfo(result.user);
			onClose();
		} catch (error) {
			console.error("Error signing in with Google:", error);
		}
	};

	return (
		<Modal open={isOpen} title="sign in to tronimoes" onCancel={onClose} footer={null} centered width={800}>
			<div className="flex flex-col items-center gap-4 p-4">
				<Button onClick={signInWithGoogle} size="large">
					Sign in with Google
				</Button>
				<Button onClick={signInAsGuest} size="large">
					Anonymous
				</Button>
			</div>
		</Modal>
	);
}
