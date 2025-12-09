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
import type z from "zod";
import {
	type CreatePollCommandWritable,
	type GetTodayPollsQueryResponse,
	postPolls,
	serverClient,
} from "@/api";
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
	userPoll?: GetTodayPollsQueryResponse;
};

const MINUTES_BEFORE_NEXT_HOUR = 45;
const TOMORROW_CLOSE_HOUR = 6;
const TOMORROW_CLOSE_MINUTE = 0;

/**
 * Calculate the next available time slot based on current time.
 * If current minutes exceed the threshold, move to the next hour.
 */
const getNextAvailableSlot = (baseTime: Date, now: Date): Date => {
	if (getMinutes(now) > MINUTES_BEFORE_NEXT_HOUR) {
		return addHours(baseTime, 1);
	}
	return baseTime;
};

const CreatePollDialog = ({ userPoll }: Props) => {
	const [open, setOpen] = useState(false);
	const { user } = ProtectedRoute.useRouteContext();
	const queryClient = useQueryClient();
	const getFoodsServer = useServerFn(getFoodsServerFn);
	const createPollServer = useServerFn(createPollServerFn);

	const existed = !!userPoll;
	const disabled = userPoll && userPoll.closed_at !== null;

	const createPoll = useMutation({
		mutationFn: createPollServer,
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

	const { orderDate, defaultCloseTime, timeSlots } = useMemo(() => {
		const orderDate = new Date(add(Date.now(), { days: 1 }));

		// Calculate default close time (next hour from now)
		const now = new Date();
		const defaultCloseTime = getNextAvailableSlot(
			startOfHour(addHours(now, 1)),
			now,
		);

		// Generate time slots
		const slots = [];
		let currentSlot = getNextAvailableSlot(startOfHour(addHours(now, 1)), now);

		const tomorrow = addDays(now, 1);
		const maxTime = setMinutes(
			setHours(tomorrow, TOMORROW_CLOSE_HOUR),
			TOMORROW_CLOSE_MINUTE,
		);

		while (
			isBefore(currentSlot, maxTime) ||
			currentSlot.getTime() === maxTime.getTime()
		) {
			if (isAfter(currentSlot, now)) {
				slots.push(new Date(currentSlot).toISOString());
			}
			currentSlot = addHours(currentSlot, 1);
		}

		return { orderDate, defaultCloseTime, timeSlots: slots };
	}, []);

	const form = useForm({
		defaultValues: {
			food_ids: [] as string[],
			order_date: orderDate.toISOString(),
			scheduled_close_at: defaultCloseTime.toISOString(),
			strategy: "ORDER_PERSONAL_CHOICE",
		} as z.infer<typeof zCreatePollCommandWritable>,
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

	// Group slots into "Today" and "Tomorrow" for better UI
	const todaySlots = timeSlots.filter(
		(d) => new Date(d).getDate() === new Date().getDate(),
	);
	const tomorrowSlots = timeSlots.filter(
		(d) => new Date(d).getDate() !== new Date().getDate(),
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
							<form.Field name="food_ids">
								{(field) => {
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
							</form.Field>

							<form.Field name="scheduled_close_at">
								{(field) => {
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
							</form.Field>

							<form.Field name="strategy">
								{(field) => {
									const isInvalid =
										field.state.meta.isTouched && !field.state.meta.isValid;
									return (
										<Field data-invalid={isInvalid} className="grid gap-1">
											<FieldLabel htmlFor={field.name}>Poll Type</FieldLabel>
											<RadioGroup
												value={field.state.value}
												onValueChange={(value) =>
													field.handleChange(
														value as CreatePollCommandWritable["strategy"],
													)
												}
												onBlur={field.handleBlur}
											>
												<div className="flex items-center gap-3">
													<RadioGroupItem
														value={
															"ORDER_PERSONAL_CHOICE" as CreatePollCommandWritable["strategy"]
														}
														id="ORDER_PERSONAL_CHOICE"
													/>
													<Label htmlFor="ORDER_PERSONAL_CHOICE">
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
							</form.Field>
						</FieldGroup>
					</div>
					{existed && !disabled && (
						<Alert variant="warning" appearance="light" className="mt-4">
							<AlertIcon>
								<LucideAlertTriangle />
							</AlertIcon>
							<AlertContent>
								<AlertTitle>Warning</AlertTitle>
								<AlertDescription>
									A poll for tomorrow order already exists, it will be replaced
									and all existing votes will be lost.
								</AlertDescription>
							</AlertContent>
						</Alert>
					)}
					{disabled && (
						<Alert variant="destructive" appearance="light" className="mt-4">
							<AlertIcon>
								<LucideAlertTriangle />
							</AlertIcon>
							<AlertContent>
								<AlertTitle>Error</AlertTitle>
								<AlertDescription>
									You cannot replace a poll that has already been closed.
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
								disabled={!canSubmit || disabled}
							>
								{isSubmitting ? (
									<>
										<Spinner />
										"Submitting..."
									</>
								) : (
									"Submit"
								)}
							</Button>
						)}
					/>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
};

export default CreatePollDialog;
