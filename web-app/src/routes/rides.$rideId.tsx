import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { LoadingSpinner } from "../lib/components/Spinner";
import { useUserStore } from "../lib/stores";
import { useQuery } from "@tanstack/react-query";
import { QUERY_KEYS } from "../lib/query";
import { Ride } from "../lib/models/ride";
import { AuthTokens } from "../lib/models/user";

export const Route = createFileRoute("/rides/$rideId")({
  component: RouteComponent,
});

function RouteComponent() {
  const { user } = useUserStore();
  const { rideId } = Route.useParams();
  const navigate = useNavigate();

  if (user.type !== "logged-in") {
    navigate({ to: "/" });
    return <LoadingSpinner content={<span>Redirecting to login...</span>} />;
  }

  return (
    <div className="flex h-full items-center justify-center">
      <RideData
        rideId={rideId}
        tokens={user.tokens}
      />
    </div>
  );
}

function RideData({ tokens, rideId }: { tokens: AuthTokens; rideId: string }) {
  const {
    isPending,
    error,
    data: ride,
  } = useQuery({
    queryKey: [QUERY_KEYS.rideSingle],
    queryFn: async () => {
      const res = await fetch(
        `${import.meta.env.VITE_API_URI}/rides/by-id/${rideId}`,
        {
          method: "GET",
          headers: {
            Authorization: tokens.accessToken,
            Accept: "application/json",
          },
        },
      );

      if (res.status === 404) {
        return {
          type: "not-found",
        };
      }

      const ride = await res.json();
      return {
        type: "found",
        data: ride as Ride,
      };
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
  return (
    <div className="relative flex aspect-video min-w-[320px] flex-col gap-2 rounded bg-neutral-200 p-4 shadow-lg dark:bg-neutral-800 dark:shadow-none">
      <h1 className="absolute left-0 top-0 translate-x-[-10%] translate-y-[-60%] text-3xl font-bold">
        Ride
      </h1>
      <span>Driver: {r.driverEmail}</span>
      <span>To: {r.locationTo}</span>
      <span>From: {r.locationFrom}</span>
      <span>When: {new Date(r.tackingPlaceAt).toLocaleString()}</span>
      <span>Created By: {r.createdByEmail}</span>
      <span>Created At: {new Date(r.createdAt).toLocaleString()}</span>
    </div>
  );
}
