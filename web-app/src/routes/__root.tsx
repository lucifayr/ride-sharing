import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";
import { ToastContainer } from "react-toastify";

export const Route = createRootRoute({
  component: () => {
    const theme: "dark" | "light" = window.matchMedia(
      "(prefers-color-scheme: dark)",
    ).matches
      ? "dark"
      : "light";

    return (
      <>
        <Outlet />
        <TanStackRouterDevtools />
        <ToastContainer
          toastStyle={{ background: theme === "light" ? "#f5f5f5" : "#262626" }}
        />
      </>
    );
  },
});
