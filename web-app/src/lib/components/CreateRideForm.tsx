import { useForm } from "@tanstack/react-form";
import { LoadingSpinner } from "./Spinner";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useUserStore } from "../stores";
import { useNavigate } from "@tanstack/react-router";
import { QUERY_KEYS, STYLES } from "../utils";

export function CreateRideForm({ afterSubmit }: { afterSubmit?: () => void }) {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { user } = useUserStore();

  if (user.type !== "logged-in") {
    navigate({ to: "/" });
    return <LoadingSpinner content={<span>Redirecting to login...</span>} />;
  }

  const defaultValues = {
    locationFrom: "",
    locationTo: "",
    tackingPlaceAt: "",
    transportLimit: "4",
  };

  const createRide = useMutation({
    mutationKey: ["create-ride-from-submmit"],
    mutationFn: async (params: typeof defaultValues) => {
      // TODO: check status
      await fetch(`${import.meta.env.VITE_API_URI}/rides`, {
        method: "POST",
        body: JSON.stringify({
          ...params,
          tackingPlaceAt: new Date(params.tackingPlaceAt).toISOString(),
          transportLimit: parseInt(params.transportLimit),
          driver: user.id,
        }),
        headers: {
          Authorization: user.tokens.accessToken,
          Accept: "application/json",
        },
      });

      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.rideItems] });
    },
  });

  const form = useForm({
    defaultValues,
    onSubmit: async ({ value }) => {
      await createRide.mutateAsync(value);
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
