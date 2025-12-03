import { client } from "./api-client"
import { getBaseUrl } from "./hey-api"

client.setConfig({
  baseUrl: getBaseUrl(),
});

export { client } from "./api-client"
export { serverClient } from "./api-server-client"

export * from "@/client/sdk.gen"
export * from "@/client/types.gen"
