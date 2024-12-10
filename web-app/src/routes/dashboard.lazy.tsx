import { createLazyFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useUserStore } from "../lib/stores";
import { LoadingSpinner } from "../lib/components/Spinner";
import { CreateRideForm } from "../lib/components/CreateRideForm";
import { useQuery } from "@tanstack/react-query";
import { RideEvent, RideSchedule } from "../lib/models/ride";
import { AuthTokens } from "../lib/models/user";
import { ReactNode, useRef } from "react";
import { STYLES, QUERY_KEYS, isRestErr, toastRestErr } from "../lib/utils";
import openLinkIcon from "../assets/open-link.svg";
import editIcon from "../assets/edit.svg";

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
  const { setUser } = useUserStore();

  const columns: {
    [K in keyof RideEvent]?: {
      label: string;
      mapField: (field: RideEvent[K]) => string | ReactNode;
    };
  } = {
    locationTo: {
      label: "To",
      mapField: (to) => (
        <a
          href={`https://www.google.com/maps/search/Austria+${encodeURIComponent(to.replace(" ", "+"))}`}
          target="_blank"
        >
          {to}
        </a>
      ),
    },
    locationFrom: {
      label: "From",
      mapField: (from) => (
        <a
          href={`https://www.google.com/maps/search/Austria+${encodeURIComponent(from.replace(" ", "+"))}`}
          target="_blank"
        >
          {from}
        </a>
      ),
    },
    tackingPlaceAt: {
      label: "When",
      mapField: (at) => new Date(at).toLocaleString(),
    },
    driverEmail: {
      label: "Driver",
      mapField: (email) => email,
    },
    status: {
      label: "status",
      mapField: (status) => status,
    },
    schedule: {
      label: "Recurs",
      mapField: displaySchedule,
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
              event={ride}
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
  event,
}: {
  values: (string | ReactNode)[];
  isHeading?: boolean;
  isLast?: boolean;
  event?: RideEvent;
}) {
  const { user } = useUserStore();
  const canEdit = user.type === "logged-in" && user.id === event?.createdBy;

  return (
    <tr
      className={`border-neutral-300 dark:border-neutral-600 ${!isLast && !isHeading ? "border-b" : ""} ${isHeading ? "sticky top-0 bg-neutral-200 dark:bg-neutral-700" : "bg-neutral-100 dark:bg-neutral-800"}`}
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
      <td className="px-4 py-2">
        {isHeading ? null : (
          <Link
            to="/rides/$rideId"
            params={{ rideId: event!.rideEventId }}
          >
            <img
              className="h-6 w-6 dark:invert"
              src={canEdit ? editIcon : openLinkIcon}
              alt={canEdit ? "edit" : "open"}
            />
          </Link>
        )}
      </td>
    </tr>
  );
}

export function displaySchedule(schedule: RideSchedule | null): string {
  if (schedule === null) {
    return "---";
  }

  if (schedule.unit === "weekdays") {
    if (schedule.weekdays === null) {
      return "---";
    }

    if (schedule.interval === 1) {
      // e.g. every monday, friday
      return `every ${schedule.weekdays.join("/")}`;
    }

    // e.g. every 4. monday/tuesday
    return `every ${schedule.interval}. ${schedule.weekdays.join("/")}`;
  }

  if (schedule.interval === 1) {
    // e.g. every day
    return `every ${schedule.unit.replace(/s$/, "")}`;
  }

  // e.g. every 3 weeks
  return `every ${schedule.interval} ${schedule.unit}`;
}
