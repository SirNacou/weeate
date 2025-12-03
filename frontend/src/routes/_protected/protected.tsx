import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useRef, useState } from "react";
import LucideUtensilsCrossed from "~icons/lucide/utensils-crossed?width=2em&height=2em";

export const Route = createFileRoute("/_protected/protected")({
	component: Info,
	staticData: {
		title: "Info",
	},
	loader: async ({ context }) => {
		return {
			user: context.user,
		};
	},
});

function Info() {
	return (
		<div>
			<MenuGrid />
		</div>
	);
}

// 1. The Smart Image Component
export const FoodImage = ({ src, alt }: { src: string; alt: string }) => {
	const [status, setStatus] = useState<"loading" | "error" | "success">(
		"loading",
	);
	const imgRef = useRef<HTMLImageElement>(null);

	// 1. THE FIX: Check if image is already loaded when component mounts
	useEffect(() => {
		const img = imgRef.current;
		if (img && img.complete) {
			// If complete and has width, it's a success. If width is 0, it failed.
			if (img.naturalWidth > 0) {
				setStatus("success");
			} else {
				setStatus("error");
			}
		}
	}, []); // Empty dependency array = runs once on mount

	return (
		<div className="relative w-full h-40 rounded-xl overflow-hidden bg-gray-100">
			{/* Loading State */}
			{status === "loading" && (
				<div className="absolute inset-0 bg-gray-200 animate-pulse flex items-center justify-center z-10">
					<span className="text-gray-400 text-xs font-medium">Loading...</span>
				</div>
			)}

			{/* Error State */}
			{status === "error" && (
				<div className="absolute inset-0 bg-gray-100 flex flex-col items-center justify-center text-gray-300 z-10">
					<LucideUtensilsCrossed className="size-20" />
					<span className="font-medium text-gray-400 uppercase">No Image</span>
				</div>
			)}

			{/* Actual Image */}
			<img
				ref={imgRef} // Attach the ref here
				src={src}
				alt={alt}
				className={`w-full h-full object-cover transition-opacity duration-500 ${
					status === "success" ? "opacity-100" : "opacity-0"
				}`}
				onLoad={() => setStatus("success")}
				onError={() => setStatus("error")}
			/>
		</div>
	);
};

// 2. The Demo Container
function MenuGrid() {
	return (
		<div className="p-6 max-w-4xl mx-auto">
			<h2 className="text-xl font-bold mb-6">Menu Preview</h2>

			<div className="grid grid-cols-1 md:grid-cols-3 gap-6">
				{/* Card 1: Success State */}
				<div className="space-y-2">
					<FoodImage
						src="https://images.unsplash.com/photo-1568901346375-23c9450c58cd?auto=format&fit=crop&w=500&q=60"
						alt="Gourmet Burger"
					/>
					<h3 className="font-semibold">Classic Burger</h3>
					<p className="text-sm text-gray-500">$12.00 • American</p>
				</div>

				{/* Card 2: Error State (Broken Link) */}
				<div className="space-y-2">
					<FoodImage
						src="https://broken-link-example.com/tacos.jpg"
						alt="Spicy Tacos"
					/>
					<h3 className="font-semibold">Street Tacos</h3>
					<p className="text-sm text-gray-500">$9.50 • Mexican</p>
				</div>

				{/* Card 3: Loading State (Simulated) */}
				<div className="space-y-2">
					{/* Manually forcing loading state for demo purposes */}
					<div className="w-full h-40 bg-gray-200 animate-pulse rounded-xl" />
					<h3 className="font-semibold">Pasta Alfredo</h3>
					<p className="text-sm text-gray-500">$14.00 • Italian</p>
				</div>
			</div>
		</div>
	);
}
