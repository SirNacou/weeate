import { createClient } from "@/lib/server";
import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/auth/oauth")({
  validateSearch: (search: Record<string, unknown>) => {
    return {
      code: (search.code as string) || undefined,
      next: (search.next as string) || "/",
    };
  },
  beforeLoad: async ({ search }) => {
    const { code, next } = search;

    if (!code) {
      throw redirect({
        to: "/auth/auth-code-error",
        search: { error: "No code provided" },
      });
    }

    const supabase = createClient();
    const { error } = await supabase.auth.exchangeCodeForSession(code);

    if (error) {
      console.error("OAuth error:", error);
      throw redirect({
        to: "/auth/auth-code-error",
        search: { error: error.message },
      });
    }

    // Ensure next is a relative URL
    const safeNext = next.startsWith("/") ? next : "/";
    throw redirect({ to: safeNext });
  },
});
