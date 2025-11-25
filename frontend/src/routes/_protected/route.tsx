import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/animate-ui/components/radix/sidebar";
import AppSidebar from "@/components/app-sidebar";
import { CentrifugoProvider } from "@/lib/centrifugo/centrifugo-context";
import { fetchUser } from "@/lib/fetch-user-server-fn";
import { ProtectedRouteContext } from "@/types/route-context";
import {
  createFileRoute,
  Outlet,
  redirect,
  useMatches,
} from "@tanstack/react-router";

export const Route = createFileRoute("/_protected")({
  beforeLoad: async () => {
    const user = await fetchUser();
    console.log("Fetch user");

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
  const { context } = matches[matches.length - 1] as {
    context: ProtectedRouteContext;
  };

  return (
    <CentrifugoProvider>
      <SidebarProvider>
        <AppSidebar />
        <SidebarInset>
          <SidebarTrigger className="size-12" />
          <div className="container mx-auto py-6 px-6 md:p-10 max-w-4xl">
            {context.pageTitle && (
              <div className="mb-6">
                <h1 className="text-3xl font-bold">{context.pageTitle}</h1>
              </div>
            )}
            <Outlet />
          </div>
        </SidebarInset>
      </SidebarProvider>
    </CentrifugoProvider>
  );
}
