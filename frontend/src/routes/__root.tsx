import { ImageKitProvider } from "@imagekit/react";
import { TanStackDevtools } from "@tanstack/react-devtools";
import { FormDevtoolsPanel } from "@tanstack/react-form-devtools";
import type { QueryClient } from "@tanstack/react-query";
import { ReactQueryDevtoolsPanel } from "@tanstack/react-query-devtools";
import {
	ClientOnly,
	createRootRouteWithContext,
	HeadContent,
	Scripts,
	useRouterState,
} from "@tanstack/react-router";
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools";
import { MotionConfig } from "motion/react";
import { Toaster } from "@/components/ui/sonner";
import { env } from "@/env/client";
import useIsMobile from "@/hooks/use-is-mobile";
import appCss from "../styles.css?url";

export interface MyRouterContext {
	queryClient: QueryClient;
	pageTitle: string;
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
	staticData: {
		title: "Weeate",
	},
	head: () => {
		return {
			meta: [
				{
					charSet: "utf-8",
				},
				{
					name: "viewport",
					content: "width=device-width, initial-scale=1",
				},
				{
					title: "Weeate",
				},
				{
					property: "og:title",
					content: "Weeate",
				},
				{
					property: "og:description",
					content: "Weeate - Your personal recipe manager",
				},
				{
					property: "og:image",
					content: `/logo512.png?v=${Date.now()}`,
				},
				{
					property: "og:url",
					content: "https://weeate.nacou.uk",
				},
			],
			links: [
				{
					rel: "stylesheet",
					href: appCss,
				},
				{
					rel: "icon",
					href: `/logo64.png?v=${Date.now()}`,
				},
			],
		};
	},
	shellComponent: RootDocument,
});

function RootDocument({ children }: { children: React.ReactNode }) {
	const { matches } = useRouterState();
	const currentRoute = matches[matches.length - 1];
	const pageTitle = currentRoute?.staticData.title;

	const { isMobile } = useIsMobile();
	return (
		<html lang="en">
			<head>
				<HeadContent />
				<title>{`${pageTitle} | Weeate`}</title>
			</head>
			<body>
				<ImageKitProvider
					urlEndpoint={env.VITE_IMAGEKIT_PUBLIC_KEY}
					transformationPosition="query"
				>
					<MotionConfig reducedMotion="user">
						<div className="Root">
							<ClientOnly>
								<Toaster
									position={isMobile ? "top-center" : "top-right"}
									richColors={true}
									theme="light" // TODO: adapt to dark mode
								/>
							</ClientOnly>
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
							name: "Tanstack Query",
							render: <ReactQueryDevtoolsPanel />,
						},
						{
							name: "Tanstack Router",
							render: <TanStackRouterDevtoolsPanel />,
						},
						{
							name: "Tanstack Form",
							render: <FormDevtoolsPanel />,
						},
					]}
				/>
				<Scripts />
			</body>
		</html>
	);
}

declare module "@tanstack/react-router" {
	interface StaticDataRouteOption {
		title: string;
	}
}
