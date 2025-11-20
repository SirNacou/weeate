import { serverClient } from "@/api";
import { CreatePollCommand, postPolls } from "@/client";
import {
  listFoodsQueryKey,
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
import { useForm } from "@tanstack/react-form";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createServerFn } from "@tanstack/react-start";
import LucidePlus from "~icons/lucide/plus?width=2em&height=2em";
import { getFoodsServer } from "@/features/foods/functions/get-server-foods";
import { Route as ProtectedRoute } from "@/routes/_protected/route";
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { MultiSelect, MultiSelectOption } from "@/components/multi-select";
import { useMemo } from "react";
import { Label } from "@/components/ui/label";
import { add } from "date-fns";
import { DateTimePicker } from "@/components/datetime-picker";

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
  const { user } = ProtectedRoute.useRouteContext();

  const createPoll = useMutation({
    mutationFn: createPollServerFn,
    onSuccess: () =>
      useQueryClient().invalidateQueries({
        queryKey: listPollsTodayQueryKey(),
      }),
  });
  const { data: foods } = useQuery({
    queryKey: listFoodsQueryKey({ query: { user_id: user.id } }),
    queryFn: () => getFoodsServer({ data: { query: { user_id: user.id } } }),
  });

  const foodOptions = useMemo(() => {
    return foods ?
        foods.map(
          (food) =>
            ({
              label: food.name,
              value: food.id,
            }) as MultiSelectOption
        )
      : [];
  }, [foods]);

  const minDate = useMemo(() => {
    const date = add(new Date(), { hours: 1 });
    if (date.getMinutes() > 0 || date.getSeconds() > 0) {
      date.setHours(date.getHours() + 1, 0, 0, 0);
    } else {
      date.setMinutes(0, 0, 0);
    }
    return date;
  }, []);

  const maxDate = useMemo(() => {
    const date = add(new Date(), { days: 1 });
    date.setHours(6, 0, 0, 0);
    return date;
  }, []);

  const orderDate = new Date(add(Date.now(), { days: 1 }));
  const form = useForm({
    defaultValues: {
      food_ids: [] as CreatePollCommand["food_ids"],
      order_date: orderDate.toISOString(),
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

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button className="w-full sm:w-auto">
          <LucidePlus />
          Create New
        </Button>
      </DialogTrigger>
      <DialogContent className="flex flex-col gap-6">
        <DialogHeader>
          <DialogTitle>Create New Poll</DialogTitle>
          <DialogDescription>
            Fill out the form below to create a new food poll.
          </DialogDescription>
        </DialogHeader>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
        >
          <div className="flex flex-col gap-1">
            <FieldGroup>
              <form.Field
                name="food_ids"
                children={(field) => {
                  const isInvalid =
                    field.state.meta.isTouched && !field.state.meta.isValid;
                  return (
                    <Field data-invalid={isInvalid} className="grid gap-1">
                      <FieldLabel htmlFor={field.name}>Select Foods</FieldLabel>
                      <MultiSelect
                        id={field.name}
                        name={field.name}
                        variant={"default"}
                        onBlur={field.handleBlur}
                        options={foodOptions}
                        value={field.state.value || []}
                        onValueChange={field.handleChange}
                      />
                      {isInvalid && (
                        <FieldError errors={field.state.meta.errors} />
                      )}
                    </Field>
                  );
                }}
              />

              <form.Field
                name="scheduled_close_at"
                children={(field) => {
                  const isInvalid =
                    field.state.meta.isTouched && !field.state.meta.isValid;
                  return (
                    <Field data-invalid={isInvalid} className="grid gap-1">
                      <FieldLabel htmlFor={field.name}>Select Foods</FieldLabel>
                      <DateTimePicker
                        onChange={(value) =>
                          field.handleChange(value?.toISOString() || "")
                        }
                        value={new Date(field.state.value)}
                        min={minDate}
                        max={maxDate}
                        timePicker={{
                          hour: true,
                          minute: false,
                          second: false,
                        }}
                        
                      />
                      {isInvalid && (
                        <FieldError errors={field.state.meta.errors} />
                      )}
                    </Field>
                  );
                }}
              />
            </FieldGroup>
          </div>
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
