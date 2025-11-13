import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/protected")({
  component: Info,
  loader: async ({ context }) => {
    return {
      user: context.user,
    };
  },
});

function Info() {
  const { user } = Route.useLoaderData();
  return <div>Hello {user.app_metadata.display_name}</div>;
}
