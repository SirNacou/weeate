import { ForgotPasswordForm } from '@/components/forgot-password-form'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth/forgot-password')({
  component: ForgotPassword,
})

function ForgotPassword() {
  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <ForgotPasswordForm />
      </div>
    </div>
  )
}
