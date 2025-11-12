import { postFoodsMutation } from "@/client/@tanstack/react-query.gen";
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
import { FileWithPreview } from "@/hooks/use-file-upload";
import { AnyFieldApi, useForm } from "@tanstack/react-form";
import { useMutation } from "@tanstack/react-query";
import { image } from "motion/react-client";
import { useCallback } from "react";
import * as z from "zod";
import FluentAdd12Filled from "~icons/fluent/add-12-filled";

function FieldInfo({ field }: { field: AnyFieldApi }) {
  return (
    <>
      {field.state.meta.isTouched && !field.state.meta.isValid ?
        <em>{field.state.meta.errors.map((err) => err.message).join(",")}</em>
      : null}
      {field.state.meta.isValidating ? "Validating..." : null}
    </>
  );
}

const foodSchema = z.object({
  name: z.string().min(1, "Name is required"),
  price: z.number().min(0, "Price must be non-negative").multipleOf(1000),
  description: z.string().default(""),
  imageFile: z.file().nullish(),
  imageFileId: z.string().optional(),
});

const AddFoodDialog = () => {
  const addFood = useMutation(postFoodsMutation());
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
      const result = await addFood.mutateAsync({ body: value });
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

  console.log("handleFileChange callback:", handleFileChange);

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button>
          <FluentAdd12Filled />
          Add
        </Button>
      </DialogTrigger>
      <DialogContent
        className="sm:max-w-md"
        onInteractOutside={(e) => e.preventDefault()}
      >
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
                  {isSubmitting ? "Adding..." : "Add"}
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
