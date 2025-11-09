import { client } from "@/client/client.gen";
import { createClient } from "@/lib/server";
import { type AxiosError, type InternalAxiosRequestConfig } from "axios";

if (typeof window !== "undefined") {
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
    (error: AxiosError) => {
      if (error.response) {
        const { status } = error.response;
        if (status === 401) {
          const supabase = createClient();
          supabase.auth.signOut();
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
}

export { client };
