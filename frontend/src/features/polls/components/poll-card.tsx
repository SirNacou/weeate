import {
  type GetTodayPollsQueryResponse,
  postPollsByIdVote,
  serverClient,
} from "@/api"
import { listPollsTodayQueryKey } from "@/client/@tanstack/react-query.gen"
import { zPostPollsByIdVoteData } from "@/client/zod.gen"
import CloseTimer from "@/components/close-timer"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { FieldGroup } from "@/components/ui/field"
import { Route } from "@/routes/_protected/route"
import { useForm } from "@tanstack/react-form"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { createServerFn } from "@tanstack/react-start"
import { useEffect, useMemo } from "react"
import { toast } from "sonner"
import { z } from "zod"
import PollOptionRadio from "./poll-option"
import PollStrategyBadge from "./poll-strategy-badge"

const castVoteServerFn = createServerFn({ method: "POST" })
  .inputValidator(zPostPollsByIdVoteData)
  .handler(async ({ data }) => {
    const res = await postPollsByIdVote({
      client: serverClient,
      path: data.path,
      body: data.body,
    })
    if (res.error) {
      console.error(`Failed to cast vote:`, res.error)
      throw new Error("Failed to cast vote")
    }

    return res.data
  })

type Props = {
  poll: GetTodayPollsQueryResponse
}

const PollCard = ({ poll }: Props) => {
  const { user } = Route.useRouteContext()
  const {
    id,
    creator: { display_name: buyerName, avatar_url: avatarUrl },
    scheduled_closes_at,
    closed_at,
    strategy,
    poll_options,
  } = poll

  const initialSelectedOption = poll_options?.find((option) =>
    option.votes?.some((vote) => vote.voter?.id === user.id),
  )

  const queryClient = useQueryClient()

  const validator = useMemo(() => {
    return z.object({
      selectedOption: z.object({
        id: z
          .string()
          .nonempty({ message: "Please select an option" })
          .refine((val) => poll_options.some((option) => option.id === val), {
            message: "Invalid option selected",
          }),
        isSelected: z.boolean(),
      }),
    })
  }, [poll_options])

  const castVoteMutation = useMutation({
    mutationFn: castVoteServerFn,
    onSuccess: () => {
      toast.success("Your vote has been submitted!")
      queryClient.invalidateQueries({
        queryKey: listPollsTodayQueryKey(),
      })
    },
    onError: (error) => {
      toast.error(error.message || "Failed to cast vote. Please try again.")
    },
  })

  const form = useForm({
    defaultValues: {
      selectedOption: {
        id: initialSelectedOption?.id || "",
        isSelected: !!initialSelectedOption,
      },
    },
    validators: {
      onChange: validator,
      onMount: validator,
    },
    listeners: {
      onChangeDebounceMs: 500,
      onChange(props) {
        props.formApi.handleSubmit()
      },
    },
    onSubmit: async ({ value }) => {
      await castVoteMutation
        .mutateAsync({
          data: {
            path: {
              id,
            },
            body: {
              poll_option_id: value.selectedOption.id,
            },
          },
        })
        .catch(() => { })
    },
  })
  useEffect(() => { }, [form.state.values.selectedOption])

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault()
        e.stopPropagation()
        form.handleSubmit()
      }}
    >
      <Card>
        <CardHeader>
          <CardTitle className="flex flex-col justify-between items-center sm:flex-row text-lg">
            <div className="flex items-center gap-2">
              <Avatar className="size-8 sm:size-10">
                <AvatarImage src={avatarUrl} />
                <AvatarFallback>
                  {buyerName.substring(0, 2).toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <h3 className="text-base sm:text-lg leading-tight">
                <span className="text-primary font-semibold">{buyerName}</span>{" "}
                is buying breakfast for{" "}
                {new Date(poll.order_date).toLocaleDateString("vi-VN")}
              </h3>
            </div>
          </CardTitle>
          <CardDescription className="mt-1 flex flex-row flex-wrap items-center gap-3 justify-between text-sm sm:text-base">
            <CloseTimer
              className="text-sm sm:text-base shrink-0"
              closesAt={closed_at || scheduled_closes_at}
            />

            <PollStrategyBadge strategy={strategy} />
          </CardDescription>
        </CardHeader>
        <CardContent>
          <FieldGroup>
            <form.Field
              name="selectedOption"
              children={(field) => {
                return (
                  <div className={"grid grid-cols-1 lg:grid-cols-2 gap-3"}>
                    {poll_options.map((option) => (
                      <PollOptionRadio
                        key={option.id}
                        option={option}
                        disabled={!!closed_at}
                        isSelected={
                          option.id === field.state.value.id &&
                          field.state.value.isSelected
                        }
                        onSelect={(id) => {
                          if (field.state.value.id === id) {
                            // Deselecting the currently selected option
                            field.handleChange({
                              id,
                              isSelected: !field.state.value.isSelected,
                            })
                          } else {
                            // Selecting a new option
                            field.handleChange({ id, isSelected: true })
                          }
                        }}
                      />
                    ))}
                  </div>
                )
              }}
            />
          </FieldGroup>
        </CardContent>
        <CardFooter>
          <div className="flex flex-col sm:flex-row items-center justify-end w-full text-lg">
            <span className="font-bold mr-2">Total Cost:</span>
            <span className="font-medium text-slate-700">
              {Intl.NumberFormat("vi-VN", {
                style: "currency",
                currency: "VND",
              }).format(
                poll_options.reduce(
                  (acc, option) =>
                    acc + option.price_at_creation * (option.votes.length || 0),
                  0,
                ),
              )}
            </span>
          </div>
        </CardFooter>
      </Card>
    </form>
  )
}

export default PollCard
