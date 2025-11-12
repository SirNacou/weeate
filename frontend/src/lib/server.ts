import { clientEnv } from "@/env";
import { createServerClient } from "@supabase/ssr";
import { getCookies, setCookie } from "@tanstack/react-start/server";

export function createClient() {
  return createServerClient(
    clientEnv.VITE_SUPABASE_URL,
    clientEnv.VITE_SUPABASE_PUBLISHABLE_OR_ANON_KEY,
    {
      auth: {
        autoRefreshToken: true,
      },
      cookies: {
        getAll() {
          return Object.entries(getCookies()).map(
            ([name, value]) =>
              ({
                name,
                value,
              }) as { name: string; value: string }
          );
        },
        setAll(cookies) {
          cookies.forEach((cookie) => {
            setCookie(cookie.name, cookie.value);
          });
        },
      },
    }
  );
}
