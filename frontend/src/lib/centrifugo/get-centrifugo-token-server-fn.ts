import { serverClient } from "@/api";
import { getAuthCentrifugoToken } from "@/client";
import { createServerFn } from "@tanstack/react-start";

export const getCentrifugoTokenServerFn = createServerFn({
  method: "GET",
}).handler(async () => {
  const resp = await getAuthCentrifugoToken({ client: serverClient });
  if (resp.error) {
    console.error("Failed to get Centrifugo token:", resp.error);
    throw new Error("Failed to get Centrifugo token");
  }
  return resp.data;
});
