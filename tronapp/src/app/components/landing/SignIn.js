import { Modal } from "antd";
import { signInAnonymously } from "firebase/auth";
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
