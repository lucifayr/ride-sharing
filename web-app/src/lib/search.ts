export const SEARCH_FILTERS = {
  source: "from",
  destination: "to",
  datetimeBefore: "before",
  datetimeAfter: "after",
  driver: "driver",
  owner: "owner",
  participants: "participants",
} as const;

type SearchFilters = {
  source?: string;
  destination?: string;
  datetimeBefore?: Date;
  datetimeAfter?: Date;
  driver?: string;
  owner?: string;
  participants?: string[];
};

export function parseSearchString(rawSearchString: string): SearchFilters {
  throw new Error("unimplemented");
}

export function recommendSearchFilters(partialSearchString: string): string[] {
  throw new Error("unimplemented");
}
