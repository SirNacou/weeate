import { client } from "@/client/client.gen";
import { env } from "@/env";
import { createClient } from "@/lib/server";

client.interceptors.request.use((request) => {
  return request;
});

// Handle successful responses
client.interceptors.response.use((response, request, options) => {
  // Store original request info in case we need to retry
  if (response.status === 401) {
    (response as any)._originalRequest = request;
    (response as any)._originalOptions = options;
  }
  return response;
});

// Handle errors including 401 unauthorized
client.interceptors.error.use(async (error) => {
  const response = error as Response & {
    _originalRequest?: Request;
    _originalOptions?: any;
  };

  // Handle 401 Unauthorized - retry with token refresh
  if (response?.status === 401 && response._originalRequest) {
    const supabase = createClient();
    let retryCount = 0;
    const maxRetries = 3;

    while (retryCount < maxRetries) {
      const { error: refreshError } = await supabase.auth.refreshSession();

      if (!refreshError) {
        // Token refreshed successfully, retry the original request
        const originalRequest = response._originalRequest;
        const originalOptions = response._originalOptions;

        try {
          // Clone the original request
          const clonedRequest = originalRequest.clone();

          // Retry the request with the new token (cookies auto-updated by Supabase)
          const _fetch = originalOptions?.fetch || fetch;
          const retryResponse = await _fetch(clonedRequest);

          if (retryResponse.ok) {
            return retryResponse;
          }
        } catch (retryError) {
          console.error("Retry failed:", retryError);
        }
      }

      retryCount++;

      // Wait before retrying (exponential backoff: 1s, 2s, 3s)
      if (retryCount < maxRetries) {
        await new Promise((resolve) => setTimeout(resolve, 1000 * retryCount));
      }
    }

    // All retries failed, log out the user
    await supabase.auth.signOut();
    window.cookieStore?.delete(env.VITE_AUTH_COOKIE_NAME);
    console.error("Unauthorized Access - Logging out...");

    if (window.location.pathname !== "/login") {
      window.location.href = `/login?redirect=${encodeURIComponent(
        window.location.pathname
      )}`;
    }
  }

  // Handle 500 Server Error
  if (response?.status === 500) {
    console.error("Server Error - Please try again later.");
  }

  // Re-throw the error so it can be handled by the caller
  return error;
});

export { client };
