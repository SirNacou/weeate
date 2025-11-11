import { env } from "@/env";
import { createBrowserClient } from "@supabase/ssr";

export function createClient() {
  return createBrowserClient(
    env.VITE_SUPABASE_URL,
    env.VITE_SUPABASE_PUBLISHABLE_OR_ANON_KEY,
    {
      auth: {
        autoRefreshToken: true,
      },
    }
  );
}
