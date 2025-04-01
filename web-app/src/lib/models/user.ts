export type User = UserLoggedIn | UserLoggedOut | UserBeforeLogin;

export type UserLoggedIn = {
  type: "logged-in";
  id: string;
  name: string;
  email: string;
  isAdmin: boolean;
  isBlocked: boolean;
  tokens: AuthTokens;
};

export type UserLoggedOut = { type: "logged-out" };

export type UserBeforeLogin = {
  type: "before-logged-in";
  tokens: AuthTokens;
};

export type AuthTokens = {
  accessToken: string;
  refreshToken: string;
};
