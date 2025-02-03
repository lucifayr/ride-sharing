import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { LoadingSpinner } from "../lib/components/Spinner";
import { useUserStore } from "../lib/stores";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { isRestErr, QUERY_KEYS, STYLES, toastRestErr } from "../lib/utils";
import { RideEvent, RideSchedule } from "../lib/models/ride";
import { UserLoggedIn } from "../lib/models/user";
import { displaySchedule } from "./dashboard";
import { toast } from "react-toastify";
import { parseRecuring } from "../lib/components/CreateRideForm";
import { useEffect, useRef } from "react";

export const Route = createFileRoute("/rides/$rideId")({
  component: RouteComponent,
});

function RouteComponent() {
  const { user } = useUserStore();
  const { rideId } = Route.useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  if (user.type !== "logged-in") {
    navigate({ to: "/" });
    return <LoadingSpinner content={<span>Redirecting to login...</span>} />;
  }

  return (
    <div className="flex h-full items-center justify-center">
      <RideData
        rideId={rideId}
        user={user}
      />
    </div>
  );
}

function RideData({ user, rideId }: { user: UserLoggedIn; rideId: string }) {
  const { setUser } = useUserStore();
  const queryClient = useQueryClient();
  const inputScheduleRef = useRef<HTMLInputElement>(null);

  const {
    isPending,
    error,
    data: ride,
  } = useQuery({
    queryKey: [QUERY_KEYS.rideSingle, `ride-${rideId}`],
    queryFn: async () => {
      const res = await fetch(
        `${import.meta.env.VITE_API_URI}/rides/by-id/${rideId}`,
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
        throw new Error("Failed to load ride event.");
      }

      return {
        type: "found",
        data: data as RideEvent,
      };
    },
  });

  const updateSchedule = useMutation({
    mutationKey: ["update-ride-schedule"],
    onError: (err) => {
      console.error(err);
    },
    mutationFn: async ({
      schedule,
      rideEventId,
    }: {
      schedule: RideSchedule;
      rideEventId: string;
    }) => {
      const res = await fetch(`${import.meta.env.VITE_API_URI}/rides/update`, {
        method: "POST",
        headers: {
          Authorization: user.tokens.accessToken,
          Accept: "application/json",
        },
        body: JSON.stringify({
          schedule,
          rideEventId,
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

      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.rideSingle] });
      toast("Update ride schedule", { type: "success" });
    },
  });

  const joinRide = useMutation({
    mutationKey: ["join-ride"],
    onError: (err) => {
      console.error(err);
    },
    mutationFn: async (rideEventId: string) => {
      const res = await fetch(`${import.meta.env.VITE_API_URI}/rides/join`, {
        method: "POST",
        headers: {
          Authorization: user.tokens.accessToken,
          Accept: "application/json",
        },
        body: JSON.stringify({
          userId: user.id,
          rideEventId,
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

      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.rideSingle] });
      toast("Joined ride", { type: "success" });
    },
  });

  const cancelRide = useMutation({
    mutationKey: ["cancel-ride"],
    onError: (err) => {
      console.error(err);
    },
    mutationFn: async (rideEventId: string) => {
      const res = await fetch(`${import.meta.env.VITE_API_URI}/rides/update`, {
        method: "POST",
        headers: {
          Authorization: user.tokens.accessToken,
          Accept: "application/json",
        },
        body: JSON.stringify({
          rideEventId,
          status: "canceled",
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

      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.rideSingle] });
      toast("Canceled ride", { type: "success" });
    },
  });

  if (isPending) {
    return <LoadingSpinner content={<span>Getting ride...</span>} />;
  }

  if (error) {
    console.error(error);
    return <span className="text-red-500">Failed to load ride</span>;
  }

  if (ride.type === "not-found" || ride.data === undefined) {
    return <span>No ride with the ID '{rideId}' was found</span>;
  }

  const r = ride.data;
  const canEdit = r.createdBy === user.id;
  const canJoin = !r.participants.some((p) => p.userId === user.id);

  return (
    <div className="flex aspect-video min-w-[320px] flex-col gap-4 rounded bg-neutral-200 p-4 text-xl shadow-lg dark:bg-neutral-800 dark:shadow-none">
      <div className="flex justify-between">
        <h1 className="mb-2 text-4xl font-bold">
          Ride &nbsp;
          {canEdit ? <em className="text-2xl font-normal">(owned)</em> : null}
        </h1>
        {canJoin && (
          <button
            className={STYLES.button}
            disabled={joinRide.isPending}
            onClick={() => {
              joinRide.mutate(r.rideEventId);
            }}
          >
            {joinRide.isPending ? <LoadingSpinner /> : <>Join</>}
          </button>
        )}
      </div>

      <div className="flex">
        <div className="flex w-1/2 flex-col">
          <span className="font-semibold">To: </span>
          <span className="ml-2 p-1">{r.locationTo}</span>
        </div>
        <div className="flex w-1/2 flex-col">
          <span className="font-semibold">From: </span>
          <span className="ml-2 p-1">{r.locationFrom}</span>
        </div>
      </div>

      <div className="flex flex-col">
        <span className="font-semibold">When: </span>
        <span className="ml-2 p-1">
          {new Date(r.tackingPlaceAt).toLocaleString()}
        </span>
      </div>

      <div className="flex flex-col">
        <span className="font-semibold">Status: </span>
        <div className="flex justify-between">
          <span className="ml-2 p-1">{r.status}</span>
          {r.status === "upcoming" ? (
            <button
              className="font-bold text-red-500"
              onClick={() => {
                cancelRide.mutate(r.rideEventId);
              }}
            >
              X
            </button>
          ) : null}
        </div>
      </div>

      <div className="flex flex-col">
        <span className="font-semibold">Driver: </span>
        <span className="ml-2 p-1">{r.driverEmail}</span>
      </div>

      <div className="flex flex-col">
        <span className="font-semibold">Recurs: </span>
        <input
          disabled={!canEdit}
          ref={inputScheduleRef}
          className="ml-2 border-b border-neutral-200 bg-transparent p-1 focus:border-cyan-500 focus:outline-none disabled:border-none dark:border-neutral-500"
          defaultValue={displaySchedule(r.schedule)}
          onKeyDown={(e) => {
            if (e.key !== "Enter") {
              return;
            }

            const schedule = parseRecuring(inputScheduleRef.current!.value);
            if (!schedule) {
              toast("Invalid ride schedule.", { type: "warning" });
              return;
            }

            updateSchedule.mutate({ schedule, rideEventId: r.rideEventId });
          }}
        />
      </div>

      <div className="flex flex-col">
        <span className="font-semibold">Participants: </span>
        {r.participants.length > 0 ? (
          r.participants.map((p) => {
            return <span className="ml-2 p-1">{p.email}</span>;
          })
        ) : (
          <span className="ml-2 p-1">---</span>
        )}
      </div>
    </div>
  );
}
