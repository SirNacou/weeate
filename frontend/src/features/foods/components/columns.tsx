import { GetFoodsResponseItem } from "@/client";
import { ColumnDef } from "@tanstack/react-table";
import { Image } from "@imagekit/react";
import TableActionMenuDialog from "./table-action-menu-dialog";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

export const createColumns = (
  currentUserId?: string
): ColumnDef<GetFoodsResponseItem>[] => [
  {
    id: "name",
    header: "Name",
    accessorFn: (row) => row.name,
  },
  {
    id: "image_url",
    header: "Image",
    accessorFn: (row) => row.image_url,
    cell: (info) =>
      info.getValue() ?
        <Image
          src={info.getValue() as string}
          alt={info.row.getValue("name")}
        />
      : <span>No image</span>,
  },
  {
    id: "price",
    header: () => <div className="text-right">Price</div>,
    accessorFn: (row) => row.price,
    cell: ({ getValue }) => {
      const formatted = new Intl.NumberFormat("vi-VN", {
        style: "currency",
        currency: "VND",
      }).format(getValue() as number);

      return <div className="text-right font-medium">{formatted}</div>;
    },
  },
  {
    id: "user",
    header: "User",
    accessorFn: (row) => row.user,
    cell: ({ getValue }) => {
      const user = getValue() as GetFoodsResponseItem["user"];
      return (
        <div className="flex items-center gap-2">
          <Avatar>
            <AvatarImage src={user.avatar_url} alt={user.display_name} />
            <AvatarFallback>{user.display_name.slice(0, 2).toUpperCase()}</AvatarFallback>
          </Avatar>

          <span>{user.display_name}</span>
        </div>
      );
    },
  },
  {
    id: "description",
    header: "Description",
    accessorFn: (row) => row.description,
  },
  {
    id: "actions",
    header: "Actions",
    cell: ({ row }) => {
      // Only show actions if the current user owns this food item
      if (row.original.user.id !== currentUserId) {
        return null;
      }

      return (
        <TableActionMenuDialog
          data={{
            id: row.original.id,
            description: row.original.description,
            name: row.original.name,
            price: row.original.price,
            imageFileUrl: row.original.image_url,
          }}
        />
      );
    },
  },
];
