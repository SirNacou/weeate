import { getAuthCentrifugoToken, serverClient } from "@/api"
import { createServerFn } from "@tanstack/react-start"

export const getCentrifugoTokenServerFn = createServerFn({
  method: "GET",
}).handler(async () => {
  const resp = await getAuthCentrifugoToken({ client: serverClient });
  if (resp.error) {
    console.error("Failed to get Centrifugo token:", resp.error);
    return {
      error:
        resp.error.errors?.at(0)?.message ?? "Failed to get Centrifugo token",
    };
  }
  return { data: resp.data };
});
