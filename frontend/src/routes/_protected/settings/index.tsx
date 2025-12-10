import { getAuthImagekitToken } from '@/api'
import { getAuthImagekitTokenQueryKey } from '@/client/@tanstack/react-query.gen'
import AvatarUpload from '@/components/avatar-upload'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter } from '@/components/ui/card'
import { FieldError } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Spinner } from '@/components/ui/spinner'
import { getFormSubmissionStatus } from '@/lib/form-utils'
import { createSupabaseClient } from '@/lib/supabase'
import { fetchUser as fetchUserServerFn } from '@/lib/supabase/fetch-user-server-fn'
import { useForm } from '@tanstack/react-form'
import { useMutation, useQuery } from '@tanstack/react-query'
import { createFileRoute, useRouter } from '@tanstack/react-router'
import { createServerFn } from '@tanstack/react-start'
import { toast } from 'sonner'
import z from 'zod'

export const Route = createFileRoute('/_protected/settings/')({
  component: RouteComponent,
  staticData: {
    title: 'Settings',
  }
})

const FormSchema = z.object({
  displayName: z.string()
    .min(2, 'Display name must be at least 2 characters long')
    .max(30, 'Display name must be at most 30 characters long'),
})

const updateUserProfileServerFn = createServerFn({ method: 'POST' })
  .inputValidator(FormSchema)
  .handler(async ({ data }) => {
    const supabase = await createSupabaseClient()
    const user = await fetchUserServerFn()
    const res = await supabase.from('user_profiles')
      .update({
        display_name: data.displayName
      })
      .eq('id', user?.id)
      .select()

    if (res.error) {
      throw res.error
    }
    console.log('User profile updated:', res)

    return res.status === 200
  })

const getImageKitTokenServerFn = createServerFn({ method: 'GET' })
  .handler(async () => {
    const res = await getAuthImagekitToken()
    if (res.error) {
      throw res.error
    }
    return res.data
  })


function RouteComponent() {
  const router = useRouter()
  const { user } = Route.useRouteContext()
  const updateUserProfileMutation = useMutation({
    mutationFn: (data: z.infer<typeof FormSchema>) => updateUserProfileServerFn({ data }),
    onSuccess: async () => {
      router.invalidate()
      toast.success('Profile updated successfully')
    }
  })
  const { refetch } = useQuery({
    queryKey: getAuthImagekitTokenQueryKey(),
    queryFn: getImageKitTokenServerFn,
  })

  const form = useForm({
    defaultValues: {
      displayName: user.app_metadata.display_name,
    },
    validators: {
      onMount: FormSchema,
      onChange: FormSchema,
    },
<<<<<<< ours
    onSubmit: async ({ value }) => {
      try {
        await updateUserProfileServerFn({ data: { displayName: value.displayName } })
      } catch (err) {
        // Optionally, you can return or throw the error for the form to display
        throw err
      }
    }
||||||| ancestor
    onSubmit: async ({ value }) => {
      try {
        await updateUserProfileServerFn({ displayName: value.displayName })
      } catch (err) {
        // Optionally, you can return or throw the error for the form to display
        throw err
      }
    }
=======
    onSubmit: ({ value }) => updateUserProfileMutation.mutateAsync({
      displayName: value.displayName
    })
>>>>>>> theirs
  })

  async function onAvatarSave(croppedImage: string): Promise<void> {
    const { data } = await refetch()

    if (!data) {
      toast.error('Failed to get upload token')
      return
    }

    console.log('Uploading avatar with token:', data)
  }

  return <div className='flex flex-col gap-6'>
    <AvatarUpload
      initialImage={user.app_metadata.avatar_url}
      onSave={onAvatarSave}
    />

    <form onSubmit={(e) => {
      e.preventDefault()
      form.handleSubmit()
    }}>
      <Card>
        <CardContent>
          <div className='grid md:grid-cols-2 gap-4'>
            <form.Field name='displayName'>
              {(field) => (
                <div className='flex flex-col gap-1'>
                  <Label className='text-sm font-medium'>Display Name</Label>
                  <Input
                    id={field.name}
                    name={field.name}
                    type='text'
                    value={field.state.value}
                    onChange={(e) => field.setValue(e.target.value)}
                    className='w-full rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500'
                  />
                  {field.state.meta.errors && (
                    <FieldError errors={field.state.meta.errors} />
                  )}
                </div>
              )}
            </form.Field>
          </div>
        </CardContent>
        <CardFooter>
          <form.Subscribe selector={getFormSubmissionStatus}>
            {({ canSubmit, isSubmitting }) => (
              <Button type='submit' disabled={!canSubmit}>
                {isSubmitting ?
                  <>
                    <Spinner />
                    Saving...
                  </>
                  : 'Save Changes'}
              </Button>
            )}
          </form.Subscribe>
        </CardFooter>
      </Card>
    </form>
  </div>
}
