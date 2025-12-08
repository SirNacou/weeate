import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/settings/notifications")({
  component: RouteComponent,
  staticData: {
    title: "Notifications",
  },
});

function RouteComponent() {
  return <div>Hello "/_settings/notifications"!</div>;
}
