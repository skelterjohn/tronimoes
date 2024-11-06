export default function ChickenFoot({ color }) {
	const colorMap = {
        red: "bg-red-500",
        blue: "bg-blue-500",
        green: "bg-green-500",
        indigo: "bg-indigo-500",
        orange: "bg-orange-500",
        yellow: "bg-yellow-500",
        purple: "bg-purple-500",
        transparent: "bg-transparent"
    };
	return (
        <div className="w-full h-full flex items-center justify-center">
            <div className={`w-3/4 h-3/4 rounded-full ${colorMap[color]} opacity-50`}></div>
        </div>
    )
}
