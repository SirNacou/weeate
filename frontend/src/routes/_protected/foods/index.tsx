import { getFoodsOptions } from "@/client/@tanstack/react-query.gen";
import { createFileRoute } from "@tanstack/react-router";

import { useQuery } from "@tanstack/react-query";
import { createColumns } from "@/features/foods/components/columns";
import { DataTable } from "@/components/simple-data-table/data-table";
import AddFoodDialog from "@/features/foods/components/add-food-dialog";
import { useMemo } from "react";
import { getServerFoods } from "@/features/foods/functions/get-server-foods";

export const Route = createFileRoute("/_protected/foods/")({
  component: Foods,
  loader: async () => ({ initialData: await getServerFoods() }),
});

function Foods() {
  const { initialData } = Route.useLoaderData();
  const { user } = Route.useRouteContext();
  const { data } = useQuery({ ...getFoodsOptions(), initialData });

  const columns = useMemo(() => createColumns(user?.id), [user?.id]);

  return (
    <div className="flex flex-col gap-4">
      <div className="flex justify-end">
        <AddFoodDialog />
      </div>
      <DataTable columns={columns} data={data?.result || []} />
    </div>
  );
}
