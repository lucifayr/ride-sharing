import { toast } from "react-toastify";

export const QUERY_KEYS = {
  rideItems: "ride-items",
  groupItems: "group-items",
  rideSingle: "ride-single",
  groupSingle: "group-single",
} as const;

export const STYLES = {
  button:
    "rounded min-w-[92px] bg-cyan-700 py-2 px-4 duration-150 hover:bg-cyan-600 font-bold",
  buttonDanger:
    "rounded min-w-[92px] bg-red-700 py-2 px-4 duration-150 hover:bg-red-800 font-bold",
} as const;

export type RestError = {
  errors: {
    title: string;
    details?: string;
  }[];
};

export function isRestErr(data: any): data is RestError {
  return "errors" in data && Array.isArray(data.errors);
}

export function toastRestErr(err: RestError) {
  for (const { title, details } of err.errors) {
    toast(title, { type: "error" });
    if (details !== undefined) {
      console.error(`['${title}']: ${details}`);
    }
  }
}
