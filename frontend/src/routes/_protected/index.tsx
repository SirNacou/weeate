import { DataTable } from "@/components/data-table/data-table";
import { Card } from "@/components/ui/card";
import { createFileRoute } from "@tanstack/react-router";
import { getCoreRowModel, useReactTable } from "@tanstack/react-table";

export const Route = createFileRoute("/_protected/")({
  component: Index,
  beforeLoad: () => {
    return {
      pageTitle: "Home",
    };
  },
});

function Index() {
  const table = useReactTable({
    data: [],
    columns: [],
    getCoreRowModel: getCoreRowModel(),
  });
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
      <Card></Card>

      <Card></Card>

      <DataTable className="md:col-span-2" table={table}></DataTable>
    </div>
  );
}
