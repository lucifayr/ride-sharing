import { useRef } from "react";
import { CreateRideForm } from "../lib/components/CreateRideForm";
import {
  createFileRoute,
  Link,
  ReactNode,
  useNavigate,
} from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import openLinkIcon from "../assets/open-link.svg";
import editIcon from "../assets/edit.svg";
import { STYLES, QUERY_KEYS, isRestErr, toastRestErr } from "../lib/utils";
import { useUserStore } from "../lib/stores";
import { LoadingSpinner } from "../lib/components/Spinner";
import { AuthTokens } from "../lib/models/user";
import { RideEvent, RideSchedule } from "../lib/models/ride";
import { Group } from "../lib/models/models";
import { useForm } from "@tanstack/react-form";

export const Route = createFileRoute("/dashboard")({
  component: Dashboard,
});

function Dashboard() {
  const dialogRefRide = useRef<HTMLDialogElement>(null);
  const dialogRefGroup = useRef<HTMLDialogElement>(null);
  const { user } = useUserStore();

  const navigate = useNavigate();
  if (user.type !== "logged-in") {
    navigate({ to: "/" });
    return <LoadingSpinner content={<span> Redirecting to login...</span>} />;
  }

  return (
    <div className="flex h-full gap-8">
      <div className="flex h-full w-full flex-col gap-2">
        <div className="flex min-h-32 items-center justify-center gap-2">
          <button
            className={`text-2xl ${STYLES.button}`}
            onClick={() => {
              dialogRefRide.current?.showModal();
            }}
          >
            Create a Ride
          </button>
          <button
            className={`text-2xl ${STYLES.button}`}
            onClick={() => {
              dialogRefGroup.current?.showModal();
            }}
          >
            Create a Group
          </button>
        </div>
        <div className="flex h-full gap-8">
          <div className="flex-grow p-8">
            <RideList tokens={user.tokens} />
          </div>
          <div className="h-full min-w-80 border-l-2 border-solid border-neutral-300 p-8 dark:border-neutral-600">
            <GroupList tokens={user.tokens} />
          </div>
        </div>
        <dialog
          className="bg-transparent"
          ref={dialogRefRide}
        >
          <CreateRideForm afterSubmit={() => dialogRefRide.current?.close()} />
        </dialog>
        <dialog
          className="bg-transparent"
          ref={dialogRefGroup}
        >
          <CreateGroupForm
            tokens={user.tokens}
            afterSubmit={() => dialogRefGroup.current?.close()}
          />
        </dialog>
      </div>
    </div>
  );
}

