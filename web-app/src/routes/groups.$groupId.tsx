import * as React from "react";
import banIcon from "../assets/ban.svg";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Group } from "../lib/models/models";
import { UserLoggedIn } from "../lib/models/user";
import { useUserStore } from "../lib/stores";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { isRestErr, QUERY_KEYS, STYLES, toastRestErr } from "../lib/utils";
import { LoadingSpinner } from "../lib/components/Spinner";
import { toast } from "react-toastify";
import checkmarkIcon from "../assets/checkmark.svg";

export const Route = createFileRoute("/groups/$groupId")({
  component: RouteComponent,
});

function RouteComponent() {
  const { user } = useUserStore();
  const { groupId } = Route.useParams();
  const navigate = useNavigate();

  if (user.type !== "logged-in") {
    navigate({ to: "/" });
    return <LoadingSpinner content={<span>Redirecting to login...</span>} />;
  }

  return (
    <div className="flex h-full items-center justify-center">
      <GroupData
        groupId={groupId}
        user={user}
      />
    </div>
  );
}

function GroupData({ groupId, user }: { groupId: string; user: UserLoggedIn }) {
  const { setUser } = useUserStore();
  const queryClient = useQueryClient();
  const inputNameRef = React.useRef<HTMLInputElement>(null);
  const inputDescriptionRef = React.useRef<HTMLInputElement>(null);

  const {
    isPending,
    error,
    data: group,
  } = useQuery({
    queryKey: [QUERY_KEYS.groupSingle, `group-${groupId}`],
    queryFn: async () => {
      const res = await fetch(
        `${import.meta.env.VITE_API_URI}/groups/by-id/${groupId}`,
        {
          method: "GET",
          headers: {
            Authorization: user.tokens.accessToken,
            Accept: "application/json",
          },
        },
      );

      if (res.status === 401) {
        setUser({ type: "logged-out" });
      }

      if (res.status === 404) {
        return {
          type: "not-found",
        };
      }

      const data = await res.json();
      if (isRestErr(data)) {
        toastRestErr(data);
        throw new Error("Failed to load group event.");
      }

      return {
        type: "found",
        data: data as Group,
      };
    },
  });

  const updateName = useMutation({
    mutationKey: [`group-update-name-${groupId}`],
    onError: (err) => {
      console.error(err);
    },
    mutationFn: async (name: string) => {
      const res = await fetch(`${import.meta.env.VITE_API_URI}/groups/update`, {
        method: "POST",
        headers: {
          Authorization: user.tokens.accessToken,
          Accept: "application/json",
        },
        body: JSON.stringify({
          groupId,
          name,
        }),
      });

      if (res.status === 401) {
        setUser({ type: "logged-out" });
      }

      if (!res.ok) {
        const data = await res.json();
        if (isRestErr(data)) {
          toastRestErr(data);
          return;
        }
      }

      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.groupSingle, `group-${groupId}`],
      });
      toast("Updated name", { type: "success" });
    },
  });

  const updateDescription = useMutation({
    mutationKey: [`group-update-description-${groupId}`],
    onError: (err) => {
      console.error(err);
    },
    mutationFn: async (description: string) => {
      const res = await fetch(`${import.meta.env.VITE_API_URI}/groups/update`, {
        method: "POST",
        headers: {
          Authorization: user.tokens.accessToken,
          Accept: "application/json",
        },
        body: JSON.stringify({
          groupId,
          description: description.length !== 0 ? description : null,
        }),
      });

      if (res.status === 401) {
        setUser({ type: "logged-out" });
      }

      if (!res.ok) {
        const data = await res.json();
        if (isRestErr(data)) {
          toastRestErr(data);
          return;
        }
      }

      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.groupSingle, `group-${groupId}`],
      });
      toast("Updated description", { type: "success" });
    },
  });

  const joinGroup = useMutation({
    mutationKey: [`join-group-${groupId}`],
    onError: (err) => {
      console.error(err);
    },
    mutationFn: async (groupId: string) => {
      const res = await fetch(
        `${import.meta.env.VITE_API_URI}/groups/by-id/${groupId}/members/join`,
        {
          method: "POST",
          headers: {
            Authorization: user.tokens.accessToken,
            Accept: "application/json",
          },
          body: JSON.stringify({
            userId: user.id,
          }),
        },
      );

      if (res.status === 401) {
        setUser({ type: "logged-out" });
      }

      if (!res.ok) {
        const data = await res.json();
        if (isRestErr(data)) {
          toastRestErr(data);
          return;
        }
      }

      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.groupSingle, `group-${groupId}`],
      });
      toast("Approval from the group owner is pending...", { type: "info" });
    },
  });

  const leaveGroup = useMutation({
    mutationKey: [`leave-group-${groupId}`],
    onError: (err) => {
      console.error(err);
    },
    mutationFn: async (groupId: string) => {
      const res = await fetch(
        `${import.meta.env.VITE_API_URI}/groups/by-id/${groupId}/members/leave`,
        {
          method: "POST",
          headers: {
            Authorization: user.tokens.accessToken,
            Accept: "application/json",
          },
          body: JSON.stringify({
            userId: user.id,
          }),
        },
      );

      if (res.status === 401) {
        setUser({ type: "logged-out" });
      }

      if (!res.ok) {
        const data = await res.json();
        if (isRestErr(data)) {
          toastRestErr(data);
          return;
        }
      }

      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.groupSingle, `group-${groupId}`],
      });
      toast(`Left group ${group?.data?.name}`, { type: "info" });
    },
  });

  const groupApproveOrBanMember = useMutation({
    mutationKey: [`group-member-approve-or-ban-${groupId}`],
    onError: (err) => {
      console.error(err);
    },
    mutationFn: async ({
      userId,
      action,
    }: {
      userId: string;
      action: "approve" | "ban";
    }) => {
      const res = await fetch(
        `${import.meta.env.VITE_API_URI}/groups/by-id/${groupId}/members/${action}`,
        {
          method: "POST",
          headers: {
            Authorization: user.tokens.accessToken,
            Accept: "application/json",
          },
          body: JSON.stringify({
            userId,
          }),
        },
      );

      if (res.status === 401) {
        setUser({ type: "logged-out" });
      }

      if (!res.ok) {
        const data = await res.json();
        if (isRestErr(data)) {
          toastRestErr(data);
          return;
        }
      }

      queryClient.invalidateQueries({
        queryKey: [QUERY_KEYS.groupSingle, `group-${groupId}`],
      });

      const message =
        action === "approve"
          ? "Added user to group."
          : "Banned user from group.";
      toast(message, { type: "success" });
    },
  });

  if (isPending) {
    return <LoadingSpinner content={<span>Getting group...</span>} />;
  }

  if (error) {
    console.error(error);
    return <span className="text-red-500">Failed to load group</span>;
  }

  if (group.type === "not-found" || group.data === undefined) {
    return <span>No ride with the ID '{groupId}' was found</span>;
  }

  const g = group.data;
  const canEdit = g.createdBy === user.id;
  const isMember = !g.members.some((m) => m.userId === user.id);

  return (
    <div className="flex aspect-video min-w-[480px] flex-col gap-4 rounded bg-neutral-200 p-4 text-xl shadow-lg dark:bg-neutral-800 dark:shadow-none">
      <div className="flex justify-between">
        <h1 className="mb-2 text-4xl font-bold">
          Group &nbsp;
          {canEdit ? <em className="text-2xl font-normal">(owned)</em> : null}
        </h1>
        {isMember ? (
          <button
            className={STYLES.button}
            disabled={joinGroup.isPending}
            onClick={() => {
              joinGroup.mutate(g.groupId);
            }}
          >
            {joinGroup.isPending ? <LoadingSpinner /> : <>Join</>}
          </button>
        ) : (
          group.data.createdBy !== user.id && (
            <button
              className={STYLES.buttonDanger}
              disabled={leaveGroup.isPending}
              onClick={() => {
                leaveGroup.mutate(g.groupId);
              }}
            >
              {leaveGroup.isPending ? <LoadingSpinner /> : <>Leave</>}
            </button>
          )
        )}
      </div>

      <div className="flex flex-col gap-4">
        <div className="flex w-full flex-col">
          <span className="font-semibold">Name: </span>
          <input
            disabled={!canEdit}
            ref={inputNameRef}
            className="ml-2 border-b border-neutral-200 bg-transparent p-1 focus:border-cyan-500 focus:outline-none disabled:border-none dark:border-neutral-500"
            defaultValue={g.name}
            onKeyDown={(e) => {
              if (e.key !== "Enter") {
                return;
              }

              if (inputNameRef?.current) {
                updateName.mutate(inputNameRef.current.value);
              }
            }}
          />
        </div>

        <div className="flex w-full flex-col">
          <span className="font-semibold">Description: </span>
          <input
            disabled={!canEdit}
            ref={inputDescriptionRef}
            className="ml-2 border-b border-neutral-200 bg-transparent p-1 focus:border-cyan-500 focus:outline-none disabled:border-none dark:border-neutral-500"
            defaultValue={g.description ?? ""}
            onKeyDown={(e) => {
              if (e.key !== "Enter") {
                return;
              }

              if (inputDescriptionRef?.current) {
                updateDescription.mutate(inputDescriptionRef.current.value);
              }
            }}
          />
        </div>

        <div className="flex w-full flex-col">
          <span className="font-semibold">Members: </span>
          {g.members.length > 0 ? (
            g.members.map((m, idx) => {
              return (
                <div
                  key={idx}
                  className="flex w-full justify-between"
                >
                  <span
                    className={`ml-2 p-1 ${m.joinStatus === "banned" ? "text-red-500 line-through" : ""} ${m.joinStatus === "pending" ? "text-neutral-300 dark:text-neutral-600" : ""}`}
                  >
                    {m.email}
                  </span>
                  {!canEdit || m.userId === user.id ? null : (
                    <div className="flex gap-2">
                      {m.joinStatus !== "member" && (
                        <button
                          disabled={groupApproveOrBanMember.isPending}
                          onClick={() => {
                            groupApproveOrBanMember.mutate({
                              userId: m.userId,
                              action: "approve",
                            });
                          }}
                        >
                          <img
                            className="h-6 w-6"
                            src={checkmarkIcon}
                            alt="add"
                          />
                        </button>
                      )}
                      {m.joinStatus !== "banned" && (
                        <button
                          disabled={groupApproveOrBanMember.isPending}
                          onClick={() => {
                            groupApproveOrBanMember.mutate({
                              userId: m.userId,
                              action: "ban",
                            });
                          }}
                        >
                          <img
                            className="h-6 w-6"
                            src={banIcon}
                            alt="ban"
                          />
                        </button>
                      )}
                    </div>
                  )}
                </div>
              );
            })
          ) : (
            <span className="ml-2 p-1">---</span>
          )}
        </div>
      </div>
    </div>
  );
}
