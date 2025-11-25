import { listFoodsOptions } from "@/client/@tanstack/react-query.gen";
import { createFileRoute } from "@tanstack/react-router";

import { DataTable } from "@/components/simple-data-table/data-table";
import AddFoodDialog from "@/features/foods/components/add-food-dialog";
import { createColumns } from "@/features/foods/components/columns";
import { getFoodsServer } from "@/features/foods/functions/get-server-foods";
import { useQuery } from "@tanstack/react-query";
import { useMemo } from "react";

export const Route = createFileRoute("/_protected/foods/")({
  component: Foods,
  beforeLoad: ({ context }) => {
    context.pageTitle = "Foods";
  },
  loader: async () => ({
    initialData: await getFoodsServer({ data: { query: {} } }),
  }),
});

function Foods() {
  const { initialData } = Route.useLoaderData();
  const { user } = Route.useRouteContext();
  const { data } = useQuery({ ...listFoodsOptions(), initialData });

  const columns = useMemo(() => createColumns(user?.id), [user?.id]);

  return (
    <div className="flex flex-col gap-4">
      <div className="flex justify-end">
        <AddFoodDialog />
      </div>
      <DataTable columns={columns} data={data || []} />
    </div>
  );
}
