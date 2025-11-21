import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import PollStrategyBadge from "./poll-strategy-badge";
import CloseTimer from "@/components/close-timer";
import { RadioGroup } from "@/components/ui/radio-group";
import { useForm } from "@tanstack/react-form";
import { z } from "zod";
import { useMemo } from "react";
import { FieldGroup } from "@/components/ui/field";
import { Button } from "@/components/ui/button";
import PollOption, { Option } from "./poll-option";
import { Spinner } from "@/components/ui/spinner";
import { toast } from "sonner";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createServerFn } from "@tanstack/react-start";
import { zPostPollsByIdVoteData } from "@/client/zod.gen";
import { serverClient } from "@/api";
import { postPollsByIdVote } from "@/client";
import { getFormSubmissionStatus } from "../../../lib/form-utils";
import { listPollsTodayQueryKey } from "@/client/@tanstack/react-query.gen";

const castVoteServerFn = createServerFn({ method: "POST" })
  .inputValidator(zPostPollsByIdVoteData)
  .handler(async ({ data }) => {
    const res = await postPollsByIdVote({
      client: serverClient,
      path: data.path,
      body: data.body,
    });
    if (res.error) {
      console.error(`Failed to cast vote:`, res.error);
      throw new Error("Failed to cast vote");
    }

    return res.data;
  });

type Props = {
  pollId: string;
  buyerName: string;
  avatarUrl: string;
  scheduled_close_at: Date;
  closed_at: Date | null;
  strategy: "ORDER_CONSENSUS_ITEM" | "ORDER_PERSONAL_CHOICE";
  options: Option[];
  initialSelectedOption?: string;
};

const PollCard = ({
  pollId,
  buyerName,
  avatarUrl,
  scheduled_close_at,
  closed_at,
  strategy,
  options,
  initialSelectedOption,
}: Props) => {
  const validator = useMemo(() => {
    return z.object({
      selectedOption: z
        .string()
        .nonempty({ message: "Please select an option" })
        .refine((val) => options.some((option) => option.id === val), {
          message: "Invalid option selected",
        }),
    });
  }, [options]);
  const queryClient = useQueryClient();

  const castVoteMutation = useMutation({
    mutationFn: castVoteServerFn,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: listPollsTodayQueryKey(),
      });
    },
  });

  const form = useForm({
    defaultValues: {
      selectedOption: initialSelectedOption || "",
    },
    validators: {
      onChange: validator,
      onMount: validator,
    },
    listeners: {
      onChange(props) {
        console.log("Form changed:", props.fieldApi.state.value);
      },
    },
    onSubmit: async ({ value }) => {
      try {
        await castVoteMutation.mutateAsync({
          data: {
            path: {
              id: pollId,
            },
            body: {
              poll_option_id: value.selectedOption,
            },
          },
        });
        toast.success("Your vote has been submitted!");
      } catch (error) {
        console.error("Error submitting vote:", error);
        toast.error("There was an error submitting your vote.");
      }
    },
  });

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        e.stopPropagation();
        form.handleSubmit();
      }}
    >
      <Card>
        <CardHeader>
          <CardTitle className="flex flex-col sm:flex-row text-lg">
            <div className="flex items-center gap-2">
              <Avatar className="size-8 sm:size-10">
                <AvatarImage src={avatarUrl} />
                <AvatarFallback>
                  {buyerName.substring(0, 2).toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <h3 className="text-base sm:text-lg leading-tight">
                <span className="text-primary font-semibold">{buyerName}</span>{" "}
                is buying breakfast
              </h3>
            </div>
          </CardTitle>
          <CardDescription className="mt-1 flex items-center justify-between text-sm sm:text-base">
            <CloseTimer
              className="text-sm sm:text-base shrink-0"
              closesAt={closed_at || scheduled_close_at}
            />

            <PollStrategyBadge strategy={strategy} />
          </CardDescription>
        </CardHeader>
        <CardContent className="p-4 sm:p-6 pt-0 sm:pt-0">
          <FieldGroup>
            <form.Field
              name="selectedOption"
              children={(field) => {
                return (
                  <RadioGroup
                    className={"grid grid-cols-1 lg:grid-cols-2 gap-3"}
                    value={field.state.value}
                    onValueChange={field.handleChange}
                  >
                    {options.map((option) => (
                      <PollOption
                        key={option.id}
                        option={option}
                        disabled={!!closed_at}
                        isSelected={option.id === field.state.value}
                        onSelect={(id) => field.handleChange(id)}
                      />
                    ))}
                  </RadioGroup>
                );
              }}
            />
          </FieldGroup>
        </CardContent>
        <CardFooter>
          <form.Subscribe
            selector={getFormSubmissionStatus}
            children={({ canSubmit, isSubmitting }) => (
              <Button
                type="submit"
                className="w-full"
                disabled={!canSubmit || !!closed_at}
              >
                {isSubmitting ?
                  <>
                    <Spinner />
                    Submitting...
                  </>
                : "Submit"}
              </Button>
            )}
          />
        </CardFooter>
      </Card>
    </form>
  );
};

export default PollCard;
