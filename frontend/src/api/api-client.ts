import { env } from "@/env";
import { type CreateClientConfig } from "@/client/client.gen";

const BASE_URL = env.VITE_BACKEND_URL || "http://localhost:8080/api";

export const createClientConfig: CreateClientConfig = (config) => ({
  ...config,
  baseURL: BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true,
});
