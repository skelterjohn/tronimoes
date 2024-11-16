import { Modal } from "antd";

export default function SignIn({userInfo, setUserInfo, isOpen, onClose}) {
	return (
		<Modal open={isOpen} title="sign in to tronimoes" onCancel={onClose} footer={null} centered width={800}>
			<div>Sign In</div>
		</Modal>
	);
}
