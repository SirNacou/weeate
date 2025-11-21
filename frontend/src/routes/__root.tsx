import {
  HeadContent,
  Scripts,
  createRootRouteWithContext,
} from "@tanstack/react-router";
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools";
import { TanStackDevtools } from "@tanstack/react-devtools";
import { FormDevtoolsPlugin } from "@tanstack/react-form-devtools";

import TanStackQueryDevtools from "../integrations/tanstack-query/devtools";

import appCss from "../styles.css?url";

import type { QueryClient } from "@tanstack/react-query";
import { ImageKitProvider } from "@imagekit/react";
import { MotionConfig } from "motion/react";
import { env } from "@/env/client";
import { Toaster } from "sonner";
import useIsMobile from "@/hooks/use-is-mobile";

interface MyRouterContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  head: () => ({
    meta: [
      {
        charSet: "utf-8",
      },
      {
        name: "viewport",
        content: "width=device-width, initial-scale=1",
      },
      {
        title: "TanStack Start Starter",
      },
    ],
    links: [
      {
        rel: "stylesheet",
        href: appCss,
      },
    ],
  }),
  shellComponent: RootDocument,
});

function RootDocument({ children }: { children: React.ReactNode }) {
  const { isMobile } = useIsMobile();
  return (
    <html lang="en">
      <head>
        <HeadContent />
      </head>
      <body>
        <ImageKitProvider
          urlEndpoint={env.VITE_IMAGEKIT_PUBLIC_KEY}
          transformationPosition="query"
        >
          <MotionConfig reducedMotion="user">
            <div className="Root">
              <Toaster
                position={isMobile ? "top-center" : "top-right"}
                richColors={true}
              />
              {children}
            </div>
          </MotionConfig>
        </ImageKitProvider>
        <TanStackDevtools
          config={{
            position: "bottom-right",
          }}
          plugins={[
            {
              name: "Tanstack Router",
              render: <TanStackRouterDevtoolsPanel />,
            },
            TanStackQueryDevtools,
            FormDevtoolsPlugin(),
          ]}
        />
        <Scripts />
      </body>
    </html>
  );
}
