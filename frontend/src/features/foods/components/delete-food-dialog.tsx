import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogProps,
  AlertDialogTitle,
} from "@/components/animate-ui/components/radix/alert-dialog";
import { useQueryClient } from "@tanstack/react-query";
import { SyntheticEvent } from "react";

type DeleteFoodDialog = AlertDialogProps & {
  foodId: string;
  onOpenChange?: (open: boolean) => void;
};

const DeleteFoodDialog = ({
  foodId,
  onOpenChange,
  ...props
}: DeleteFoodDialog) => {
  const queryClient = useQueryClient();

  const handleDelete = (_e: SyntheticEvent<HTMLButtonElement, Event>): void => {
    // TODO: Implement actual delete API call when backend endpoint is ready
    console.log("Deleting food:", foodId);

    // For now, just invalidate the query and close
    queryClient.invalidateQueries({ queryKey: ["get", "/foods"] });
    onOpenChange?.(false);
  };

  return (
    <AlertDialog {...props}>
      <AlertDialogContent className="sm:max-w-[425px]">
        <AlertDialogHeader>
          <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
          <AlertDialogDescription>
            This action cannot be undone. This will permanently delete your food
            from our servers.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction onSelect={handleDelete}>
            Continue
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
};

export default DeleteFoodDialog;
