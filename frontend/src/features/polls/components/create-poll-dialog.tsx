import { serverClient } from "@/api";
import { CreatePollCommand, listPollsToday, postPolls } from "@/client";
import {
  listOrdersTodayOptions,
  listPollsTodayOptions,
  listPollsTodayQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { zCreatePollCommandWritable } from "@/client/zod.gen";
import { Button } from "@/components/animate-ui/components/buttons/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/animate-ui/components/radix/dialog";
import { LogoutButton } from "@/components/logout-button";
import { useForm } from "@tanstack/react-form";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createServerFn, useServerFn } from "@tanstack/react-start";
import LucidePlus from "~icons/lucide/plus?width=2em&height=2em";
import { listPollsTodayServerFn } from "../api/get-today-polls";

const createPollServerFn = createServerFn({ method: "POST" })
  .inputValidator(zCreatePollCommandWritable)
  .handler(async ({ data }) => {
    const result = await postPolls({
      client: serverClient,
      body: {
        food_ids: data.food_ids,
        order_date: new Date(data.order_date),
        scheduled_close_at: new Date(data.scheduled_close_at),
        strategy: data.strategy,
      },
    });
    if (result.error) {
      throw new Error(result.error.detail || "Failed to create poll");
    }
    return result.data;
  });

type Props = {};

const CreatePollDialog = ({}: Props) => {
  const createPoll = useMutation({
    mutationFn: createPollServerFn,
    onSuccess: () =>
      useQueryClient().invalidateQueries({
        queryKey: listPollsTodayQueryKey(),
      }),
  });
  const { data, isLoading, error } = useQuery({
    queryKey: listPollsTodayQueryKey(),
    queryFn: listPollsTodayServerFn,
  });

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error loading polls: {(error as Error).message}</div>;
  }

  const form = useForm({
    defaultValues: {
      food_ids: [] as CreatePollCommand["food_ids"],
      order_date: new Date(Date.now()).toISOString(),
      scheduled_close_at: new Date(Date.now()).toISOString(),
      strategy: "ORDER_MULTIPLE_ITEMS" as CreatePollCommand["strategy"],
    },
    validators: {
      onChange: zCreatePollCommandWritable,
    },
    onSubmit: async ({ value }) => {
      console.log("Submitted poll data:", value);
      await createPoll.mutateAsync({
        data: value,
      });
    },
  });
  LogoutButton;

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button className="w-full sm:w-auto">
          <LucidePlus />
          Create New
        </Button>
      </DialogTrigger>
      <DialogContent>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
        >
          <DialogHeader>
            <DialogTitle>Create New Poll</DialogTitle>
            <DialogDescription>
              Fill out the form below to create a new food poll.
            </DialogDescription>
          </DialogHeader>
          <div className="flex flex-col gap-1"></div>
          <DialogFooter>
            <DialogClose asChild>
              <Button variant={"outline"}>Cancel</Button>
            </DialogClose>
            <Button type="submit">Submit</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export default CreatePollDialog;
