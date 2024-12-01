import { RideGroup, ScheduledRideGroup } from "./models/models";
import { User } from "./models/user";
import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";

type UserStore = {
  user: User;
  setUser: (user: User) => void;
};

type ScheduledRideGroupsStore = {
  groups: ScheduledRideGroup[];
  setGroups: (ride: ScheduledRideGroup[]) => void;
};

type RideGroupsStore = {
  groups: RideGroup[];
  setGroups: (ride: RideGroup[]) => void;
};

export const useUserStore = create<UserStore>()(
  persist(
    (set) => ({
      user: {
        type: "logged-out",
      },
      setUser: (user) => set({ user }),
    }),
    {
      name: "user",
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
