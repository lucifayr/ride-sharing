import { useForm } from "@tanstack/react-form";
import { LoadingSpinner } from "./Spinner";

export function CreateRideForm() {
  const form = useForm({
    defaultValues: {
      locationFrom: "",
      locationTo: "",
      tackingPlaceAt: "",
    },
    onSubmit: async ({ value }) => {
      console.log(value);
    },
  });

  return (
    <div className="flex h-full flex-col items-center gap-8">
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

          <form.Subscribe
            selector={(state) => [state.canSubmit, state.isSubmitting]}
            children={([canSubmit, isSubmitting]) => (
              <button
                type="submit"
                className="mt-2 rounded bg-cyan-700 p-2"
                disabled={!canSubmit}
              >
                {isSubmitting ? <LoadingSpinner /> : "Submit"}
              </button>
            )}
          />
        </form>
      </div>
    </div>
  );
}
