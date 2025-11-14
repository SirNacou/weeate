import {
  DropdownMenuTrigger,
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
} from "@/components/animate-ui/components/radix/dropdown-menu";
import { Button } from "@/components/ui/button";
import FluentMoreHorizontal32Filled from "~icons/fluent/more-horizontal-32-filled";
import EditFoodDialog from "./edit-food-dialog";
import { useState } from "react";
import DeleteFoodDialog from "./delete-food-dialog";

const TableActionMenuDialog = ({
  data,
}: {
  data: {
    id: string;
    name: string;
    price: number;
    description: string;
    imageFileId?: string;
    imageFileUrl?: string;
  };
}) => {
  const [editOpen, setEditOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);

  return (
    <>
      <DropdownMenu modal={false}>
        <DropdownMenuTrigger asChild>
          <Button variant={"outline"} size={"icon"} aria-label="Menu">
            <FluentMoreHorizontal32Filled />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" side="bottom">
          <DropdownMenuItem
            onSelect={(e) => {
              e.preventDefault();
              setEditOpen(true);
            }}
          >
            Edit
          </DropdownMenuItem>
          <DropdownMenuItem
            variant="destructive"
            onSelect={(e) => {
              e.preventDefault();
              setDeleteOpen(true);
            }}
          >
            Delete
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <EditFoodDialog data={data} open={editOpen} onOpenChange={setEditOpen} />
      <DeleteFoodDialog
        foodId={data.id}
        open={deleteOpen}
        onOpenChange={setDeleteOpen}
      />
    </>
  );
};

export default TableActionMenuDialog;
