import { env } from "@/env";
import { createClient } from "@/lib/client";
import axios, { type AxiosError, type InternalAxiosRequestConfig } from "axios";

const BASE_URL = env.VITE_BACKEND_URL || "http://localhost:8080/api";

const apiClient = axios.create({
  baseURL: BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true,
});

apiClient.interceptors.request.use(
  async (config: InternalAxiosRequestConfig) => {
    const supabase = createClient();
    const { data, error } = await supabase.auth.getSession();
    if (data.session && !error) {
      try {
        config.headers["Authorization"] = `Bearer ${data.session.access_token}`;
      } catch (error) {
        console.error("Error parsing auth token:", error);
      }
    }

    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error: AxiosError) => {
    if (error.response) {
      const { status } = error.response;
      if (status === 401) {
        // Handle unauthorized access, e.g., redirect to login
        console.error("Unauthorized Access - Logging out...");

        // Clear all Supabase auth keys
        const authKeys = Object.keys(localStorage).filter(
          (key) => key.startsWith("sb-") && key.includes("-auth-token")
        );
        // authKeys.forEach((key) => localStorage.removeItem(key));

        // if (window.location.pathname !== "/login") {
        //   window.location.href = `/login?redirect=${encodeURIComponent(
        //     window.location.pathname
        //   )}`;
        // }
      }

      if (status === 500) {
        console.error("Server Error - Please try again later.");
      }
    }

    return Promise.reject(error);
  }
);

export default apiClient;
