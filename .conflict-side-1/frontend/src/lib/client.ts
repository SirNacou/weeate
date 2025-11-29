import { clientEnv } from "@/env";
import { createBrowserClient } from "@supabase/ssr";

export function createClient() {
  return createBrowserClient(
    clientEnv.VITE_SUPABASE_URL,
    clientEnv.VITE_SUPABASE_PUBLISHABLE_OR_ANON_KEY,
    {
      auth: {
        autoRefreshToken: true,
      },
    }
  );
}
