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
import {
  Alert,
  AlertContent,
  AlertDescription,
  AlertIcon,
  AlertTitle,
} from "@/components/ui/alert";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Label } from "@/components/ui/label";
import {
  MultiSelect,
  MultiSelectContent,
  MultiSelectGroup,
  MultiSelectItem,
  MultiSelectTrigger,
  MultiSelectValue,
} from "@/components/ui/multi-select";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
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
import { getFoodsServer as getFoodsServerFn } from "@/features/foods/functions/get-server-foods";
import { getFormSubmissionStatus } from "@/lib/form-utils";
import { Route as ProtectedRoute } from "@/routes/_protected/route";
import { useForm } from "@tanstack/react-form";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createServerFn, useServerFn } from "@tanstack/react-start";
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
import { useMemo, useState } from "react";
import { toast } from "sonner";
import LucideAlertTriangle from "~icons/lucide/alert-triangle";
import LucidePlus from "~icons/lucide/plus?width=2em&height=2em";

const createPollServerFn = createServerFn({ method: "POST" })
  .inputValidator(zCreatePollCommandWritable)
  .handler(async ({ data }) => {
    const result = await postPolls({
      client: serverClient,
      body: {
        food_ids: data.food_ids,
        order_date: data.order_date,
        scheduled_close_at: data.scheduled_close_at,
        strategy: data.strategy,
      },
    });
    if (result.error) {
      throw result.error.errors?.at(0) || new Error("Failed to create poll");
    }
    return result.data;
  });

type Props = {
  userPollExists?: boolean;
};

const CreatePollDialog = ({ userPollExists }: Props) => {
  const [open, setOpen] = useState(false);
  const { user } = ProtectedRoute.useRouteContext();
  const queryClient = useQueryClient();
  const getFoodsServer = useServerFn(getFoodsServerFn);

  const createPoll = useMutation({
    mutationFn: createPollServerFn,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: listPollsTodayQueryKey(),
      });
      toast.success("Poll created successfully");
      setOpen(false);
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });
  const { data: foods } = useQuery({
    queryKey: listFoodsQueryKey({ query: { user_id: user.id } }),
    queryFn: () => getFoodsServer({ data: { query: { user_id: user.id } } }),
  });

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
      onMount: zCreatePollCommandWritable,
    },
    onSubmit: async ({ value }) => {
      await createPoll
        .mutateAsync({
          data: value,
        })
        .catch(() => {
          // Error is handled in useMutation onError
        });
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
        slots.push(new Date(currentSlot).toISOString());
      }

      // Increment by 1 hour (Change to 30 for half-hour slots)
      currentSlot = addHours(currentSlot, 1);
    }

    return slots;
  }, []);

  // Group slots into "Today" and "Tomorrow" for better UI
  const todaySlots = timeSlots.filter(
    (d) => new Date(d).getDate() === new Date().getDate()
  );
  const tomorrowSlots = timeSlots.filter(
    (d) => new Date(d).getDate() !== new Date().getDate()
  );

  return (
    <Dialog open={open} onOpenChange={setOpen}>
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
                        defaultValues={field.state.value ?? []}
                        values={field.state.value ?? []}
                        onValuesChange={field.handleChange}
                      >
                        <MultiSelectTrigger className="w-full">
                          <MultiSelectValue placeholder="Select food..." />
                        </MultiSelectTrigger>
                        <MultiSelectContent>
                          <MultiSelectGroup>
                            {foods &&
                              foods.map((food) => (
                                <MultiSelectItem key={food.id} value={food.id}>
                                  {food.name}
                                </MultiSelectItem>
                              ))}
                          </MultiSelectGroup>
                        </MultiSelectContent>
                      </MultiSelect>
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
                                <SelectItem key={date} value={date}>
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
                                <SelectItem key={date} value={date}>
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
          {userPollExists && (
            <Alert variant="warning" appearance="light" className="mt-4">
              <AlertIcon>
                <LucideAlertTriangle />
              </AlertIcon>
              <AlertContent>
                <AlertTitle>Warning</AlertTitle>
                <AlertDescription>
                  A poll already exists for today, it will be replaced and all
                  existing votes will be lost.
                </AlertDescription>
              </AlertContent>
            </Alert>
          )}
        </form>
        <DialogFooter>
          <DialogClose asChild>
            <Button variant={"outline"}>Cancel</Button>
          </DialogClose>
          <form.Subscribe
            selector={getFormSubmissionStatus}
            children={({ canSubmit, isSubmitting }) => (
              <Button
                form="create-poll-form"
                type="submit"
                disabled={!canSubmit}
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
