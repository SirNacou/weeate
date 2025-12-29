import { createRouter, ErrorComponent } from "@tanstack/react-router"
import { setupRouterSsrQueryIntegration } from "@tanstack/react-router-ssr-query"
import * as TanstackQuery from "./integrations/tanstack-query/root-provider"

// Import the generated route tree
import { routeTree } from "./routeTree.gen"

// Setup client interceptors (only runs on client)
export const getRouter = () => {
	const rqContext = TanstackQuery.getContext()

	const router = createRouter({
		routeTree,
		context: { ...rqContext, pageTitle: "Weeate" },
		defaultPreload: "intent",
		defaultErrorComponent: ({ error }: { error: Error }) => <ErrorComponent error={error} />,
		Wrap: (props: { children: React.ReactNode }) => {
			return (
				<TanstackQuery.Provider {...rqContext}>
					{props.children}
				</TanstackQuery.Provider>
			)
		},
	})

	setupRouterSsrQueryIntegration({
		router,
		queryClient: rqContext.queryClient,
	})

	return router
}

// This code is only for TypeScript
declare global {
	interface Window {
		__TANSTACK_QUERY_CLIENT__: import("@tanstack/query-core").QueryClient
	}
}

if (typeof window !== "undefined") {
	// This code is for all users
	window.__TANSTACK_QUERY_CLIENT__ = TanstackQuery.getContext().queryClient
}
