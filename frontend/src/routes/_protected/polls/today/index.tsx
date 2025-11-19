import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { getPageTitle } from "@/lib/head-utils";
import { createFileRoute } from "@tanstack/react-router";
import LucidePlus from "~icons/lucide/plus?width=2em&height=2em";
import LucideCalendar from "~icons/lucide/calendar";
import LucideUserRound from "~icons/lucide/user-round";
import { Badge } from "@/components/ui/badge";
import CloseTimer from "@/components/close-timer";
import { add } from "date-fns";
import PollStrategyBadge from "@/features/polls/components/poll-strategy-badge";
import PollCard from "@/features/polls/components/poll-card";

export const Route = createFileRoute("/_protected/polls/today/")({
  component: RouteComponent,
  head: () => {
    return {
      meta: [
        {
          title: getPageTitle("Today Polls"),
        },
      ],
    };
  },
});

function RouteComponent() {
  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-3">
        <div className="flex flex-col sm:flex-row justify-between items-center gap-3">
          <div className="flex flex-col self-start">
            <h1 className="text-2xl font-bold">Today Food Polls</h1>
            <p>Vote for your favorite meals and let your voice be heard!</p>
          </div>

          <Button className="w-full sm:w-auto">
            <LucidePlus />
            Create New
          </Button>
        </div>

        <div className="flex gap-2 text-center">
          <Card className="basis-1/2 md:basis-[200px] gap-2 py-4">
            <CardHeader>
              <CardTitle>Today Polls</CardTitle>
            </CardHeader>
            <CardContent>
              <span className="font-bold">10</span>
            </CardContent>
          </Card>
          <Card className="basis-1/2 md:basis-[200px] gap-2 py-4">
            <CardHeader>
              <CardTitle>Your Votes</CardTitle>
            </CardHeader>
            <CardContent>
              <span className="font-bold">5</span>
            </CardContent>
          </Card>
        </div>
      </div>

      <div className="grid grid-cols-1 xl:grid-cols-2 gap-6">
        <PollCard
          avatarUrl=""
          buyerName="Anh"
          strategy="ORDER_CONSENSUS_ITEM"
          options={[
            {
              id: "1",
              foodName: "Pizza",
              foodImageUrl: "",
              price: 100000,
            },
            {
              id: "2",
              foodName: "Burger",
              foodImageUrl: "",
              price: 80000,
            },
          ]}
        />
        <Card>
          <CardHeader>
            <CardTitle className="flex gap-1 items-center text-lg">
              <LucideCalendar />
              <span>
                <b>Tomorrow's Breakfast</b> - Vote now!
              </span>

              <Badge className="text-sm">Eat what you pick</Badge>
            </CardTitle>
            <CardDescription className="flex gap-1 items-center text-base">
              <LucideUserRound />
              <span>
                <b>User</b> will buy the winning option for <b>everyone</b>
              </span>
            </CardDescription>
          </CardHeader>
          <CardContent>
            {/* Replace with actual poll list */}
            <p>No active polls at the moment.</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex gap-2 items-center text-lg">
              <div className="flex gap-1">
                <LucideCalendar />
                <span>
                  <b>Tomorrow's Breakfast</b> - Vote now!
                </span>
              </div>

              <PollStrategyBadge strategy="ORDER_CONSENSUS_ITEM" />

              <CloseTimer
                className="text-sm ml-auto px-3"
                closesAt={add(Date.now(), { minutes: 120 })}
              />
            </CardTitle>
            <CardDescription className="flex gap-1 items-center text-base">
              <LucideUserRound />
              <span>
                <b>User</b> will buy the winning option for <b>everyone</b>
              </span>
            </CardDescription>
          </CardHeader>
          <CardContent>
            {/* Replace with actual poll list */}
            <p>No active polls at the moment.</p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
