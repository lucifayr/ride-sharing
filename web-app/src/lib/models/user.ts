export type User = UserLoggedIn | UserLoggedOut;

export type UserLoggedIn = {
  type: "logged-in";
  id: string;
  name: string;
  email: string;
};
export type UserLoggedOut = { type: "logged-out" };
