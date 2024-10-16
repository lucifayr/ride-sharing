import { IUser } from "../models/models";

export type UserStore = {
  user: IUser;
  setUser: (user: IUser) => void;
};
