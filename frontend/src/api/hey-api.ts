import type { CreateClientConfig } from "@/client/client.gen"
import { clientEnv, serverEnv } from "@/env"
import { createIsomorphicFn } from "@tanstack/react-start"

export const getBaseUrl = createIsomorphicFn()
.server(() => serverEnv.BACKEND_URL)
.client(() => clientEnv.VITE_BACKEND_URL);

export const createClientConfig: CreateClientConfig = (config) => {
  return ({
    ...config,
    baseURL: getBaseUrl(),
  })
};
