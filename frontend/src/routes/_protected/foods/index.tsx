import { getFoodsOptions } from "@/client/@tanstack/react-query.gen";
import { createFileRoute } from "@tanstack/react-router";

import { useQuery } from "@tanstack/react-query";
import { columns } from "@/features/foods/components/columns";
import { DataTable } from "@/components/simple-data-table/data-table";
import { GetFoodsResponseItem } from "@/client/types.gen";
import AddFoodDialog from "@/features/foods/components/add-food-dialog";

export const Route = createFileRoute("/_protected/foods/")({
  loader: ({ context }) => {
    return context.queryClient.ensureQueryData(getFoodsOptions());
  },
  component: Foods,
});

const data: GetFoodsResponseItem[] = [
  {
    id: "1",
    name: "Pho Bo",
    image_url: "https://example.com/images/pho_bo.jpg",
    price: 50000,
    user_id: "user_123",
    description: "Traditional Vietnamese beef noodle soup.",
  },
];

function Foods() {
  const initialData = Route.useLoaderData();
  // const { data } = useQuery({
  //   ...getFoodsOptions(),
  //   initialData,
  // });

  return (
    <div className="container mx-auto py-10 flex flex-col gap-4">
      <div className="flex justify-end">
        <AddFoodDialog />
      </div>
      <DataTable columns={columns} data={data} />
    </div>
  );
}
