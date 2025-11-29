import { UpdatePasswordForm } from '@/components/update-password-form'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth/update-password')({
  component: UpdatePassword,
})

function UpdatePassword() {
  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <UpdatePasswordForm />
      </div>
    </div>
  )
}
