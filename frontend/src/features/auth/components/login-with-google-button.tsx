import { Button } from "@/components/ui/button";
import { createClient } from "@/lib/client";
import { useState } from "react";
import { LogosGoogleIcon } from "@/components/icons";

type Props = {} & React.HTMLAttributes<HTMLDivElement>;

export function LoginWithGoogleButton({}: Props) {
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleSocialLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    const supabase = createClient();
    setIsLoading(true);
    setError(null);

    try {
      const { error } = await supabase.auth.signInWithOAuth({
        provider: "google",
        options: {
          redirectTo: `${window.location.origin}/auth/oauth?next=/protected`,
        },
      });

      if (error) throw error;
    } catch (error: unknown) {
      setError(error instanceof Error ? error.message : "An error occurred");
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSocialLogin}>
      <div className="flex flex-col gap-6">
        {error && <p className="text-sm text-destructive-500">{error}</p>}

        <Button type="submit" className="w-full" disabled={isLoading}>
          {isLoading ?
            "Logging in..."
          : <>
              <LogosGoogleIcon className="mr-2 h-5 w-5" />
              Google
            </>
          }
        </Button>
      </div>
    </form>
  );
}
