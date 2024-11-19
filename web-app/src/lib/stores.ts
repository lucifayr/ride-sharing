import { RideGroup, ScheduledRideGroup } from "./models/models";
import { User } from "./models/user";
import { create } from "zustand";

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
};

type ScheduledRideGroupsStore = {
  groups: ScheduledRideGroup[];
  setGroups: (ride: ScheduledRideGroup[]) => void;
};

type RideGroupsStore = {
  groups: RideGroup[];
  setGroups: (ride: RideGroup[]) => void;
};

export const useUserStore = create<UserStore>((set) => ({
  user: {
    id: "(/986)",
    type: "logged-in",
    name: "TestUser",
    email: "test@test.com",
  },
  setUser: (user) => set({ user }),
}));

export const useAuthStore = create<AuthStore>((set) => ({
  tokens: undefined,
  setTokens: (tokens) => set({ tokens }),
}));

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
