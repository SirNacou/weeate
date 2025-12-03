import { defineConfig } from "@hey-api/openapi-ts"

export default defineConfig({
  input: {
    path: "http://localhost:8080/openapi.yaml",
    watch: true,
  },
  interactive: true,
  output: {
    clean: true,
    path: "./src/client",
    format: "biome",
    lint: "biome",
  },
  logs: {
    path: "/openapi-logs",
  },
  plugins: [
    "@hey-api/typescript",
    {
      name: "@hey-api/sdk",
      validator: true,
      transformer: true,
    },
    {
      name: "@hey-api/transformers",
      bigInt: false,
      dates: false,
    },
    "@hey-api/schemas",
    {
      name: "@hey-api/client-fetch",
      runtimeConfigPath: "@/api/hey-api.ts",
      dates: {
        offset: true,
      },
    },
    {
      name: "zod",
      dates: { offset: true },
    },
    "@tanstack/react-query",
  ],
});
