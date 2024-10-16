import { create } from "zustand";
import { IUser } from "../models/models";
import { UserStore } from "./storeTypes";

//Users that have not logged in yet will get this dummy user with
//basically no ability to do anything but login. Identifyer will always be -1
export const useUserStore = create<UserStore>((set) => ({
  user: {
    userId: -1,
    username: "",
    email: "",
    age: 0,
    driver: false,
    residence: "",
  },
  setUser: (user: IUser) => set({ user: user }),
}));
