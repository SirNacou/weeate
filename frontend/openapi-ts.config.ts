import { defineConfig } from "@hey-api/openapi-ts";

export default defineConfig({
  input: {
    path: "http://localhost:8080/openapi.yaml",
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
    { name: "@hey-api/transformers", bigInt: false },
    "@hey-api/schemas",
    {
      name: "@hey-api/client-fetch",
      runtimeConfigPath: "../api/api-client-config.ts",
    },
    "zod",
    "@tanstack/react-query",
  ],
});
