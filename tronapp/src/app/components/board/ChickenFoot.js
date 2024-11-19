import { useEffect } from "react";
export default function ChickenFoot({ url, color }) {
	const colorMap = {
        red: "bg-red-700",
        blue: "bg-blue-700",
        green: "bg-green-700",
        yellow: "bg-yellow-700",
        orange: "bg-orange-700",
        yellow: "bg-yellow-700",
        purple: "bg-purple-700",
        transparent: "bg-transparent"
    };

	return (
        <div className="w-full h-full flex items-center justify-center">
            <div className={`w-3/4 h-3/4 rounded-lg ${colorMap[color]} border border-black relative`}>
            </div>
            {url && <div className="w-3/4 h-3/4 absolute">
                <img src={url} alt="Chicken Foot" className="w-full h-full object-contain" />
            </div>}
        </div>
    )
}
