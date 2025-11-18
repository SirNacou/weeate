import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/animate-ui/components/radix/sidebar";
import AppSidebar from "@/components/app-sidebar";
import { fetchUser } from "@/lib/fetch-user-server-fn";
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected")({
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
  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <SidebarTrigger className="size-12" />
        <div className="container mx-auto p-10 xxl:px-0">
          <Outlet />
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
