import { env } from "@/env";
import { createClient, type User } from "@supabase/supabase-js";
import { createContext, useContext, useEffect, useState } from "react";

export const supabase = createClient(
  env.VITE_SUPABASE_URL,
  env.VITE_SUPABASE_API_KEY
);

export interface SupabaseAuthState {
  isAuthenticated: boolean;
  user: any;
  login: (email: string, password: string) => Promise<void>;
  loginWithGoogle: (options?: {
    /** A URL to send the user to after they are confirmed. */
    redirectTo?: string;
    /** A space-separated list of scopes granted to the OAuth application. */
    scopes?: string;
    /** An object of query params */
    queryParams?: { [key: string]: string };
    /** If set to true does not immediately redirect the current browser context to visit the OAuth authorization page for the provider. */
    skipBrowserRedirect?: boolean;
  }) => Promise<void>;
  logout: () => Promise<void>;
  isLoading: boolean;
}

const SupabaseAuthContext = createContext<SupabaseAuthState | undefined>(
  undefined
);

export function SupabaseAuthProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [user, setUser] = useState<User | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  useEffect(() => {
    // Get initial session
    supabase.auth.getSession().then(({ data: { session } }) => {
      setUser(session?.user ?? null);
      setIsAuthenticated(!!session?.user);
      setIsLoading(false);
    });

    // Listen for auth changes
    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange((_event, session) => {
      setUser(session?.user ?? null);
      setIsAuthenticated(!!session?.user);
      setIsLoading(false);
    });

    return () => subscription.unsubscribe();
  }, []);

  const login = async (email: string, password: string) => {
    const { error } = await supabase.auth.signInWithPassword({
      email,
      password,
    });
    if (error) throw error;
  };

  const loginWithGoogle = async (options?: {
    /** A URL to send the user to after they are confirmed. */
    redirectTo?: string;
    /** A space-separated list of scopes granted to the OAuth application. */
    scopes?: string;
    /** An object of query params */
    queryParams?: { [key: string]: string };
    /** If set to true does not immediately redirect the current browser context to visit the OAuth authorization page for the provider. */
    skipBrowserRedirect?: boolean;
  }) => {
    const { error } = await supabase.auth.signInWithOAuth({
      provider: "google",
      options,
    });

    console.log((await supabase.auth.getUser()).data.user);
    if (error) throw error;
  };

  const logout = async () => {
    const { error } = await supabase.auth.signOut();
    if (error) throw error;
  };

  return (
    <SupabaseAuthContext.Provider
      value={{
        isAuthenticated,
        user,
        login,
        loginWithGoogle,
        logout,
        isLoading,
      }}
    >
      {children}
    </SupabaseAuthContext.Provider>
  );
}

export function useSupabaseAuth() {
  const context = useContext(SupabaseAuthContext);
  if (context === undefined) {
    throw new Error("useSupabaseAuth must be used within SupabaseAuthProvider");
  }
  return context;
}
