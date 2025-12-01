import { LogosGoogleIcon } from "@/components/icons";
import { Button } from "@/components/ui/button";
import { createClient } from "@/lib/client";
import { useState } from "react";

type Props = {
  redirectTo?: string;
} & React.HTMLAttributes<HTMLFormElement>;

export function LoginWithGoogleButton({ redirectTo = "/", ...props }: Props) {
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const supabase = createClient();

  const handleSocialLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);

    try {
      const { error } = await supabase.auth.signInWithOAuth({
        provider: "google",
        options: {
          redirectTo: `${window.location.origin}/auth/oauth?next=${redirectTo}`,
          queryParams: {
            access_type: "offline",
            prompt: "select_account",
          },
        },
      });

      if (error) throw error;
    } catch (error: unknown) {
      setError(error instanceof Error ? error.message : "An error occurred");
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSocialLogin} {...props}>
      <div className="flex flex-col gap-6">
        {error && (
          <div className="rounded-lg bg-destructive/10 p-3 text-sm text-destructive">
            {error}
          </div>
        )}

        <Button
          type="submit"
          className="w-full"
          variant="outline"
          disabled={isLoading}
        >
          {isLoading ?
            "Connecting to Google..."
          : <>
              <LogosGoogleIcon className="mr-2 h-5 w-5" />
              Continue with Google
            </>
          }
        </Button>
      </div>
    </form>
  );
}
