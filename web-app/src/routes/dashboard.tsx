import { useEffect, useMemo, useRef, useState } from "react";
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
import { AuthTokens, UserLoggedIn } from "../lib/models/user";
import { RideEvent, RideSchedule } from "../lib/models/ride";
import { Group, GroupMessage } from "../lib/models/models";
import { useForm } from "@tanstack/react-form";
import { parseSearchString, SearchFilters } from "../lib/search";

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
        <div className="flex h-full gap-8">
          <div className="flex flex-grow flex-col gap-8 p-8">
            <button
              className={`text-2xl ${STYLES.button} max-w-80`}
              onClick={() => {
                dialogRefRide.current?.showModal();
              }}
            >
              Create a Ride
            </button>
            <RideList user={user} />
          </div>
          <div className="flex max-h-full min-w-[600px] flex-1 flex-col gap-8 overflow-y-auto border-l-2 border-solid border-neutral-300 p-4 dark:border-neutral-600">
            <button
              className={`text-2xl ${STYLES.button} w-96 max-w-96`}
              onClick={() => {
                dialogRefGroup.current?.showModal();
              }}
            >
              Create a Group
            </button>
            <GroupCol user={user} />
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

function GroupCol({ user }: { user: UserLoggedIn }) {
  const { setUser } = useUserStore();
  const [activeGroup, setActiveGroup] = useState<Group | undefined>(undefined);

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
          Authorization: user.tokens.accessToken,
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
    <div className="flex max-h-full flex-1 flex-grow flex-col gap-16 overflow-y-auto">
      <SearchInput
        items={groups}
        entryMap={(group) => ({
          value: group,
          ordinal: `${group.name}${group.description ?? ""}`,
          display: group.name,
        })}
        onConfirm={(group) => {
          setActiveGroup(group);
        }}
        extraProps={{
          inputField: { placeholder: "Choose a group..." },
        }}
      />
      <GroupChat
        group={activeGroup}
        user={user}
      />
    </div>
  );
}

const msgs = [
  {
    messageId: "1",
    content: "first message",
    sentBy: "NnCaPHQLC9",
    sentByEmail: "test@example.com",
    createdAt: "todo",
  },
  {
    messageId: "2",
    content: "second message",
    sentBy: "NnCaPHQLC9",
    sentByEmail: "test@example.com",
    createdAt: "todo",
  },
  {
    messageId: "3",
    content: "reply to first message",
    sentBy: "nmBSHcxyvn",
    sentByEmail: "2@other.com",
    createdAt: "todo",
    repliesTo: "2",
  },
  {
    messageId: "3",
    content: "reply to first message",
    sentBy: "nmBSHcxyvn",
    sentByEmail: "2@other.com",
    createdAt: "todo",
    repliesTo: "2",
  },
  {
    messageId: "3",
    content: "reply to first message",
    sentBy: "nmBSHcxyvn",
    sentByEmail: "2@other.com",
    createdAt: "todo",
    repliesTo: "2",
  },
  {
    messageId: "2",
    content: "second message",
    sentBy: "NnCaPHQLC9",
    sentByEmail: "test@example.com",
    createdAt: "todo",
  },
  {
    messageId: "2",
    content: "second message",
    sentBy: "NnCaPHQLC9",
    sentByEmail: "test@example.com",
    createdAt: "todo",
  },
  {
    messageId: "2",
    content: "second message",
    sentBy: "NnCaPHQLC9",
    sentByEmail: "test@example.com",
    createdAt: "todo",
  },
  {
    messageId: "3",
    content: "reply to first message",
    sentBy: "nmBSHcxyvn",
    sentByEmail: "2@other.com",
    createdAt: "todo",
    repliesTo: "2",
  },
  {
    messageId: "3",
    content: "reply to first message",
    sentBy: "nmBSHcxyvn",
    sentByEmail: "2@other.com",
    createdAt: "todo",
    repliesTo: "2",
  },
  {
    messageId: "3",
    content: "reply to first message",
    sentBy: "nmBSHcxyvn",
    sentByEmail: "2@other.com",
    createdAt: "todo",
    repliesTo: "2",
  },
  {
    messageId: "3",
    content: "reply to first message",
    sentBy: "nmBSHcxyvn",
    sentByEmail: "2@other.com",
    createdAt: "todo",
    repliesTo: "2",
  },
  {
    messageId: "3",
    content: "reply to first message",
    sentBy: "nmBSHcxyvn",
    sentByEmail: "2@other.com",
    createdAt: "todo",
    repliesTo: "2",
  },
  {
    messageId: "3",
    content: "reply to first message",
    sentBy: "nmBSHcxyvn",
    sentByEmail: "2@other.com",
    createdAt: "todo",
    repliesTo: "2",
  },
  {
    messageId: "3",
    content: "reply to first message",
    sentBy: "nmBSHcxyvn",
    sentByEmail: "2@other.com",
    createdAt: "todo",
    repliesTo: "2",
  },
  {
    messageId: "3",
    content: "reply to first message",
    sentBy: "nmBSHcxyvn",
    sentByEmail: "2@other.com",
    createdAt: "todo",
    repliesTo: "2",
  },
  {
    messageId: "3",
    content: "reply to first message",
    sentBy: "nmBSHcxyvn",
    sentByEmail: "2@other.com",
    createdAt: "todo",
    repliesTo: "2",
  },
] satisfies GroupMessage[];

