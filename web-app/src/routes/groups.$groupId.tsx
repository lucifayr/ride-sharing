import * as React from "react";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Group } from "../lib/models/models";
import { UserLoggedIn } from "../lib/models/user";
import { useUserStore } from "../lib/stores";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { isRestErr, QUERY_KEYS, toastRestErr } from "../lib/utils";
import { LoadingSpinner } from "../lib/components/Spinner";
import { toast } from "react-toastify";

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
    mutationKey: ["group-update-name"],
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
        queryKey: [QUERY_KEYS.groupSingle, QUERY_KEYS.groupItems],
      });
      toast("Updated name", { type: "success" });
    },
  });

  const updateDescription = useMutation({
    mutationKey: ["group-update-description"],
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
        queryKey: [QUERY_KEYS.groupSingle, QUERY_KEYS.groupItems],
      });
      toast("Updated description", { type: "success" });
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

  return (
    <div className="flex aspect-video min-w-[320px] flex-col gap-4 rounded bg-neutral-200 p-4 text-xl shadow-lg dark:bg-neutral-800 dark:shadow-none">
      <h1 className="mb-2 text-4xl font-bold">
        Group &nbsp;
        {canEdit ? <em className="text-2xl font-normal">(owned)</em> : null}
      </h1>

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
      </div>
    </div>
  );
}
