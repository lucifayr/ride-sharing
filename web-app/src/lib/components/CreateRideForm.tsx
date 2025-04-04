import { useForm } from "@tanstack/react-form";
import { LoadingSpinner } from "./Spinner";
import { useQueryClient } from "@tanstack/react-query";
import { useUserStore } from "../stores";
import { useNavigate } from "@tanstack/react-router";
import { STYLES } from "../utils";
import { RideSchedule } from "../models/ride";
import { createRide } from "../createRide";

export function CreateRideForm({ afterSubmit }: { afterSubmit?: () => void }) {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { user, setUser } = useUserStore();

  if (user.type !== "logged-in") {
    navigate({ to: "/" });
    return <LoadingSpinner content={<span>Redirecting to login...</span>} />;
  }

  const defaultValues = {
    locationFrom: "",
    locationTo: "",
    tackingPlaceAt: "",
    transportLimit: "4",
    recurs: "",
  };


  const form = useForm({
    defaultValues,
    onSubmit: async ({ value }) => {
      await createRide(defaultValues, user, setUser, queryClient).mutateAsync(value);
      afterSubmit?.();
    },
  });

  return (
    <div className="flex h-full flex-col items-center gap-8 bg-neutral-200 dark:bg-neutral-800 dark:text-white">
      <div className="neutral-cyan-700 rounded-lg border-2 p-10">
        <h2 className="doto-h2">Create a Ride</h2>
        <form
          className="flex flex-col gap-3"
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
        >
          <form.Field
            name="locationFrom"
            children={(field) => {
              return (
                <div className="flex flex-col">
                  <label
                    className="font-bold"
                    htmlFor={field.name}
                  >
                    From
                  </label>
                  <input
                    className="w-full appearance-none rounded border-2 border-gray-200 bg-gray-200 px-4 py-2 leading-tight text-gray-700 focus:border-purple-500 focus:bg-white focus:outline-none"
                    placeholder="Kaindorf"
                    required={true}
                    id={field.name}
                    name={field.name}
                    value={field.state.value}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                  />
                </div>
              );
            }}
          />

          <form.Field
            name="locationTo"
            children={(field) => {
              return (
                <div className="flex flex-col">
                  <label
                    className="font-bold"
                    htmlFor={field.name}
                  >
                    To
                  </label>
                  <input
                    className="w-full appearance-none rounded border-2 border-gray-200 bg-gray-200 px-4 py-2 leading-tight text-gray-700 focus:border-purple-500 focus:bg-white focus:outline-none"
                    placeholder="Murek"
                    required={true}
                    id={field.name}
                    name={field.name}
                    value={field.state.value}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                  />
                </div>
              );
            }}
          />

          <form.Field
            name="tackingPlaceAt"
            children={(field) => {
              return (
                <div className="flex flex-col">
                  <label
                    className="font-bold"
                    htmlFor={field.name}
                  >
                    When
                  </label>
                  <input
                    id={field.name}
                    name={field.name}
                    required={true}
                    type="datetime-local"
                    className="rounded border-2 border-gray-200 bg-gray-200 px-4 py-2 leading-tight text-gray-700 focus:border-cyan-700 focus:bg-white focus:outline-none"
                    value={field.state.value}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                  />
                </div>
              );
            }}
          />

          <form.Field
            name="transportLimit"
            children={(field) => {
              return (
                <div className="flex flex-col">
                  <label
                    className="font-bold"
                    htmlFor={field.name}
                  >
                    Max. Passengers
                  </label>
                  <input
                    className="w-full appearance-none rounded border-2 border-gray-200 bg-gray-200 px-4 py-2 leading-tight text-gray-700 focus:border-purple-500 focus:bg-white focus:outline-none"
                    placeholder="4"
                    id={field.name}
                    name={field.name}
                    value={field.state.value}
                    onBlur={field.handleBlur}
                    required={true}
                    type="number"
                    onChange={(e) => field.handleChange(e.target.value)}
                  />
                </div>
              );
            }}
          />

          <form.Field
            name="recurs"
            validators={{
              onChange: ({ value }) => {
                if (value.length === 0) {
                  return undefined;
                }

                if (parseRecuring(value) === undefined) {
                  return "Invalid interval";
                }
              },
            }}
            children={(field) => {
              return (
                <div className="flex flex-col">
                  <label
                    className="font-bold"
                    htmlFor={field.name}
                  >
                    Recurs?
                  </label>
                  <input
                    className="w-full appearance-none rounded border-2 border-gray-200 bg-gray-200 px-4 py-2 leading-tight text-gray-700 focus:border-purple-500 focus:bg-white focus:outline-none"
                    placeholder="2 days"
                    required={false}
                    id={field.name}
                    name={field.name}
                    value={field.state.value}
                    onBlur={field.handleBlur}
                    onChange={(e) => field.handleChange(e.target.value)}
                  />
                  {field.state.meta.errors.length ? (
                    <em className="text-red-500">
                      {field.state.meta.errors.join(",")}
                    </em>
                  ) : (
                    <>&nbsp;</>
                  )}
                </div>
              );
            }}
          />

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

const validBasicUnits = ["day", "week", "month", "year"];
const validWeekdays = [
  "monday",
  "tuesday",
  "wednesday",
  "thursday",
  "friday",
  "saturday",
  "sunday",
];

export function parseRecuring(text: string): RideSchedule | undefined {
  let t = text.trim().toLowerCase();
  if (t.startsWith("every")) {
    t = t.substring("every".length).trim();
  }

  const parts = t.split(" ", 2);
  let partUnit: string;
  let interval = 1;

  if (parts.length === 0) {
    return undefined;
  }

  if (parts.length === 1) {
    partUnit = parts[0];
  } else {
    partUnit = parts[1];

    interval = parseInt(parts[0]);
    const intervalInt = Math.floor(interval);
    if (
      isNaN(interval) ||
      interval <= 0 ||
      intervalInt.toString() != parts[0]
    ) {
      return undefined;
    }
  }

  const unit = validBasicUnits.find(
    (u) => u === partUnit || u + "s" === partUnit,
  );
  if (unit) {
    return {
      unit: unit + "s",
      interval,
      weekdays: null,
    };
  }

  const weekdays: string[] = [];
  const days = partUnit.split(",");
  for (const day of days) {
    const d = day.trim();
    const duplicate = weekdays.includes(d);
    if (!validWeekdays.includes(d) || duplicate) {
      return undefined;
    }
    weekdays.push(d);
  }

  return {
    unit: "weekdays",
    interval,
    weekdays,
  };
}
