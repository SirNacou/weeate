import { useEffect, useRef } from "react";
import { z } from "zod";
import { useCentrifugo } from "./centrifugo-context";
import type { Channel } from "./channel";

export const useStrictSubscription = <T extends z.ZodSchema>(
	channel: Channel,
	schema: T,
	onData: (data: z.infer<T>) => void,
) => {
	const { subscribe, isConnected } = useCentrifugo();

	// Keep ref to latest callback
	const onDataRef = useRef(onData);
	useEffect(() => {
		onDataRef.current = onData;
	}, [onData]);

	useEffect(() => {
		if (!isConnected || !channel) return;

		// The unsafe callback receives "any" from Centrifugo
		const unsafeCallback = (rawJson: any) => {
			// 1. RUNTIME VALIDATION
			const result = schema.safeParse(rawJson);

			if (result.success) {
				// 2. SUCCESS: Data is safe, pass to component
				if (onDataRef.current) {
					onDataRef.current(result.data);
				}
			} else {
				// 3. FAILURE: Log error, do NOT crash the UI
				console.error(
					`[Centrifugo] Schema validation failed for channel "${channel}":`,
					z.treeifyError(result.error),
				);
			}
		};

		const unsubscribe = subscribe(channel, unsafeCallback);

		return () => {
			unsubscribe();
		};
	}, [channel, isConnected, subscribe, schema]);
};
