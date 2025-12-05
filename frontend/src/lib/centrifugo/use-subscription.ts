import { useEffect, useRef } from "react";
import { useCentrifugo } from "./centrifugo-context";
import type { Channel } from "./channel";

// Generic T allows you to type the incoming data
export const useSubscription = <T = any>(
	channel: Channel,
	onData: (data: T) => void,
) => {
	const { subscribe, isConnected } = useCentrifugo();

	// Keep a ref to the latest callback to avoid re-subscribing on render
	const onDataRef = useRef(onData);

	useEffect(() => {
		onDataRef.current = onData;
	}, [onData]);

	useEffect(() => {
		if (!isConnected || !channel) return;

		const safeCallback = (data: any) => {
			if (onDataRef.current) {
				onDataRef.current(data as T);
			}
		};

		const unsubscribe = subscribe(channel, safeCallback);

		return () => {
			unsubscribe();
		};
	}, [channel, isConnected, subscribe]);
};
