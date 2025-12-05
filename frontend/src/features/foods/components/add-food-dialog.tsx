import { postFoods } from "@/api"
import {
	listFoodsQueryKey,
	postFoodsMutation,
} from "@/client/@tanstack/react-query.gen"
import { zPostFoodsData } from "@/client/zod.gen"
import ImageUpload from "@/components/image-upload"
import { Button } from "@/components/ui/button"
import {
	Dialog,
	DialogClose,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog"
import {
	Field,
	FieldError,
	FieldGroup,
	FieldLabel,
} from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { Spinner } from "@/components/ui/spinner"
import type { FileWithPreview } from "@/hooks/use-file-upload"
import { useForm } from "@tanstack/react-form"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { createServerFn, useServerFn } from "@tanstack/react-start"
import { useCallback, useState } from "react"
import * as z from "zod"
import FluentAdd32Filled from "~icons/fluent/add-32-filled"

const addFoodServerFn = createServerFn({ method: "POST" })
	.inputValidator(zPostFoodsData)
	.handler(async ({ data }) => {
		const res = await postFoods({
			body: {
				...data.body,
				price: Number(data.body.price),
			}
		})
		if (res.error) {
			throw new Error(res.error.errors?.at(0)?.message || "Failed to add food")
		}
		return res.data
	})
const foodSchema = z.object({
	name: z.string().nonempty("Name is required"),
	price: z.number().min(0, "Price must be non-negative").multipleOf(1000),
	description: z.string(),
	imageFileId: z.string(),
})

const AddFoodDialog = () => {
	const [open, setOpen] = useState(false)
	const queryClient = useQueryClient()
	const addFoodServer = useServerFn(addFoodServerFn)

	const addFood = useMutation({
		...postFoodsMutation(),
		mutationFn: ({ body }) => addFoodServer({ data: { body: body } }),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: listFoodsQueryKey() })
			setOpen(false)
			form.reset()
		},
	})
	const form = useForm({
		defaultValues: {
			name: "",
			price: 0,
			description: "",
			imageFileId: "",
		},
		validators: {
			onChange: foodSchema,
		},
		onSubmit: async ({ value }) => {
			console.log("Submitting form with value:", value)
			await addFood.mutateAsync({
				body: {
					name: value.name,
					price: value.price,
					description: value.description,
					image_file_id: value.imageFileId,
				},
			})
		},
	})

	const handleFileChange = useCallback(
		(files: FileWithPreview[]) => {
			const file = files.at(0)
			if (file?.file instanceof File) {
				// Update the imageFile field
				setTimeout(() => {
					form.setFieldValue("imageFileId", file.file.name)
				}, 5000)
			}
		},
		[form],
	)

	return (
		<Dialog open={open} onOpenChange={setOpen}>
			<DialogTrigger asChild>
				<Button>
					<FluentAdd32Filled />
					Add
				</Button>
			</DialogTrigger>
			<DialogContent className="sm:max-w-md" aria-description="Add food dialog">
				<form
					className="grid gap-4"
					onSubmit={(e) => {
						e.preventDefault()
						e.stopPropagation()
						form.handleSubmit()
					}}
				>
					<DialogHeader>
						<DialogTitle>New Food</DialogTitle>
					</DialogHeader>
					<div className="grid gap-4">
						<FieldGroup>
							<form.Field name="name">
								{(field) => {
									const isInvalid =
										field.state.meta.isTouched && !field.state.meta.isValid
									return (
										<Field data-invalid={isInvalid}>
											<FieldLabel htmlFor={field.name}>Name</FieldLabel>
											<Input
												id={field.name}
												name={field.name}
												value={field.state.value}
												onBlur={field.handleBlur}
												onChange={(e) => field.handleChange(e.target.value)}
											/>
											{isInvalid && (
												<FieldError errors={field.state.meta.errors} />
											)}
										</Field>
									)
								}}
							</form.Field>

							<form.Field name="price">
								{(field) => {
									const isInvalid =
										field.state.meta.isTouched && !field.state.meta.isValid
									return (
										<Field data-invalid={isInvalid}>
											<FieldLabel htmlFor={field.name}>Price</FieldLabel>
											<Input
												id={field.name}
												name={field.name}
												value={field.state.value}
												onBlur={field.handleBlur}
												type="number"
												onChange={(e) =>
													field.handleChange(Number(e.target.value))
												}
											/>
											{isInvalid && (
												<FieldError errors={field.state.meta.errors} />
											)}
										</Field>
									)
								}}
							</form.Field>

							<form.Field name="description">
								{(field) => {
									const isInvalid =
										field.state.meta.isTouched && !field.state.meta.isValid
									return (
										<Field data-invalid={isInvalid}>
											<FieldLabel htmlFor={field.name}>Description</FieldLabel>
											<Input
												id={field.name}
												name={field.name}
												value={field.state.value}
												onBlur={field.handleBlur}
												type="text"
												onChange={(e) => field.handleChange(e.target.value)}
											/>
											{isInvalid && (
												<FieldError errors={field.state.meta.errors} />
											)}
										</Field>
									)
								}}
							</form.Field>

							<form.Field name="imageFileId">
								{(field) => {
									const isInvalid =
										field.state.meta.isTouched && !field.state.meta.isValid

									return (
										<Field data-invalid={isInvalid}>
											<FieldLabel htmlFor={field.name}>Image</FieldLabel>

											<ImageUpload
												maxSizeMB={5}
												onFilesChange={handleFileChange}
											/>

											<Input
												hidden={true}
												id={field.name}
												name={field.name}
												value={field.state.value}
												readOnly
											/>
											{isInvalid && (
												<FieldError errors={field.state.meta.errors} />
											)}
										</Field>
									)
								}}
							</form.Field>
						</FieldGroup>
					</div>
					<DialogFooter>
						<DialogClose asChild>
							<Button variant="outline">Cancel</Button>
						</DialogClose>
						<form.Subscribe
							selector={(state) => [
								state.canSubmit,
								state.isPristine,
								state.isSubmitting,
							]}
						>
							{([canSubmit, isPristine, isSubmitting]) => (
								<Button
									type="submit"
									disabled={!canSubmit || isPristine || isSubmitting}
								>
									{isSubmitting ? (
										<>
											<Spinner />
											Adding...
										</>
									) : (
										"Add"
									)}
								</Button>
							)}
						</form.Subscribe>
					</DialogFooter>
				</form>
			</DialogContent>
		</Dialog>
	)
}

export default AddFoodDialog
