import { serverClient } from "@/api";
import { getOrderedItems, getShoppingOrder } from "@/client";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from "@/components/ui/empty";
import MealCard from "@/features/orders/components/meal-card";
import ShoppingCard from "@/features/orders/components/shopping-card";
import { createFileRoute, Link, redirect } from "@tanstack/react-router";
import { createServerFn } from "@tanstack/react-start";
import { useEffect } from "react";
import z from "zod";
import LucideSandwich from "~icons/lucide/sandwich";

const getOrderedItemsServerFn = createServerFn({ method: "GET" })
  .inputValidator(z.object({ date: z.date() }))
  .handler(async ({ data }) => {
    const resp = await getOrderedItems({
      client: serverClient,
      query: {
        date: data.date.toISOString(),
      },
    });
    if (resp.error) {
      return {
        error: resp.error.errors?.at(0)?.message ?? "No ordered items found",
      };
    }
    return {
      data: resp.data,
    };
  });

const getShoppingOrderServerFn = createServerFn({ method: "GET" })
  .inputValidator(z.object({ date: z.date() }))
  .handler(async ({ data }) => {
    const resp = await getShoppingOrder({
      client: serverClient,
      query: {
        date: data.date.toISOString(),
      },
    });
    if (resp.error) {
      return {
        error: resp.error.errors?.at(0)?.message ?? "No shopping order found",
      };
    }
    return {
      data: resp.data,
    };
  });

export const Route = createFileRoute("/_protected/")({
  component: Index,
  staticData: {
    title: "Home",
  },
  beforeLoad: () => {
    return {
      pageTitle: "Home",
    };
  },
  validateSearch: z.object({
    date: z.iso.date().default(new Date().toISOString().split("T")[0]),
  }),
  loaderDeps: ({ search }) => ({ search }),
  loader: async ({ deps: { search } }) => {
    // Auto-navigate to today's date if no date provided
    const today = new Date().toISOString().split("T")[0];
    if (search.date !== today) {
      throw redirect({
        to: "/",
        search: { date: today },
        replace: true,
      });
    }

    const [orderedItemsRes, shoppingOrderRes] = await Promise.all([
      getOrderedItemsServerFn({
        data: {
          date: new Date(search.date),
        },
      }),
      getShoppingOrderServerFn({
        data: {
          date: new Date(search.date),
        },
      }),
    ]);

    return {
      orderedItems: orderedItemsRes.data?.items || [],
      shoppingOrder: shoppingOrderRes.data,
      error: orderedItemsRes.error || shoppingOrderRes.error,
    };
  },
});

function Index() {
  const { orderedItems, shoppingOrder, error } = Route.useLoaderData();

  useEffect(() => {
    if (error) {
      console.error("Error loading data:", error);
    }
  }, [error]);

  return (
    <div className="flex flex-col gap-6">
      {orderedItems.length === 0 ?
        <Card>
          <Empty>
            <EmptyHeader>
              <EmptyMedia variant={"icon"}>
                <LucideSandwich className="size-20" />
              </EmptyMedia>
              <EmptyTitle>No meals ordered yet</EmptyTitle>
              <EmptyDescription>
                You have not ordered any meals yet. Start voting to get meals!
              </EmptyDescription>
            </EmptyHeader>
            <EmptyContent>
              <Button asChild>
                <Link to="/polls/today">Go to voting page</Link>
              </Button>
            </EmptyContent>
          </Empty>
        </Card>
      : <MealCard orderedItems={orderedItems} />}

      {shoppingOrder && shoppingOrder.items.length > 0 && (
        <ShoppingCard shoppingOrder={shoppingOrder} />
      )}
    </div>
  );
}
