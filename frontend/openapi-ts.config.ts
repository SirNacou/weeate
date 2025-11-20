import { defineConfig } from "@hey-api/openapi-ts";

export default defineConfig({
  input: {
    path: "http://localhost:8080/openapi.yaml",
    watch: false,
  },
  interactive: true,
  output: {
    clean: true,
    path: "./src/client",
    format: "biome",
    lint: "biome",
  },
  plugins: [
    "@hey-api/typescript",
    { name: "@hey-api/sdk", validator: true },
    { name: "@hey-api/transformers", bigInt: false },
    "@hey-api/schemas",
    {
      name: "@hey-api/client-fetch",
      runtimeConfigPath: "../api/api-client-config.ts",
    },
    { name: "zod", requests: true, responses: true },
    "@tanstack/react-query",
  ],
});
