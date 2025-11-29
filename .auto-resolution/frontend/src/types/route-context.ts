/**
 * Extended context for protected routes
 */
export interface ProtectedRouteContext {
  pageTitle?: string;
  user: {
    id: string;
    display_name?: string;
    avatar_url?: string;
    created_at: string;
  };
}