function GroupChat({ group, user }: { group?: Group; user: UserLoggedIn }) {
  const msgContainer = useRef<HTMLDivElement>(null);

  useEffect(() => {
    msgContainer.current?.scrollTo({
      behavior: "instant",
      top: msgContainer.current.scrollHeight,
    });
  }, []);

  if (group === undefined) {
    return null;
  }

  return (
    <div className="flex max-h-full flex-1 flex-grow flex-col gap-4 overflow-y-auto">
      <span className="truncate text-wrap text-2xl">
        Messages from{" "}
        <Link
          to="/groups/$groupId"
          className="italic underline"
          params={{ groupId: group.groupId }}
        >
          {group.name}
        </Link>
      </span>
      <div
        ref={msgContainer}
        className="flex max-h-full flex-1 flex-grow flex-col gap-2 overflow-y-auto p-4"
      >
        {msgs.map((msg, idx) => {
          return (
            <div
              key={idx}
              className={`relative max-w-[80%] rounded-md p-2 ${
                msg.sentBy === user.id
                  ? "self-end border-b-4 border-l-4 border-cyan-900 bg-cyan-800"
                  : "self-start border-b-4 border-r-4 border-neutral-400 bg-neutral-300 dark:border-neutral-800 dark:bg-neutral-700"
              }`}
            >
              <div
                className={`absolute top-[-12px] flex h-8 w-8 items-center justify-center rounded-full bg-neutral-400 dark:bg-neutral-800 ${
                  msg.sentBy === user.id ? "left-[-24px]" : "right-[-24px]"
                }`}
              >
                <span className="text-lg font-semibold">
                  {msg.sentByEmail
                    .toUpperCase()
                    .substring(0, Math.min(2, msg.sentByEmail.indexOf("@")))}
                </span>
              </div>
              <div>
                {msg.repliesTo !== undefined && (
                  <Reply
                    msgs={msgs}
                    reply={msg}
                    user={user}
                  />
                )}
                <span className="text-lg">{msg.content}</span>
              </div>
            </div>
          );
        })}
      </div>
      <input
        className="w-full border-b border-neutral-200 bg-transparent p-1 text-xl focus:border-cyan-500 focus:outline-none disabled:border-none dark:border-neutral-500"
        type="text"
        autoComplete="off"
        placeholder="Send a message..."
      />
    </div>
  );
}

function Reply({
  msgs,
  reply,
  user,
}: {
  msgs: GroupMessage[];
  reply: GroupMessage;
  user: UserLoggedIn;
}) {
  const originalMsg = msgs.find((msg) => msg.messageId === reply.repliesTo);
  if (originalMsg === undefined) {
    return;
  }

  return (
    <div
      className={`rounded border-neutral-100 p-2 dark:border-neutral-400 ${
        reply.sentBy === user.id
          ? "border-r-2 bg-cyan-900"
          : "border-l-2 bg-neutral-400 dark:bg-neutral-800"
      }`}
    >
      <span className="text-lg">{originalMsg.content}</span>
    </div>
  );
}

type SearchInputEntry<T> = { value: T; display: string; ordinal: string };

