import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { UserLoggedIn } from "../lib/models/user";
import { useUserStore } from "../lib/stores";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { isRestErr, STYLES, toastRestErr } from "../lib/utils";
import { LoadingSpinner } from "../lib/components/Spinner";
import { toast } from "react-toastify";

export const Route = createFileRoute("/user/$userId")({
  component: RouteComponent,
});

function RouteComponent() {
  const { user } = useUserStore();
  const { userId } = Route.useParams();
  const navigate = useNavigate();

  if (user.type !== "logged-in") {
    navigate({ to: "/" });
    return <LoadingSpinner content={<span>Redirecting to login...</span>} />;
  }

  return (
    <div className="flex h-full items-center justify-center">
      <UserData
        userId={userId}
        user={user}
      />
    </div>
  );
}

function UserData({ userId, user }: { userId: string; user: UserLoggedIn }) {
  const queryClient = useQueryClient();
  const { setUser } = useUserStore();

  const {
    isPending,
    error,
    data: userData,
  } = useQuery({
    queryKey: [`user-${userId}`],
    queryFn: async () => {
      const res = await fetch(
        `${import.meta.env.VITE_API_URI}/users/by-id/${userId}`,
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
        throw new Error("Failed to load user.");
      }

      return {
        type: "found",
        data: data as {
          id: string;
          name: string;
          email: string;
          isBlocked: boolean;
        },
      };
    },
  });

  const setBanned = useMutation({
    mutationKey: [`user-set-banned-${userId}`],
    onError: (err) => {
      console.error(err);
    },
    mutationFn: async ({
      banStatus,
    }: {
      banStatus: "banned" | "un-banned";
    }) => {
      const res = await fetch(
        `${import.meta.env.VITE_API_URI}/users/by-id/${userId}/ban-status`,
        {
          method: "POST",
          headers: {
            Authorization: user.tokens.accessToken,
            Accept: "application/json",
          },
          body: JSON.stringify({
            isBanned: banStatus === "banned" ? true : false,
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
        queryKey: [`user-${userId}`],
      });

      const message =
        banStatus === "banned" ? "Banned user." : "Un-Banned user.";
      toast(message, { type: "success" });
    },
  });

  if (isPending) {
    return <LoadingSpinner content={<span>Getting user...</span>} />;
  }

  if (error) {
    console.error(error);
    return <span className="text-red-500">Failed to load user</span>;
  }

  if (userData.type === "not-found" || userData.data === undefined) {
    return <span>No user with the ID '{userId}' was found</span>;
  }

  const u = userData.data;

  return (
    <div className="flex aspect-video min-w-[480px] flex-col gap-4 rounded bg-neutral-200 p-4 text-xl shadow-lg dark:bg-neutral-800 dark:shadow-none">
      <div className="flex justify-between">
        <h1 className="mb-2 text-4xl font-bold">User &nbsp;</h1>
        {user.isAdmin && u.id !== user.id ? (
          u.isBlocked ? (
            <button
              disabled={setBanned.isPending}
              className={STYLES.button}
              onClick={() => {
                setBanned.mutate({ banStatus: "un-banned" });
              }}
            >
              Un-Ban
            </button>
          ) : (
            <button
              disabled={setBanned.isPending}
              className={STYLES.buttonDanger}
              onClick={() => {
                setBanned.mutate({ banStatus: "banned" });
              }}
            >
              Ban
            </button>
          )
        ) : null}
      </div>

      <div className="flex flex-col gap-4">
        <div className="flex w-full flex-col">
          <span className="font-semibold">Name: </span>
          <span>{u.name}</span>
        </div>

        <div className="flex w-full flex-col">
          <span className="font-semibold">Email: </span>
          <span>{u.email}</span>
        </div>
      </div>
    </div>
  );
}
