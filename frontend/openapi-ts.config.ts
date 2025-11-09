import { defineConfig } from "@hey-api/openapi-ts";

export default defineConfig({
  input: {
    path: "../openapi.yaml",
    watch: true,
  },
  output: {
    path: "./src/client",
    format: "biome",
    lint: "biome",
  },
  plugins: [
    "@hey-api/typescript",
    "@hey-api/sdk",
    "@hey-api/transformers",
    "@hey-api/schemas",
    {
      name: "@hey-api/client-axios",
      runtimeConfigPath: "../api/api-client.ts",
    },
    "zod",
    "@tanstack/react-query",
  ],
});
