// routes/health.ts
import { createFileRoute } from "@tanstack/react-router";
import { json } from "@tanstack/react-start";

export const Route = createFileRoute("/health")({
  staticData: {
    title: "Health Check",
  },
  server: {
    handlers: {
      GET: async () => {
        const checks = {
          status: "healthy",
          timestamp: new Date().toISOString(),
          uptime: process.uptime(),
          memory: process.memoryUsage(),
          // database: await checkDatabase(),
          version: process.env.npm_package_version,
        };

        return json(checks);
      },
    },
  },
});

// async function checkDatabase() {
//   try {
//     await db.raw('SELECT 1')
//     return { status: 'connected', latency: 0 }
//   } catch (error) {
//     return { status: 'error', error: error.message }
//   }
// }
