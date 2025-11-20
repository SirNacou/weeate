import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { getPageTitle } from "@/lib/head-utils";
import { createFileRoute } from "@tanstack/react-router";
import LucideCalendar from "~icons/lucide/calendar";
import LucideUserRound from "~icons/lucide/user-round";
import { Badge } from "@/components/ui/badge";
import PollCard from "@/features/polls/components/poll-card";
import { listPollsTodayServerFn } from "@/features/polls/api/get-today-polls";
import CreatePollDialog from "@/features/polls/components/create-poll-dialog";

export const Route = createFileRoute("/_protected/polls/today/")({
  component: RouteComponent,
  loader: async () => {
    const polls = await listPollsTodayServerFn();
    return { polls };
  },
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
  const { polls } = Route.useLoaderData();
  const { user } = Route.useRouteContext();

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-3">
        <div className="flex flex-col sm:flex-row justify-between items-center gap-3">
          <div className="flex flex-col self-start">
            <h1 className="text-2xl font-bold">Today Food Polls</h1>
            <p>Vote for your favorite meals and let your voice be heard!</p>
          </div>

          <CreatePollDialog />
        </div>

        <div className="flex gap-2 text-center">
          <Card className="basis-1/2 md:basis-[200px] gap-2 py-4">
            <CardHeader>
              <CardTitle>Today Polls</CardTitle>
            </CardHeader>
            <CardContent>
              <span className="font-bold">{polls?.length || 0}</span>
            </CardContent>
          </Card>
          <Card className="basis-1/2 md:basis-[200px] gap-2 py-4">
            <CardHeader>
              <CardTitle>Your Votes</CardTitle>
            </CardHeader>
            <CardContent>
              <span className="font-bold">
                {polls?.reduce((acc, poll) => {
                  const hasVoted = poll.poll_options?.some((option) =>
                    option.votes?.some((vote) => vote.voter?.id === user?.id)
                  );
                  return acc + (hasVoted ? 1 : 0);
                }, 0) || 0}
              </span>
            </CardContent>
          </Card>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6">
        {polls?.map((poll) => {
          const myVote = poll.poll_options?.find((option) =>
            option.votes?.some((vote) => vote.voter?.id === user?.id)
          );

          return (
            <PollCard
              key={poll.id}
              avatarUrl={poll.creator?.avatar_url || ""}
              buyerName={poll.creator?.display_name || "Unknown"}
              strategy={poll.strategy as any}
              initialSelectedOption={myVote?.id}
              options={
                poll.poll_options?.map((option) => ({
                  id: option.id || "",
                  foodName: option.food?.name || "Unknown Food",
                  foodImageUrl: "", // API doesn't seem to return image url in poll option yet
                  price: option.price_at_creation || 0,
                  votes:
                    option.votes?.map((vote) => ({
                      userId: vote.voter?.id || "",
                      userName: vote.voter?.display_name || "Unknown",
                      userAvatarUrl: vote.voter?.avatar_url || "",
                    })) || [],
                })) || []
              }
            />
          );
        })}

        {(!polls || polls.length === 0) && (
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
              <p>No active polls at the moment.</p>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
