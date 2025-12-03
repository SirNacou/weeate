import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/_settings/notifications")({
  component: RouteComponent,
  staticData: {
    title: "Notification Settings",
  },
});

function RouteComponent() {
  return <div>Hello "/_settings/notifications"!</div>;
}
