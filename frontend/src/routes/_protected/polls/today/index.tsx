import { GetTodayPollsQueryResponse, Vote } from "@/client";
import {
  listPollsTodayOptions,
  listPollsTodayQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { listPollsTodayServerFn } from "@/features/polls/api/get-today-polls";
import CreatePollDialog from "@/features/polls/components/create-poll-dialog";
import PollCard from "@/features/polls/components/poll-card";
import { useSubscription } from "@/lib/centrifugo/use-subscription";
import { getPageTitle } from "@/lib/head-utils";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { toast } from "sonner";

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

  const { data, refetch } = useQuery({
    ...listPollsTodayOptions(),
    queryFn: listPollsTodayServerFn,
    initialData: polls,
  });

  const queryClient = useQueryClient();

  useSubscription("public:polls", async (data) => {
    console.log("Vote moved event received", data);
    if (data.type === "poll_created") {
      toast("New poll created, refetching today's polls...");
      await refetch();
    }
    if (data.type === "vote_added") {
      queryClient.setQueryData(
        listPollsTodayQueryKey(),
        (oldData: GetTodayPollsQueryResponse[] | null) => {
          if (!oldData) return oldData;

          return oldData.map((poll) => {
            if (poll.id !== data.data.poll_id) return poll;
            return {
              ...poll,
              poll_options: poll.poll_options?.map((option) => {
                if (option.id !== data.data.option_id) return option;

                if (
                  option.votes?.some((v) => v.voter?.id === data.data.user_id)
                ) {
                  return option;
                }

                return {
                  ...option,
                  votes: [
                    ...(option.votes || []),
                    {
                      voter: {
                        id: data.data.user_id,
                        display_name: data.data.user_display_name,
                        avatar_url: data.data.user_avatar_url,
                        created_at: new Date().toISOString(),
                      },
                    },
                  ] as Array<Vote>,
                };
              }),
            };
          });
        }
      );
    }
    if (data.type === "vote_removed") {
      queryClient.setQueryData(
        listPollsTodayQueryKey(),
        (oldData: GetTodayPollsQueryResponse[] | null) => {
          if (!oldData) return oldData;

          return oldData.map((poll) => {
            if (poll.id !== data.data.poll_id) return poll;
            return {
              ...poll,
              poll_options: poll.poll_options?.map((option) => {
                if (option.id !== data.data.option_id) return option;
                return {
                  ...option,
                  votes: option.votes?.filter(
                    (v) => v.voter?.id !== data.data.user_id
                  ),
                };
              }),
            };
          });
        }
      );
    }
    if (data.type === "vote_moved") {
      queryClient.setQueryData(
        listPollsTodayQueryKey(),
        (oldData: GetTodayPollsQueryResponse[] | null) => {
          if (!oldData) return oldData;

          return oldData.map((poll) => {
            if (poll.id !== data.data.poll_id) return poll;
            return {
              ...poll,
              poll_options: poll.poll_options?.map((option) => {
                // Remove the user's vote from all options first to ensure no duplicates
                const votesWithoutUser = option.votes?.filter(
                  (v) => String(v.voter?.id) !== String(data.data.user_id)
                );

                if (option.id === data.data.new_option_id) {
                  // Check if the user is already in the votes list (in case of race conditions or duplicate events)
                  const isAlreadyAdded = votesWithoutUser?.some(
                    (v) => v.voter?.id === data.data.user_id
                  );

                  if (isAlreadyAdded) {
                    return {
                      ...option,
                      votes: votesWithoutUser,
                    };
                  }

                  return {
                    ...option,
                    votes: [
                      ...(votesWithoutUser || []),
                      {
                        voter: {
                          id: data.data.user_id,
                          display_name: data.data.user_display_name,
                          avatar_url: data.data.user_avatar_url,
                          created_at: new Date().toISOString(),
                        },
                      },
                    ] as Array<Vote>,
                  };
                }

                return {
                  ...option,
                  votes: votesWithoutUser,
                };
              }),
            };
          });
        }
      );
    }
  });

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-3">
        <div className="flex flex-col sm:flex-row justify-between items-center gap-3">
          <div className="flex flex-col self-start">
            <h1 className="text-2xl font-bold">Today Food Polls</h1>
            <p>Vote for your favorite meals and let your voice be heard!</p>
          </div>

          <CreatePollDialog
            userPollExists={data?.some((poll) => poll.creator?.id === user?.id)}
          />
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
                {data?.reduce((acc, poll) => {
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
        {data?.map((poll) => {
          return (
            <PollCard
              key={poll.id}
              pollId={poll.id}
              avatarUrl={poll.creator?.avatar_url}
              buyerName={poll.creator?.display_name || "Unknown"}
              scheduled_close_at={poll.scheduled_closes_at}
              closed_at={poll.closed_at}
              strategy={poll.strategy as any}
              options={
                poll.poll_options?.map((option) => ({
                  id: option.id || "",
                  foodName: option.food?.name || "Unknown Food",
                  foodImageUrl: "",
                  price: option.price_at_creation || 0,
                  votes:
                    option.votes?.map((vote) => ({
                      userId: vote.voter?.id || "",
                      userName: vote.voter?.display_name || "Unknown",
                      userAvatarUrl: vote.voter?.avatar_url,
                    })) || [],
                })) || []
              }
            />
          );
        })}
      </div>
    </div>
  );
}
