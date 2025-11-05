import { env } from "@/env";
import axios, { type AxiosError, type InternalAxiosRequestConfig } from "axios";

const BASE_URL = env.VITE_BACKEND_URL || "http://localhost:8080";

const apiClient = axios.create({
  baseURL: BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Get Supabase session from localStorage
    const authKey = Object.keys(localStorage).find(
      (key) => key.startsWith("sb-") && key.endsWith("-auth-token")
    );

    if (authKey) {
      try {
        const authData = localStorage.getItem(authKey);
        if (authData) {
          const session = JSON.parse(authData);
          // Supabase stores the token in session.access_token
          const token = session?.access_token;

          if (token) {
            config.headers["Authorization"] = `Bearer ${token}`;
          }
        }
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
