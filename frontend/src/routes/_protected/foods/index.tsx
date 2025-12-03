import { listFoodsOptions } from "@/client/@tanstack/react-query.gen"
import { DataTable } from "@/components/simple-data-table/data-table"
import AddFoodDialog from "@/features/foods/components/add-food-dialog"
import { createColumns } from "@/features/foods/components/columns"
import { getFoodsServer as getFoodsServerFn } from "@/features/foods/functions/get-server-foods"
import { useQuery } from "@tanstack/react-query"
import { createFileRoute } from "@tanstack/react-router"
import { useServerFn } from "@tanstack/react-start"
import { useMemo } from "react"

export const Route = createFileRoute("/_protected/foods/")({
	component: Foods,
	staticData: {
		title: "Foods",
	},
	beforeLoad: ({ context }) => {
		context.pageTitle = "Foods";
	},
	loader: async () => {
		const data = await getFoodsServerFn({ data: { query: {} } });
		return { foods: data };
	},
	staleTime: 1000 * 60 * 5, // 5 minutes
});

function Foods() {
	const { foods } = Route.useLoaderData();
	const { user } = Route.useRouteContext();
  const getFoodsServer = useServerFn(getFoodsServerFn)

	const { data } = useQuery({
		...listFoodsOptions(),
    queryFn: () => getFoodsServer({data: { query: {} }}), 
		initialData: foods,
	});

	const columns = useMemo(() => createColumns(user?.id), [user?.id]);

	return (
		<div className="flex flex-col gap-4">
			<div className="flex justify-end">
				<AddFoodDialog />
			</div>
			<div className="rounded-lg">
				<DataTable columns={columns} data={data ?? []} />
			</div>
		</div>
	);
}
