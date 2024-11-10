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

export const useUserStore = create<UserStore>((set) => ({
  user: {
    type: "logged-out",
  },
  setUser: (user) => set({ user }),
}));

export const useAuthStore = create<AuthStore>((set) => ({
  tokens: undefined,
  setTokens: (tokens) => set({ tokens }),
}));
