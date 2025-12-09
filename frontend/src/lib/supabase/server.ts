import { createServerClient } from "@supabase/ssr";
import { getCookies, setCookie } from "@tanstack/react-start/server";
import { clientEnv } from "@/env";

export function createClient() {
	return createServerClient(
		clientEnv.VITE_SUPABASE_URL,
		clientEnv.VITE_SUPABASE_PUBLISHABLE_OR_ANON_KEY,
		{
			cookies: {
				getAll() {
					return Object.entries(getCookies()).map(
						([name, value]) =>
							({
								name,
								value,
							}) as { name: string; value: string },
					);
				},
				setAll(cookies) {
					cookies.forEach((cookie) => {
						const { name, value, options } = cookie;
						// Handle cookie deletion (maxAge: 0) and set proper cookie options
						setCookie(name, value, options);
					});
				},
			},
		},
	);
}
