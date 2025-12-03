import { createSupabaseClient } from "@/lib/supabase";
import { createFileRoute, redirect } from "@tanstack/react-router";
import { createServerFn } from "@tanstack/react-start";
import { getRequest } from "@tanstack/react-start/server";
import { useEffect } from "react";

const confirmFn = createServerFn({ method: "GET" })
  .inputValidator((searchParams: unknown) => {
    if (
      searchParams &&
      typeof searchParams === "object" &&
      "code" in searchParams &&
      "next" in searchParams
    ) {
      return searchParams;
    }
    throw new Error("Invalid search params");
  })
  .handler(async (ctx) => {
    const request = getRequest();

    if (!request) {
      return { success: false, error: "No request", redirectTo: null };
    }

    const searchParams = ctx.data;
    const code = searchParams["code"] as string;
    const _next = (searchParams["next"] ?? "/") as string;
    const next = _next?.startsWith("/") ? _next : "/";

    if (!code) {
      return { success: false, error: "No code found", redirectTo: null };
    }

    const supabase = await createSupabaseClient();

    const { error } = await supabase.auth.exchangeCodeForSession(code);

    if (error) {
      return { success: false, error: error.message, redirectTo: null };
    }

    return { success: true, error: null, redirectTo: next };
  });

export const Route = createFileRoute("/auth/oauth")({
  preload: false,
  staticData: {
    title: "OAuth Callback",
  },
  loader: async (opts) => {
    const result = await confirmFn({ data: opts.location.search });

    if (!result.success) {
      throw redirect({
        to: "/auth/error",
        search: { error: result.error || "Unknown error" },
      });
    }

    return result;
  },
  component: OAuthCallback,
});

function OAuthCallback() {
  const { redirectTo } = Route.useLoaderData();
  const navigate = Route.useNavigate();

  useEffect(() => {
    if (redirectTo) {
      navigate({ to: redirectTo, replace: true });
    }
  }, [redirectTo]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4" />
        <p>Completing sign in...</p>
      </div>
    </div>
  );
}
