import { createFileRoute } from "@tanstack/react-router";
import { useEffect } from "react";
import { AuthTokens } from "../lib/models/user";

export const Route = createFileRoute("/authenticate")({
  validateSearch: validateAuthTokens,
  component: ReceiveAuthRedirect,
});

function ReceiveAuthRedirect() {
  const { accessToken, refreshToken } = Route.useSearch();
  useEffect(() => {
    if (window.opener) {
      const data = { accessToken, refreshToken };
      window.opener.postMessage(data);
      window.close();
    }
  }, []);

  return <span>Redirecting...</span>;
}

export function validateAuthTokens(data: Record<string, unknown>): AuthTokens {
  const accessToken = data["accessToken"];
  if (typeof accessToken !== "string") {
    throw new Error("Invalid/Missing search parameter 'accessToken'.");
  }

  const refreshToken = data["refreshToken"];
  if (typeof refreshToken !== "string") {
    throw new Error("Invalid/Missing search parameter 'refreshToken'.");
  }
  return {
    accessToken,
    refreshToken,
  };
}
