import { QueryClient, useMutation } from "@tanstack/react-query";
import { DefaultValues } from "./models/models";
import { User, UserLoggedIn } from "./models/user";
import { parseRecuring } from "./components/CreateRideForm";
import { isRestErr, QUERY_KEYS, toastRestErr } from "./utils";


export const createRide = (defaultValues: DefaultValues, user: UserLoggedIn, setUser: (user: User) => void, queryClient: QueryClient) => {
  return useMutation({
    mutationKey: ["create-ride-from-submmit"],
    mutationFn: async (params: typeof defaultValues) => {
      const res = await fetch(`${import.meta.env.VITE_API_URI}/rides`, {
        method: "POST",
        body: JSON.stringify({
          ...params,
          tackingPlaceAt: new Date(params.tackingPlaceAt).toISOString(),
          transportLimit: parseInt(params.transportLimit),
          driver: user.id,
          schedule: parseRecuring(params.recurs),
        }),
        headers: {
          Authorization: user.tokens.accessToken,
          Accept: "application/json",
        },
      });

      if (res.status === 401) {
        setUser({ type: "logged-out" });
        return;
      }

      if (res.status !== 201) {
        const data = await res.json();
        if (isRestErr(data)) {
          toastRestErr(data);
          return;
        }
      }

      queryClient.invalidateQueries({ queryKey: [QUERY_KEYS.rideItems] });
    },
  });
}

