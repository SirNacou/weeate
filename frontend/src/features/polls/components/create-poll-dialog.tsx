import { serverClient } from "@/api";
import {
  CreatePollCommand,
  CreatePollCommandWritable,
  postPolls,
} from "@/client";
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
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { MultiSelect, MultiSelectOption } from "@/components/multi-select";
import { useMemo } from "react";
import {
  add,
  addDays,
  addHours,
  format,
  getMinutes,
  isAfter,
  isBefore,
  setHours,
  setMinutes,
  startOfHour,
} from "date-fns";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Spinner } from "@/components/ui/spinner";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";

const createPollServerFn = createServerFn({ method: "POST" })
  .inputValidator(zCreatePollCommandWritable)
  .handler(async ({ data }) => {
    const result = await postPolls({
      client: serverClient,
      body: {
        food_ids: data.food_ids,
        order_date: data.order_date as unknown as Date,
        scheduled_close_at: data.scheduled_close_at as unknown as Date,
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
  const queryClient = useQueryClient();

  const createPoll = useMutation({
    mutationFn: createPollServerFn,
    onSuccess: () =>
      queryClient.invalidateQueries({
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
      try {
        await createPoll.mutateAsync({
          data: value,
        });
        toast.success("Poll created successfully");
      } catch (error) {
        toast.error(
          error instanceof Error ? error.message : "Failed to create poll"
        );
      }
    },
  });

  const timeSlots = useMemo(() => {
    const now = new Date();
    const slots = [];

    // 1. Calculate the Constraint Boundaries

    // Min: Current time + 1 hour, rounded to start of hour
    let currentSlot = startOfHour(addHours(now, 1));

    if (getMinutes(now) > 55) {
      currentSlot = addHours(currentSlot, 1);
    }

    // Max: Tomorrow at 6:00 AM
    const tomorrow = addDays(now, 1);
    const maxTime = setMinutes(setHours(tomorrow, 6), 0); // Tomorrow 06:00

    // 2. Loop to generate slots
    while (
      isBefore(currentSlot, maxTime) ||
      currentSlot.getTime() === maxTime.getTime()
    ) {
      // Safety Check: Ensure we aren't showing past times if logic drifts
      if (isAfter(currentSlot, now)) {
        slots.push(new Date(currentSlot));
      }

      // Increment by 1 hour (Change to 30 for half-hour slots)
      currentSlot = addHours(currentSlot, 1);
    }

    return slots;
  }, []);

  // Group slots into "Today" and "Tomorrow" for better UI
  const todaySlots = timeSlots.filter(
    (d) => d.getDate() === new Date().getDate()
  );
  const tomorrowSlots = timeSlots.filter(
    (d) => d.getDate() !== new Date().getDate()
  );

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
          id="create-poll-form"
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
                      <FieldLabel htmlFor={field.name}>
                        Scheduled Close At
                      </FieldLabel>
                      <Select
                        onValueChange={field.handleChange}
                        value={field.state.value}
                      >
                        <SelectTrigger className="w-full">
                          <SelectValue placeholder="Select closing time..." />
                        </SelectTrigger>
                        <SelectContent>
                          {/* TODAY GROUP */}
                          {todaySlots.length > 0 && (
                            <SelectGroup>
                              <SelectLabel>Today</SelectLabel>
                              {todaySlots.map((date) => (
                                <SelectItem
                                  key={date.toISOString()}
                                  value={date.toISOString()}
                                >
                                  {format(date, "h:00 a")}
                                </SelectItem>
                              ))}
                            </SelectGroup>
                          )}

                          {/* TOMORROW GROUP */}
                          {tomorrowSlots.length > 0 && (
                            <SelectGroup>
                              <SelectLabel>Tomorrow Morning</SelectLabel>
                              {tomorrowSlots.map((date) => (
                                <SelectItem
                                  key={date.toISOString()}
                                  value={date.toISOString()}
                                >
                                  {format(date, "h:00 a")}
                                </SelectItem>
                              ))}
                            </SelectGroup>
                          )}
                        </SelectContent>
                      </Select>
                      {isInvalid && (
                        <FieldError errors={field.state.meta.errors} />
                      )}
                    </Field>
                  );
                }}
              />

              <form.Field
                name="strategy"
                children={(field) => {
                  const isInvalid =
                    field.state.meta.isTouched && !field.state.meta.isValid;
                  return (
                    <Field data-invalid={isInvalid} className="grid gap-1">
                      <FieldLabel htmlFor={field.name}>Poll Type</FieldLabel>
                      <RadioGroup
                        value={field.state.value}
                        onValueChange={(value) =>
                          field.handleChange(
                            value as CreatePollCommandWritable["strategy"]
                          )
                        }
                        onBlur={field.handleBlur}
                      >
                        <div className="flex items-center gap-3">
                          <RadioGroupItem
                            value={
                              "ORDER_MULTIPLE_ITEMS" as CreatePollCommandWritable["strategy"]
                            }
                            id="ORDER_MULTIPLE_ITEMS"
                          />
                          <Label htmlFor="ORDER_MULTIPLE_ITEMS">
                            Order All
                          </Label>
                        </div>

                        <div className="flex items-center gap-3">
                          <RadioGroupItem
                            value={
                              "ORDER_CONSENSUS_ITEM" as CreatePollCommandWritable["strategy"]
                            }
                            id="ORDER_CONSENSUS_ITEM"
                          />
                          <Label htmlFor="ORDER_CONSENSUS_ITEM">
                            Order Most Voted
                          </Label>
                        </div>
                      </RadioGroup>

                      {isInvalid && (
                        <FieldError errors={field.state.meta.errors} />
                      )}
                    </Field>
                  );
                }}
              />
            </FieldGroup>
          </div>
        </form>
        <DialogFooter>
          <DialogClose asChild>
            <Button variant={"outline"}>Cancel</Button>
          </DialogClose>
          <form.Subscribe
            selector={(state) => [
              state.isSubmitting,
              state.canSubmit,
              state.isPristine,
            ]}
            children={([isSubmitting, canSubmit, isPristine]) => (
              <Button
                form="create-poll-form"
                type="submit"
                disabled={isSubmitting || !canSubmit || isPristine}
              >
                {isSubmitting ?
                  <>
                    <Spinner />
                    "Submitting..."
                  </>
                : "Submit"}
              </Button>
            )}
          />
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default CreatePollDialog;
