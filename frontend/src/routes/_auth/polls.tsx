import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth/polls')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/_auth/polls"!</div>
}
