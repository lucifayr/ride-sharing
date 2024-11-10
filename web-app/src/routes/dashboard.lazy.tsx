import { createLazyFileRoute, useNavigate } from "@tanstack/react-router";
import { useUserStore } from "../lib/stores";
import { LoadingSpinner } from "../lib/components/Spinner";
import { UserLoggedIn } from "../lib/models/user";
import { useForm } from "@tanstack/react-form";

export const Route = createLazyFileRoute("/dashboard")({
  component: DashBoard,
});

function DashBoard() {
  const { user } = useUserStore();
  const navigate = useNavigate();
  if (user.type !== "logged-in") {
    navigate({ to: "/" });
    return <LoadingSpinner content={<span>Redirecting to login...</span>} />;
  }

  return (
    <div className="flex h-full flex-col items-center gap-8">
      <span>
        Hello {user.name}! <br />
        Good design will soon find your heart
      </span>
      <CreateRideForm user={user} />
    </div>
  );
}

// WIP
function CreateRideForm({ user }: { user: UserLoggedIn }) {
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
    <form
      className="flex flex-col gap-4"
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
            <div className="flex flex-col gap-2">
              <label htmlFor={field.name}>From:</label>
              <input
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
            <div className="flex flex-col gap-2">
              <label htmlFor={field.name}>To:</label>
              <input
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
            <div className="flex flex-col gap-2">
              <label htmlFor={field.name}>When:</label>
              <input
                id={field.name}
                name={field.name}
                type="datetime-local"
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
            disabled={!canSubmit}
          >
            {isSubmitting ? <LoadingSpinner /> : "Submit"}
          </button>
        )}
      />
    </form>
  );
}
