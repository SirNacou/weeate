import { createBrowserClient } from "@supabase/ssr";
import { clientEnv } from "@/env";

export function createClient() {
	return createBrowserClient(
		clientEnv.VITE_SUPABASE_URL,
		clientEnv.VITE_SUPABASE_PUBLISHABLE_OR_ANON_KEY,
	);
}