function GroupList({ tokens }: { tokens: AuthTokens }) {
  const { setUser } = useUserStore();

  const {
    isPending,
    error,
    data: groups,
  } = useQuery({
    queryKey: [QUERY_KEYS.groupItems],
    queryFn: async () => {
      const res = await fetch(`${import.meta.env.VITE_API_URI}/groups/many`, {
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
        throw new Error("Failed to load groups.");
      }

      return data as Group[];
    },
  });

  if (isPending) {
    return <LoadingSpinner content={<span>Getting groups...</span>} />;
  }

  if (error) {
    return <span className="text-red-500">Failed to load groups</span>;
  }

  if (groups.length === 0) {
    return <span>No groups found</span>;
  }

  return (
    <div className="flex flex-col gap-8">
      {groups.map((group, idx) => {
        return (
          <Link
            key={idx}
            className="flex flex-col gap-2 text-wrap"
            to="/groups/$groupId"
            params={{ groupId: group.groupId }}
          >
            <span className="text-4xl">{group.name}</span>
            {group.description ? (
              <em className="text-lg">{group.description}</em>
            ) : null}
          </Link>
        );
      })}
    </div>
  );
}

function RideList({ tokens }: { tokens: AuthTokens }) {
  const { setUser } = useUserStore();

  const columns: {
    [K in keyof RideEvent]?: {
      label: string;
      mapField: (field: RideEvent[K], ride: RideEvent) => string | ReactNode;
    };
  } = {
    locationTo: {
      label: "To",
      mapField: (to) => (
        <a
          href={`https://www.google.com/maps/search/Austria+${encodeURIComponent(to.replace(" ", "+"))}`}
          target="_blank"
        >
          {to ?? "---"}
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
          {from ?? "---"}
        </a>
      ),
    },
    tackingPlaceAt: {
      label: "When",
      mapField: (at) => new Date(at).toLocaleString(),
    },
    driverEmail: {
      label: "Driver",
      mapField: (email) => email ?? "---",
    },
    status: {
      label: "Status",
      mapField: (status) => status ?? "---",
    },
    schedule: {
      label: "Recurs",
      mapField: displaySchedule,
    },
    participants: {
      label: "Participants",
      mapField: (participants, ride) => {
        return participants
          ? `${participants.length}/${ride.transportLimit ?? "---"}`
          : "---";
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
              isDoneOrCanceled={ride.status !== "upcoming"}
              values={Object.entries(columns).map(([key, { mapField }]) => {
                return (mapField as any)(ride[key as keyof RideEvent], ride);
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
  isDoneOrCanceled,
  event,
}: {
  values: (string | ReactNode)[];
  isHeading?: boolean;
  isLast?: boolean;
  isDoneOrCanceled?: boolean;
  event?: RideEvent;
}) {
  const { user } = useUserStore();
  const canEdit = user.type === "logged-in" && user.id === event?.createdBy;

  return (
    <tr
      className={`border-neutral-300 dark:border-neutral-600 ${isDoneOrCanceled ? "bg-neutral-300 line-through dark:bg-neutral-500" : ""} ${!isLast && !isHeading ? "border-b" : ""} ${isHeading ? "sticky top-0 bg-neutral-200 dark:bg-neutral-700" : "bg-neutral-100 dark:bg-neutral-800"}`}
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

function CreateGroupForm({
  tokens,
  afterSubmit,
}: {
  tokens: AuthTokens;
  afterSubmit: () => void;
}) {
  const { setUser } = useUserStore();
  const queryClient = useQueryClient();

  const createGroup = useMutation({
    mutationKey: ["create-group"],
    onError: (err) => {
      console.error(err);
    },
    mutationFn: async ({
      name,
      description,
    }: {
      name: string;
      description?: string;
    }) => {
      const res = await fetch(`${import.meta.env.VITE_API_URI}/groups`, {
        method: "POST",
        headers: {
          Authorization: tokens.accessToken,
          Accept: "application/json",
        },
        body: JSON.stringify({
          name,
          description,
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

      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.groupItems] });
    },
  });

  const form = useForm({
    defaultValues: {
      name: "",
      description: undefined as string | undefined,
    },
    onSubmit: async ({ value }) => {
      await createGroup.mutateAsync({
        name: value.name,
        description: value.description,
      });
      afterSubmit?.();
    },
  });

  return (
    <div className="flex h-full flex-col items-center gap-8 bg-neutral-200 dark:bg-neutral-800 dark:text-white">
      <div className="neutral-cyan-700 rounded-lg border-2 p-10">
        <h2 className="doto-h2">Create a Group</h2>
        <form
          className="flex flex-col gap-3"
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
        >
          <form.Field
            name="name"
            children={(field) => {
              return (
                <div>
                  <label className="font-bold">Group Name:</label>
                  <input
                    placeholder="Name..."
                    className="w-full appearance-none rounded border-2 border-gray-200 bg-gray-200 px-4 py-2 leading-tight text-gray-700 focus:border-purple-500 focus:bg-white focus:outline-none"
                    id={field.name}
                    name={field.name}
                    value={field.state.value}
                    onBlur={field.handleBlur}
                    onChange={(e) => {
                      field.handleChange(e.target.value);
                    }}
                  />
                </div>
              );
            }}
          ></form.Field>

          <form.Field
            name="description"
            children={(field) => {
              return (
                <div>
                  <label className="font-bold">Group Description:</label>
                  <input
                    placeholder="Description..."
                    className="w-full appearance-none rounded border-2 border-gray-200 bg-gray-200 px-4 py-2 leading-tight text-gray-700 focus:border-purple-500 focus:bg-white focus:outline-none"
                    id={field.name}
                    name={field.name}
                    value={field.state.value}
                    onBlur={field.handleBlur}
                    onChange={(e) => {
                      field.handleChange(e.target.value);
                    }}
                  />
                </div>
              );
            }}
          ></form.Field>

          <form.Subscribe
            selector={(state) => [state.canSubmit, state.isSubmitting]}
            children={([canSubmit, isSubmitting]) => (
              <button
                type="submit"
                className={`mt-2 flex items-center justify-center ${STYLES.button}`}
                disabled={!canSubmit}
              >
                {isSubmitting ? (
                  <LoadingSpinner content={"Creating..."} />
                ) : (
                  "Submit"
                )}
              </button>
            )}
          />
        </form>
      </div>
    </div>
  );
}
