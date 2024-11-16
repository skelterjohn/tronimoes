import { Modal } from "antd";
import { signInAnonymously, GoogleAuthProvider, FacebookAuthProvider, signInWithPopup, fetchSignInMethodsForEmail, signInWithCredential, linkWithCredential } from "firebase/auth";
import { Button } from "antd";
import { auth } from "@/config";


export default function SignIn({setErrorMessage, setUserInfo, isOpen, onClose}) {
	const signInAsGuest = async () => {
		try {
			const result = await signInAnonymously(auth);
			console.log(result);
			setUserInfo(result.user);
			onClose();
		} catch (error) {
			console.error("Error signing in anonymously:", error);
			setErrorMessage("Error signing in anonymously");
			onClose();
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
			setErrorMessage("Error signing in with Google");
			onClose();
		}
	};

	const signInWithFacebook = async () => {
		try {
			const provider = new FacebookAuthProvider();
			const result = await signInWithPopup(auth, provider);
			setUserInfo(result.user);
			onClose();
		} catch (error) {
			if (error.code === 'auth/account-exists-with-different-credential') {
				setUserInfo({
					accessToken: error.customData._tokenResponse.oauthAccessToken,
					uid: error.customData._tokenResponse.localId,
				});
				onClose();
			} else {
				console.error("Error signing in with Facebook:", error);
				setErrorMessage("Error signing in with Facebook");
				onClose();
			}
		}
	};

	return (
		<Modal open={isOpen} title="sign in to tronimoes" onCancel={onClose} footer={null} centered width={800}>
			<div className="flex flex-col items-center gap-4 p-4">
				<Button onClick={signInWithGoogle} size="large">
					Sign in with Google
				</Button>
				<Button onClick={signInWithFacebook} size="large">
					Sign in with Facebook
				</Button>
				<Button onClick={signInAsGuest} size="large">
					Anonymous
				</Button>
			</div>
		</Modal>
	);
}
