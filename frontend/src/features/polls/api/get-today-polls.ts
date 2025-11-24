import { serverClient } from "@/api/api-server-client";
import { listPollsToday } from "@/client/sdk.gen";
import { createServerFn } from "@tanstack/react-start";

export const listPollsTodayServerFn = createServerFn({ method: "GET" }).handler(
  async ({ signal }) => {
    const { data, error } = await listPollsToday({
      client: serverClient,
      signal,
    });
    console.log("Today Polls Data:", data, "Error:", error);

    if (error) {
      throw error;
    }

    return data;
  }
);
