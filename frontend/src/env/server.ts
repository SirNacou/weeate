import { createEnv } from "@t3-oss/env-core";
import z from "zod";

export const env = createEnv({
  server: {
    SERVER_URL: z.url().optional(),
    BACKEND_URL: z.url(),
    IMAGE_KIT_API_KEY: z.string(),
  },
  runtimeEnv: process.env,
  emptyStringAsUndefined: true,
});
