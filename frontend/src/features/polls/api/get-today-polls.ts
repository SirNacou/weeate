import { createServerFn } from "@tanstack/react-start";
import { serverClient } from "@/api/api-server-client";
import { listPollsToday } from "@/client/sdk.gen";

export const listPollsTodayServerFn = createServerFn({ method: "GET" }).handler(
	async ({ signal }) => {
		const { data, error } = await listPollsToday({
			client: serverClient,
			signal,
		});

		if (error) {
			return {
				error: error.errors?.at(0)?.message ?? "Failed to fetch today's polls",
			};
		}

		return { data };
	},
);
