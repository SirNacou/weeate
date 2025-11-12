import { getFoodsOptions } from "@/client/@tanstack/react-query.gen";
import { createFileRoute } from "@tanstack/react-router";

import { useQuery } from "@tanstack/react-query";
import { columns } from "@/features/foods/components/columns";
import { DataTable } from "@/components/simple-data-table/data-table";
import AddFoodDialog from "@/features/foods/components/add-food-dialog";
import { buildSrc } from '@imagekit/react'

export const Route = createFileRoute("/_protected/foods/")({
  component: Foods,
});

function Foods() {
  const { data } = useQuery(getFoodsOptions());
  

  return (
    <div className="flex flex-col gap-4">
      <div className="flex justify-end">
        <AddFoodDialog />
      </div>
      <DataTable columns={columns} data={data?.result || []} />
    </div>
  );
}
