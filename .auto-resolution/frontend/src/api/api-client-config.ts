import { clientEnv } from "@/env";
import { type CreateClientConfig } from "@/client/client.gen";

const BASE_URL = clientEnv.VITE_BACKEND_URL || "http://localhost:8080/api";

export const createClientConfig: CreateClientConfig = (config) => ({
  ...config,
  baseURL: BASE_URL,
  credentials: "include",
});
