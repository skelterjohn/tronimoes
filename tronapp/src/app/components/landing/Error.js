export default function Error({message}) {
	if (!message) return null;

	return (
		<div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 
						p-3 px-5 bg-white border-2 border-red-500 rounded-md text-red-500 
						z-50">
			{message}
		</div>
	)
}
