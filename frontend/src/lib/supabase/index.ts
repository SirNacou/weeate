import { createIsomorphicFn } from "@tanstack/react-start";

export const createSupabaseClient = createIsomorphicFn()
	.server(async () => {
		const server = await import("@/lib/supabase/server");
		return server.createClient();
	})
	.client(async () => {
		const client = await import("@/lib/supabase/client");
		return client.createClient();
	});
