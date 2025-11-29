import { createIsomorphicFn, createMiddleware } from "@tanstack/react-start";
import { getCookies, setCookie } from "@tanstack/react-start/server";

export const setCookieMiddleware = createMiddleware({
  type: "function",
}).client(({ next }) => {
  const cookies = getCookies();
  return next({
    headers: {
      // Cookie: cookies.,
    },
  });
});
