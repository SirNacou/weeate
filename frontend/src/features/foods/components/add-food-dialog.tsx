import {
  getFoodsQueryKey,
  postFoodsMutation,
} from "@/client/@tanstack/react-query.gen";
import ImageUpload from "@/components/comp-545";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";
import { FileWithPreview } from "@/hooks/use-file-upload";
import { useForm } from "@tanstack/react-form";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useCallback, useState } from "react";
import * as z from "zod";
import FluentAdd32Filled from "~icons/fluent/add-32-filled";

const foodSchema = z.object({
  name: z.string().min(1, "Name is required"),
  price: z.number().min(0, "Price must be non-negative").multipleOf(1000),
  description: z.string().default(""),
  imageFile: z.file().nullish(),
  imageFileId: z.string().optional(),
});

const AddFoodDialog = () => {
  const [open, setOpen] = useState(false);
  const queryClient = useQueryClient();

  const addFood = useMutation({
    ...postFoodsMutation(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: getFoodsQueryKey() });
      setOpen(false);
      form.reset();
    },
  });
  const form = useForm({
    defaultValues: {
      name: "",
      price: 0,
      description: "",
      imageFile: null as File | null,
      imageFileId: "",
    },
    validators: {
      // @ts-ignore
      onChange: foodSchema,
    },
    onSubmit: async ({ value }) => {
      console.log("Submitted values:", value);
      const result = await addFood.mutateAsync({
        body: {
          name: value.name,
          price: value.price,
          description: value.description,
          image_file_id: value.imageFileId || undefined,
        },
      });
      console.log("Add food result:", result);
    },
  });

  const handleFileChange = useCallback(
    (files: FileWithPreview[]) => {
      console.log("File changed:", files);
      const file = files.at(0);
      if (file?.file instanceof File) {
        // Update the imageFile field
        setTimeout(() => {
          form.setFieldValue("imageFile", file.file as File);
          form.setFieldValue("imageFileId", file.file.name);
        }, 5000);
      }
      console.log("handleFileChange completed");
    },
    [form]
  );

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
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
        >
          <DialogHeader>
            <DialogTitle>New Food</DialogTitle>
          </DialogHeader>
          <div className="grid gap-4">
            <FieldGroup>
              <form.Field
                name="name"
                children={(field) => {
                  const isInvalid =
                    field.state.meta.isTouched && !field.state.meta.isValid;
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
                  );
                }}
              />

              <form.Field
                name="price"
                children={(field) => {
                  const isInvalid =
                    field.state.meta.isTouched && !field.state.meta.isValid;
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
                  );
                }}
              />

              <form.Field
                name="description"
                children={(field) => {
                  const isInvalid =
                    field.state.meta.isTouched && !field.state.meta.isValid;
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
                  );
                }}
              />

              <form.Field
                name="imageFileId"
                children={(field) => {
                  const isInvalid =
                    field.state.meta.isTouched && !field.state.meta.isValid;

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
                  );
                }}
              />
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
              children={([canSubmit, isPristine, isSubmitting]) => (
                <Button
                  type="submit"
                  disabled={!canSubmit || isPristine || isSubmitting}
                >
                  {isSubmitting ?
                    <>
                      <Spinner />
                      Adding...
                    </>
                  : "Add"}
                </Button>
              )}
            />
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export default AddFoodDialog;