// I am sorry
function SearchInput<T>({
  items,
  entryMap,
  onConfirm,
  extraProps,
}: {
  items: T[];
  entryMap: (item: T) => SearchInputEntry<T>;
  onConfirm: (item: T | undefined) => void;
  extraProps?: {
    inputField?: React.InputHTMLAttributes<HTMLInputElement>;
  };
}) {
  const containerRef = useRef<HTMLDivElement>(null);

  const [search, setSearch] = useState("");
  const [open, setOpen] = useState(false);
  const [confirmed, setConfirmed] = useState(false);

  const possibleItems = useMemo(() => {
    return items.map(entryMap).filter((entry) => {
      return entry.ordinal.toLowerCase().includes(search.toLowerCase());
    });
  }, [items, search]);

  return (
    <div
      className="relative"
      ref={containerRef}
    >
      <input
        className="w-full border-b border-neutral-200 bg-transparent p-1 text-xl focus:border-cyan-500 focus:outline-none disabled:border-none dark:border-neutral-500"
        type="text"
        autoComplete="off"
        value={search}
        onFocus={() => {
          setOpen(true);
          setConfirmed(false);
        }}
        onBlur={(e) => {
          setOpen(containerRef.current?.contains(e.relatedTarget) ?? false);
        }}
        onChange={(e) => {
          setOpen(true);
          setConfirmed(false);
          setSearch(e.target.value);
        }}
        onKeyDown={(e) => {
          if (e.key !== "Enter" || possibleItems.length !== 1) {
            return;
          }

          setConfirmed(true);
          setSearch(possibleItems[0].display);
          onConfirm(possibleItems[0].value);
        }}
        {...extraProps?.inputField}
      />
      <div className="absolute flex w-full flex-col divide-y-[1px] divide-neutral-300 dark:divide-neutral-700">
        {open && !confirmed
          ? possibleItems
              .sort((a, b) => {
                return a.ordinal.localeCompare(b.ordinal);
              })
              .map((entry, idx) => {
                return (
                  <button
                    key={idx}
                    className="bg-neutral-200 p-2 text-left text-lg font-semibold duration-150 hover:bg-neutral-300 focus:bg-neutral-300 focus:outline-none dark:bg-neutral-800 hover:dark:bg-neutral-700 focus:dark:bg-neutral-700"
                    onClick={() => {
                      setConfirmed(true);
                      setSearch(entry.display);
                      onConfirm(entry.value);
                    }}
                  >
                    <span>{entry.display}</span>
                  </button>
                );
              })
          : null}
      </div>
    </div>
  );
}

function RideList({ user }: { user: UserLoggedIn }) {
  const tokens = user.tokens;
  const inputRefRideSearch = useRef<HTMLInputElement>(null);
  const { setUser } = useUserStore();

  const [filters, setFilters] = useState<SearchFilters>({});

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
    <div className="flex flex-col gap-8">
      <div>
        <span className="text-2xl font-semibold">Search for Rides</span>
        <input
          className="w-full border-b border-neutral-200 bg-transparent p-1 text-xl focus:border-cyan-500 focus:outline-none disabled:border-none dark:border-neutral-500"
          placeholder=":from My cool city"
          type="text"
          autoComplete="off"
          ref={inputRefRideSearch}
          onKeyDown={(e) => {
            if (e.key !== "Enter" || !inputRefRideSearch.current) {
              return;
            }

            const newFilters = parseSearchString(
              inputRefRideSearch.current.value,
            );
            setFilters(newFilters);
          }}
        />
      </div>
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
          {rides
            .filter((ride) => rideMatchesFilters(ride, user, filters))
            .map((ride, idx) => {
              return (
                <RideListRow
                  key={idx}
                  isLast={idx === rides.length - 1}
                  isDoneOrCanceled={ride.status !== "upcoming"}
                  values={Object.entries(columns).map(([key, { mapField }]) => {
                    return (mapField as any)(
                      ride[key as keyof RideEvent],
                      ride,
                    );
                  })}
                  event={ride}
                />
              );
            })}
        </tbody>
      </table>
    </div>
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

function rideMatchesFilters(
  ride: RideEvent,
  me: UserLoggedIn,
  filters: SearchFilters,
): boolean {
  if (filters.source !== undefined && ride.locationFrom !== filters.source) {
    return false;
  }

  if (
    filters.destination !== undefined &&
    ride.locationTo !== filters.destination
  ) {
    return false;
  }

  if (filters.status !== undefined && ride.status !== filters.status) {
    return false;
  }

  if (
    filters.owner !== undefined &&
    aliasedEmail(me, ride.createdByEmail) !== aliasedEmail(me, filters.owner)
  ) {
    return false;
  }

  if (
    filters.driver !== undefined &&
    aliasedEmail(me, ride.driverEmail) !== aliasedEmail(me, filters.driver)
  ) {
    return false;
  }

  const rideDate = new Date(Date.parse(ride.tackingPlaceAt));
  if (filters.dateAfter !== undefined && rideDate < filters.dateAfter) {
    return false;
  }

  if (filters.dateBefore !== undefined && rideDate > filters.dateBefore) {
    return false;
  }

  if (
    filters.participants !== undefined &&
    !filters.participants.every((participant) => {
      const hasParticipant = ride.participants.some(
        (rideParticipant) =>
          aliasedEmail(me, rideParticipant.email) ===
          aliasedEmail(me, participant),
      );
      return hasParticipant;
    })
  ) {
    return false;
  }

  return true;
}

function aliasedEmail(me: UserLoggedIn, email: string): string {
  if (email === "me") {
    return me.email;
  }

  return email;
}
