import { getRequestHeader } from "@tanstack/react-start/server";
import { createClient } from "@/client/client";
import { serverEnv } from "@/env";

const serverClient = createClient({
	baseUrl: serverEnv.BACKEND_URL,
	credentials: "include",
});

serverClient.interceptors.request.use((request: Request): Request => {
	request.headers.set("Cookie", getRequestHeader("Cookie") || "");
	return request;
});

export { serverClient };
