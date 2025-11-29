import { GetTodayPollsQueryResponse, Vote } from "@/client";
import { listPollsTodayQueryKey } from "@/client/@tanstack/react-query.gen";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { listPollsTodayServerFn } from "@/features/polls/api/get-today-polls";
import CreatePollDialog from "@/features/polls/components/create-poll-dialog";
import PollCard from "@/features/polls/components/poll-card";
import { CentrifugoProvider } from "@/lib/centrifugo/centrifugo-context";
import { useSubscription } from "@/lib/centrifugo/use-subscription";
import { QueryClient, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useServerFn } from "@tanstack/react-start";
import { toast } from "sonner";

export const Route = createFileRoute("/_protected/polls/today/")({
  component: RouteComponent,
  beforeLoad: () => {
    return {
      pageTitle: "Today's Polls",
    };
  },
  loader: async ({ abortController }) => {
    const res = await listPollsTodayServerFn({
      signal: abortController.signal,
    });
    return { polls: res.data || [], error: res.error };
  },
});

function RouteComponent() {
  return (
    <CentrifugoProvider>
      <PollsContent />
    </CentrifugoProvider>
  );
}

function PollsContent() {
  const { polls } = Route.useLoaderData();
  const { user } = Route.useRouteContext();
  console.log("polls loader data", polls);

  const queryClient = useQueryClient();

  // Wrap server function for client-side use
  const fetchPolls = useServerFn(listPollsTodayServerFn);

  const { data } = useQuery({
    queryKey: listPollsTodayQueryKey(),
    queryFn: async () => {
      const res = await fetchPolls();
      if (res.error) {
        toast.error(res.error);
        throw new Error(res.error);
      }
      return res.data || [];
    },
    initialData: polls,
    staleTime: 30000, // Consider data fresh for 30s to prevent immediate refetch
  });

  useSubscription("public:polls", updatePollsPage(queryClient));

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-3">
        <div className="flex w-full flex-col items-center justify-end gap-3 sm:flex-row">
          <CreatePollDialog
            userPoll={data.find(
              (poll) =>
                poll.creator?.id === user?.id &&
                new Date(poll.order_date).toISOString().split("T")[0] ===
                  new Date().toISOString().split("T")[0]
            )}
          />
        </div>

        <div className="flex gap-2 text-center">
          <Card className="basis-1/2 md:basis-[200px] gap-2 py-4">
            <CardHeader>
              <CardTitle>Polls</CardTitle>
            </CardHeader>
            <CardContent>
              <span className="font-bold">{data?.length || 0}</span>
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
          return <PollCard key={poll.id} poll={poll} />;
        })}
      </div>
    </div>
  );
}
function updatePollsPage(queryClient: QueryClient): (data: any) => void {
  return (data) => {
    console.log("Vote moved event received", data);
    if (data.type === "poll_created") {
      toast("New poll created, refetching today's polls...");
      queryClient.invalidateQueries({ queryKey: listPollsTodayQueryKey() });
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
  };
}
