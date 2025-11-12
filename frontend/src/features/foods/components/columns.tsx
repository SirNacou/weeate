import { GetFoodsResponseItem } from "@/client";
import { ColumnDef } from "@tanstack/react-table";
import { Image } from "@imagekit/react";

export const columns: ColumnDef<GetFoodsResponseItem>[] = [
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
    id: "user_id",
    header: "User ID",
    accessorFn: (row) => row.user_id,
  },
  {
    id: "description",
    header: "Description",
    accessorFn: (row) => row.description,
  },
];
