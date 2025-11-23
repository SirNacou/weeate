// Example: authMiddleware.ts
import { createClient } from "@/lib/server";
import { createMiddleware } from "@tanstack/react-start";

export const refreshSessionMiddleware = createMiddleware().server(
  async ({ next }) => {
    const supa = createClient();
    const {
      data: { user },
      error,
    } = await supa.auth.getUser(); // This will refresh the session if needed and set new cookies

    if (error) {
      // Handle error, e.g., redirect to login
      console.error("Error getting user in middleware:", error);
      // You might want to throw an error or redirect here
    }

    return next({
      context: {
        user, // Make the user object available in the context
      },
    });
  }
);
