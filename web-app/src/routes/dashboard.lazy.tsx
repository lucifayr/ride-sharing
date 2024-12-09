import { createLazyFileRoute, useNavigate } from "@tanstack/react-router";
import { useUserStore } from "../lib/stores";
import { LoadingSpinner } from "../lib/components/Spinner";
import { CreateRideForm } from "../lib/components/CreateRideForm";
import { useQuery } from "@tanstack/react-query";
import { RideEvent } from "../lib/models/ride";
import { AuthTokens } from "../lib/models/user";
import { ReactNode, useRef } from "react";
import { STYLES, QUERY_KEYS, isRestErr, toastRestErr } from "../lib/utils";

export const Route = createLazyFileRoute("/dashboard")({
  component: DashBoard,
});

function DashBoard() {
  const { user } = useUserStore();
  const navigate = useNavigate();
  const dialogRef = useRef<HTMLDialogElement>(null);

  if (user.type !== "logged-in") {
    navigate({ to: "/" });
    return <LoadingSpinner content={<span>Redirecting to login...</span>} />;
  }

  return (
    <div className="flex gap-8">
      <div className="flex w-full flex-col gap-2">
        <div className="flex min-h-32 items-center justify-center">
          <button
            className={`text-2xl ${STYLES.button}`}
            onClick={() => {
              dialogRef.current?.showModal();
            }}
          >
            Create a Ride
          </button>
        </div>
        <RideList tokens={user.tokens} />
        <dialog
          className="bg-transparent"
          ref={dialogRef}
        >
          <CreateRideForm afterSubmit={() => dialogRef.current?.close()} />
        </dialog>
      </div>
    </div>
  );
}

function RideList({ tokens }: { tokens: AuthTokens }) {
  const navigate = useNavigate();
  const { setUser } = useUserStore();

  const columns: {
    [K in keyof RideEvent]?: {
      label: string;
      mapField: (field: RideEvent[K]) => string | ReactNode;
    };
  } = {
    driverEmail: {
      label: "Driver",
      mapField: (email) => email,
    },
    locationTo: {
      label: "To",
      mapField: (to) => to,
    },
    locationFrom: {
      label: "From",
      mapField: (from) => from,
    },
    tackingPlaceAt: {
      label: "When",
      mapField: (at) => new Date(at).toLocaleString(),
    },
    status: {
      label: "status",
      mapField: (status) => status,
    },
    schedule: {
      label: "Recurring",
      mapField: (schedule) => {
        if (schedule === null) {
          return "---";
        }

        if (schedule.unit === "weekdays") {
          if (schedule.weekdays === null) {
            return "---";
          }

          if (schedule.interval === 1) {
            // e.g. monday, friday
            return schedule.weekdays.join(", ");
          }

          // e.g. 4. monday, tuesday
          return `${schedule.interval}. ${schedule.weekdays.join(", ")}`;
        }

        if (schedule.interval === 1) {
          // e.g. day
          return schedule.unit.replace(/s$/, "");
        }

        // e.g. 3 weeks
        return `${schedule.interval} ${schedule.unit}`;
      },
    },
  };

  const {
    isPending,
    error,
    data: rides,
  } = useQuery({
    queryKey: [QUERY_KEYS.rideItems],
    queryFn: async () => {
      const res = await fetch(`${import.meta.env.VITE_API_URI}/rides/many`, {
        method: "GET",
        headers: {
          Authorization: tokens.accessToken,
          Accept: "application/json",
        },
      });

      if (res.status === 401) {
        setUser({ type: "logged-out" });
      }

      const data = await res.json();
      if (isRestErr(data)) {
        toastRestErr(data);
        throw new Error("Failed to load rides.");
      }

      return data as RideEvent[];
    },
  });

  if (isPending) {
    return <LoadingSpinner content={<span>Getting rides...</span>} />;
  }

  if (error) {
    return <span className="text-red-500">Failed to load rides</span>;
  }

  if (rides.length === 0) {
    return <span>No rides found</span>;
  }

  return (
    <table className="relative h-fit w-full max-w-full table-auto overflow-x-auto text-lg">
      <thead className="uppercase">
        <RideListRow
          isHeading={true}
          values={Object.values(columns).map(({ label }) => {
            return label;
          })}
        />
      </thead>
      <tbody>
        {rides.map((ride, idx) => {
          return (
            <RideListRow
              key={idx}
              isLast={idx === rides.length - 1}
              values={Object.entries(columns).map(([key, { mapField }]) => {
                return (mapField as any)(ride[key as keyof RideEvent]);
              })}
              onClick={() => {
                navigate({
                  to: "/rides/$rideId",
                  params: { rideId: ride.rideEventId },
                });
              }}
            />
          );
        })}
      </tbody>
    </table>
  );
}

// horrible API, very sorry
function RideListRow({
  values,
  isHeading,
  isLast,
  onClick,
}: {
  values: (string | ReactNode)[];
  isHeading?: boolean;
  isLast?: boolean;
  onClick?: () => void;
}) {
  return (
    <tr
      className={`border-neutral-300 dark:border-neutral-600 ${!isHeading ? "cursor-pointer" : ""} ${!isLast && !isHeading ? "border-b" : ""} ${isHeading ? "sticky top-0 bg-neutral-200 dark:bg-neutral-700" : "bg-neutral-100 dark:bg-neutral-800"}`}
      onClick={onClick}
    >
      {values.map((value, idx) => {
        if (isHeading) {
          return (
            <th
              key={idx}
              className="px-4 py-2 text-left"
            >
              {value}
            </th>
          );
        }

        return (
          <td
            key={idx}
            className="px-4 py-2"
          >
            {value}
          </td>
        );
      })}
    </tr>
  );
}
