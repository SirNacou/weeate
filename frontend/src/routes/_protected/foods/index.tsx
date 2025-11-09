import { GetFoodsResponseItem } from "@/client";
import { getFoodsOptions } from "@/client/@tanstack/react-query.gen";
import { DataTable } from "@/components/data-table/data-table";
import { DataTableColumnHeader } from "@/components/data-table/data-table-column-header";
import { DataTableSortList } from "@/components/data-table/data-table-sort-list";
import { DataTableToolbar } from "@/components/data-table/data-table-toolbar";
import { useDataTable } from "@/hooks/use-data-table";
import { createFileRoute } from "@tanstack/react-router";
import { createColumnHelper } from "@tanstack/react-table";

export const Route = createFileRoute("/_protected/foods/")({
  loader: ({ context }) => {
    return context.queryClient.ensureQueryData(getFoodsOptions());
  },
  component: Foods,
});

const columnHelper = createColumnHelper<GetFoodsResponseItem>();

const columns = [
  columnHelper.accessor("name", {
    id: "name",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Name" label={"Name"} />
    ),
    cell: ({ row }) => <div>{row.getValue("name")}</div>,
    meta: {
      placeholder: "Search by Name",
      variant: "text",
    },
    enableColumnFilter: true,
  }),
  columnHelper.accessor("image_url", {
    id: "image_url",
    header: ({ column }) => (
      <DataTableColumnHeader
        column={column}
        title="Image URL"
        label={"Image URL"}
      />
    ),
    cell: ({ row }) => <div>{row.getValue("image_url")}</div>,
  }),
  columnHelper.accessor("price", {
    id: "price",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Price" label={"Price"} />
    ),
    cell: ({ row }) => (
      <div>{(row.getValue("price") as number).toFixed(0)}₫</div>
    ),
    enableColumnFilter: true,
  }),
  columnHelper.accessor("user_id", {
    id: "user",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="User" label={"User"} />
    ),
    cell: ({ row }) => <div>{row.getValue("user")}</div>,
    enableColumnFilter: true,
  }),
  columnHelper.accessor("description", {
    id: "description",
    header: ({ column }) => (
      <DataTableColumnHeader
        column={column}
        title="Description"
        label={"Description"}
      />
    ),
    cell: ({ row }) => <div>{row.getValue("description")}</div>,
    enableColumnFilter: true,
  }),
  columnHelper.display({
    id: "actions",
    header: "Actions",
  }),
];

function Foods() {
  const { table } = useDataTable({
    data: [],
    columns,
    pageCount: 0,
    initialState: {
      sorting: [{ id: "name", desc: false }],
      pagination: { pageIndex: 0, pageSize: 10 },
    },
    getRowId: (row) => row.id,
  });
  return (
    <DataTable table={table}>
      <DataTableToolbar table={table}>
        <DataTableSortList table={table} />
      </DataTableToolbar>
    </DataTable>
  );
}
