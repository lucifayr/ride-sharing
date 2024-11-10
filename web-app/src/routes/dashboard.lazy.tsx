import { createLazyFileRoute, useNavigate } from "@tanstack/react-router";
import { useUserStore } from "../lib/stores";
import { LoadingSpinner } from "../lib/components/Spinner";

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

  return <div>Hello {user.name}!</div>;
}
