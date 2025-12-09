import {
	createFileRoute,
	Outlet,
	redirect,
	useMatches,
} from "@tanstack/react-router";
import {
	SidebarInset,
	SidebarProvider,
	SidebarTrigger,
} from "@/components/animate-ui/components/radix/sidebar";
import AppSidebar from "@/components/app-sidebar";
import { fetchUser } from "@/lib/supabase/fetch-user-server-fn";

export const Route = createFileRoute("/_protected")({
	staticData: {
		title: "",
	},
	beforeLoad: async () => {
		const user = await fetchUser();

		if (!user) {
			throw redirect({ to: "/login" });
		}

		return {
			user,
		};
	},
	component: ProtectedLayout,
});

function ProtectedLayout() {
	const matches = useMatches();
	const staticData = matches[matches.length - 1].staticData;

	return (
		<SidebarProvider>
			<AppSidebar />
			<SidebarInset>
				<SidebarTrigger className="size-12" />
				<div className="container mx-auto py-6 px-6 md:p-10 max-w-4xl">
					{staticData.title && (
						<div className="mb-6">
							<h1 className="text-3xl font-bold">{staticData.title}</h1>
						</div>
					)}
					<Outlet />
				</div>
			</SidebarInset>
		</SidebarProvider>
	);
}
