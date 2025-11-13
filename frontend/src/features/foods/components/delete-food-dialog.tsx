import {
  deleteFoodsByIdMutation,
  getFoodsQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { Button } from "@/components/animate-ui/components/buttons/button";
import {
  AlertDialog,
  AlertDialogClose,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { useMutation, useQueryClient } from "@tanstack/react-query";

interface DeleteFoodDialog {
  foodId: string;
  open: boolean;
  onOpenChange?: (open: boolean) => void;
}

const DeleteFoodDialog = ({ foodId, open, onOpenChange }: DeleteFoodDialog) => {
  const deleteFood = useMutation({
    ...deleteFoodsByIdMutation(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: getFoodsQueryKey() });
      onOpenChange?.(false);
    },
  });
  const queryClient = useQueryClient();

  const handleDelete = async (): Promise<void> => {
    await deleteFood.mutateAsync({ path: { id: foodId } });
  };

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent className="sm:max-w-[425px]">
        <AlertDialogHeader>
          <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
          <AlertDialogDescription>
            This action cannot be undone. This will permanently delete your food
            from our servers.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogClose>Cancel</AlertDialogClose>
          <Button variant={"destructive"} onClick={handleDelete}>
            Delete
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
};

export default DeleteFoodDialog;
