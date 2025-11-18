import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_protected/polls/active/')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/_protected/polls/active/"!</div>
}
