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
