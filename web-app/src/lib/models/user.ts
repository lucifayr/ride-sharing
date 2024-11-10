export type User =
  | {
      type: "logged-in";
      id: string;
      name: string;
      email: string;
    }
  | { type: "logged-out" };
