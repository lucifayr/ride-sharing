import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useAuthStore, useUserStore } from "../lib/stores";
import { useMutation } from "@tanstack/react-query";
import { validateAuthTokens } from "./authenticate";
import googleIcon from "../assets/google-icon.svg";
import { ReactNode } from "react";
import { LoadingSpinner } from "../lib/components/Spinner";

export const Route = createFileRoute("/")({
  component: LoginPage,
});

type LoginPageState = "already-logged-in" | "received-auth-tokens" | "idle";

function LoginPage() {
  const { user, setUser } = useUserStore();
  const { tokens } = useAuthStore();
  const navigate = useNavigate();

  const loginAsUser = useMutation({
    mutationKey: ["login-as-user"],
    mutationFn: async () => {
      if (tokens === undefined) {
        return;
      }

      const res = await fetch(`${import.meta.env.VITE_API_URI}/users`, {
        method: "GET",
        headers: {
          Authorization: tokens.accessToken,
          Accept: "application/json",
        },
      });

      if (!res.ok) {
        throw new Error(
          `Login as user failed with message ${await res.text()}`,
        );
      }

      const user = await res.json();
      // TODO: validate
      setUser({
        type: "logged-in",
        id: user.id,
        email: user.email,
        name: user.name,
      });
    },
  });

  const state: LoginPageState =
    user.type === "logged-in"
      ? "already-logged-in"
      : tokens !== undefined
        ? "received-auth-tokens"
        : "idle";

  switch (state) {
    case "already-logged-in": {
      navigate({ to: "/dashboard" });
      return (
        <Shell>
          <LoadingSpinner content={<span>Redirecting to dashboard...</span>} />
        </Shell>
      );
    }
    case "received-auth-tokens": {
      if (loginAsUser.status === "idle") {
        loginAsUser.mutate();
      }

      return (
        <Shell>
          <LoadingSpinner content={<span>Logging you in...</span>} />
        </Shell>
      );
    }
    case "idle": {
      return (
        <Shell>
          <button
            className={`light:hover:shadow flex gap-2 rounded-lg border border-slate-200 px-4 py-2 transition duration-150 hover:border-slate-400 dark:border-slate-700 dark:hover:border-slate-500`}
            onClick={() => authenticateGoogle()}
          >
            <img
              className="h-6 w-6"
              src={googleIcon}
              alt="google logo"
            />
            <span>Login with Google</span>
          </button>
        </Shell>
      );
    }
  }
}

function Shell({ children }: { children: ReactNode }) {
  return (
    <div className="flex h-full flex-col items-center justify-center">
      <div className="light:shadow-md light:shadow-gray-200 flex max-h-[360px] max-w-[480px] flex-col items-center gap-4 rounded p-8 dark:bg-neutral-800">
        <h1 className="text-4xl font-bold">Ride Sharing</h1>
        {children}
      </div>
    </div>
  );
}

async function authenticateGoogle() {
  openSignInWindow(`${import.meta.env.VITE_API_URI}/auth/google/login`);
}

let windowObjectReference: Window | null = null;
let previousUrl: string | undefined = undefined;

function openSignInWindow(url: string) {
  window.removeEventListener("message", receiveAuthMessage);

  if (windowObjectReference === null || windowObjectReference.closed) {
    windowObjectReference = window.open(url);
  } else if (previousUrl !== url) {
    windowObjectReference = window.open(url);
    windowObjectReference?.focus();
  } else {
    windowObjectReference.focus();
  }

  window.addEventListener(
    "message",
    (event) => receiveAuthMessage(event),
    false,
  );
  previousUrl = url;
}

function receiveAuthMessage(event: MessageEvent) {
  if (event.origin !== window.location.origin) {
    return;
  }

  try {
    const tokens = validateAuthTokens(event.data);
    useAuthStore.getState().setTokens(tokens);
  } catch (err) {
    console.error("Invalid authentication token data.", err);
  }
}
