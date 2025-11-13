import {
  getFoodsQueryKey,
  putFoodsByIdMutation,
} from "@/client/@tanstack/react-query.gen";
import ImageUpload from "@/components/comp-545";
import * as DialogPrimitive from "@radix-ui/react-dialog";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
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
import { useCallback } from "react";
import * as z from "zod";

const foodSchema = z.object({
  name: z.string().min(1, "Name is required"),
  price: z.number().min(0, "Price must be non-negative").multipleOf(1000),
  description: z.string().default(""),
  imageFileId: z.string().optional(),
});

type EditFoodDialogProps = React.ComponentProps<typeof DialogPrimitive.Root> & {
  data: {
    id: string;
    name: string;
    price: number;
    description: string;
    imageFileId?: string;
    imageFileUrl?: string;
  };
  onOpenChange?: (open: boolean) => void;
};

const AddFoodDialog = ({
  data,
  onOpenChange,
  ...props
}: EditFoodDialogProps) => {
  const queryClient = useQueryClient();
  const updateFood = useMutation({
    ...putFoodsByIdMutation(),
    onSuccess: () => {
      // Invalidate and refetch foods query
      queryClient.invalidateQueries({ queryKey: getFoodsQueryKey() });
      // Close the dialog
      onOpenChange?.(false);
    },
  });
  const form = useForm({
    defaultValues: {
      name: data.name || "",
      price: data.price || 0,
      description: data.description || "",
      imageFileId: data.imageFileId || "",
    },
    validators: {
      // @ts-ignore
      onChange: foodSchema,
    },
    onSubmit: async ({ value }) => {
      try {
        await updateFood.mutateAsync({
          path: { id: data.id },
          body: {
            name: value.name,
            price: value.price,
            description: value.description,
            image_file_id: value.imageFileId,
          },
        });
      } catch (error) {
        console.error("Submit error:", error);
        throw error;
      }
    },
  });

  const handleFileChange = useCallback(
    (files: FileWithPreview[]) => {
      console.log("File changed:", files);
      const file = files.at(0);
      if (file?.file instanceof File) {
        // Update the imageFile field
        setTimeout(() => {
          form.setFieldValue("imageFileId", file.file.name);
        }, 5000);
      }
      console.log("handleFileChange completed");
    },
    [form]
  );

  console.log("handleFileChange callback:", handleFileChange);

  return (
    <Dialog {...props}>
      <DialogContent
        className="sm:max-w-md"
        aria-description="Edit food dialog"
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
            <DialogTitle>Edit Food</DialogTitle>
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
                      Updating...
                    </>
                  : "Update"}
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
