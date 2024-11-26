import { RideGroup, ScheduledRideGroup } from "./models/models";
import { User } from "./models/user";
import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";

export type AuthTokens = {
  accessToken: string;
  refreshToken: string;
};

type UserStore = {
  user: User;
  setUser: (user: User) => void;
};

type AuthStore = {
  tokens: AuthTokens | undefined;
  setTokens: (tokens: AuthTokens) => void;
  clearTokens: () => void;
};

type ScheduledRideGroupsStore = {
  groups: ScheduledRideGroup[];
  setGroups: (ride: ScheduledRideGroup[]) => void;
};

type RideGroupsStore = {
  groups: RideGroup[];
  setGroups: (ride: RideGroup[]) => void;
};

const fakeLoggedInUser: User = {
  id: "(/986)",
  type: "logged-in",
  name: "TestUser",
  email: "test@test.com",
};

export const useUserStore = create<UserStore>((set) => ({
  user: {
    type: "logged-out",
  },
  setUser: (user) => set({ user }),
}));

export const useAuthStore = create<AuthStore>()(
  persist(
    (set) => ({
      tokens: undefined,
      setTokens: (tokens) => set({ tokens }),
      clearTokens: () => set({ tokens: undefined }),
    }),
    {
      name: "auth-store",
      storage: createJSONStorage(() => localStorage),
    },
  ),
);

export const useScheduledRideGroupsStore = create<ScheduledRideGroupsStore>(
  (set) => ({
    groups: [],
    setGroups: (groups) => set({ groups: groups }),
  }),
);

export const useRideGroupStore = create<RideGroupsStore>((set) => ({
  groups: [],
  setGroups: (groups) => set({ groups: groups }),
}));
