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
import { createFileRoute, Link } from "@tanstack/react-router";
import LucideSandwich from "~icons/lucide/sandwich";

export const Route = createFileRoute("/_protected/")({
  component: Index,
  beforeLoad: () => {
    return {
      pageTitle: "Home",
    };
  },
});

function Index() {
  const empty = false;
  return (
    <div className="flex flex-col gap-6">
      {empty ?
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
      : <MealCard />}

      <ShoppingCard />
    </div>
  );
}
