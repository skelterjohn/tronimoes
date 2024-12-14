import { Modal } from "antd";
import { signInAnonymously, GoogleAuthProvider, FacebookAuthProvider, signInWithPopup, linkWithCredential } from "firebase/auth";
import { Button } from "antd";
import { auth } from "@/config";

export default function SignIn({setErrorMessage, setUserInfo, isOpen, onClose}) {
	const signInAsGuest = async () => {
		try {
			const result = await signInAnonymously(auth);
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
				try {
					// 1. Get the email from the error
					const email = error.customData.email;
					
					const methods = error.customData._tokenResponse.verifiedProvider;
					
					// 3. If the user has previously signed in with Google
					if (methods.includes('google.com')) {
						// Sign in with Google
						const googleProvider = new GoogleAuthProvider();
						googleProvider.setCustomParameters({ login_hint: email });
						const googleResult = await signInWithPopup(auth, googleProvider);
						
						// Get Facebook credential
						const fbCredential = FacebookAuthProvider.credentialFromError(error);
						
						// Link the Facebook credential to the Google account
						await linkWithCredential(googleResult.user, fbCredential);
						
						// Update user info and set up token refresh
						setUserInfo(googleResult.user);
						onClose();
					} else {
						setErrorMessage(`Please sign in with ${methods[0]}`);
					}
				} catch (linkError) {
					console.error("Error linking accounts:", linkError);
					setErrorMessage("Error linking accounts");
				}
			} else {
				console.error("Error signing in with Facebook:", error);
				setErrorMessage("Error signing in with Facebook");
			}
			onClose();
		}
	};

	return (
		<Modal open={isOpen} title="sign in to tronimoes" onCancel={onClose} footer={null} centered width={800}>
			<div className="flex">
				<div className="hidden md:block w-1/2">
					<img 
						src="/fallingtiles.png" 
						alt="Falling Tiles" 
						className="object-cover w-full h-full"
					/>
				</div>
				<div className="w-full flex flex-col items-center gap-4 p-4">
					<Button onClick={signInWithGoogle} size="large" className="w-48">
						Sign in with Google
					</Button>
					<Button onClick={signInWithFacebook} size="large" className="w-48">
						Sign in with Facebook
					</Button>
					<Button onClick={signInAsGuest} size="large" className="w-48">
						Anonymous
					</Button>
				</div>
			</div>
		</Modal>
	);
}
