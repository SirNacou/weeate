import { client } from "@/client/client.gen";
import { env } from "@/env";
import { createClient } from "@/lib/server";
import { type AxiosError, type InternalAxiosRequestConfig } from "axios";

client.instance.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

client.instance.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error: AxiosError) => {
    if (error.response) {
      const { status } = error.response;
      if (status === 401) {
        const supabase = createClient();

        let retry = 0;
        while (retry < 3) {
          const { error: refreshError } = await supabase.auth.refreshSession();
          if (!refreshError) {
            return client.instance.request(
              error.config as InternalAxiosRequestConfig
            );
          }
          retry++;
        }

        supabase.auth.signOut();
        window.cookieStore.delete(env.VITE_AUTH_COOKIE_NAME);
        console.error("Unauthorized Access - Logging out...");

        if (window.location.pathname !== "/login") {
          window.location.href = `/login?redirect=${encodeURIComponent(
            window.location.pathname
          )}`;
        }
      }

      if (status === 500) {
        console.error("Server Error - Please try again later.");
      }
    }

    return Promise.reject(error);
  }
);

export { client };
