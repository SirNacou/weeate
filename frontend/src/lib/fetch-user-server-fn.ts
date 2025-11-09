import { createClient } from "@/lib/server";
import type { Factor, User } from "@supabase/supabase-js";
import { createServerFn } from "@tanstack/react-start";
type SSRSafeUser = User & {
  factors: (Factor & { factor_type: "phone" | "totp" })[];
  app_metadata: {
    provider: string | null;
    avatar_url: string;
    display_name: string;
  };
};

export const fetchUser: () => Promise<SSRSafeUser | null> = createServerFn({
  method: "GET",
}).handler(async () => {
  const supabase = createClient();
  const { data, error } = await supabase.auth.getUser();
  console.log(data);

  if (error) {
    return null;
  }

  return data.user as SSRSafeUser | null;
});
