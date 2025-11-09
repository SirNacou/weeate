import { DataGrid } from "@/components/data-grid/data-grid";
import { useDataGrid } from "@/hooks/use-data-grid";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_protected/foods/")({
  component: Foods,
});

function Foods() {
  const { table } = useDataGrid({
    data: [],
    columns: []
  });
  return <DataGrid table={table}></DataGrid>;
}
