import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_protected/settings/routes')({
  component: RouteComponent,
  staticData: {
    title: 'Settings Routes',
  },
})

function RouteComponent() {
  return <div>Hello "/_protected/settings/routes"!</div>
}
