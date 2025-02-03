export const SEARCH_FILTERS = {
  source: "from",
  destination: "to",
  dateBefore: "before",
  dateAfter: "after",
  driver: "driver",
  owner: "owner",
  participants: "participants",
} as const satisfies {
  [K in keyof SearchFilters]: any;
};

type SearchFilters = {
  source?: string;
  destination?: string;
  dateBefore?: Date;
  dateAfter?: Date;
  driver?: string;
  owner?: string;
  participants?: string[];
};

/**
 * Parses a ride search string by finding any search filters contained in it.
 * Search filters are of the format `:{filter-name} {filter-value}`, where
 * `{filter-name}` is a value such as `from`, `to`, or `driver` and
 * `{filter-value}` is the matching value for that filter.
 *
 * @param rawSearchString The search string to parse.
 */
export function parseSearchString(rawSearchString: string): SearchFilters {
  let filters: SearchFilters = {};

  let filterValue: string | undefined;
  let filterKey: keyof SearchFilters | undefined;

  for (let i = 0; i < rawSearchString.length; i++) {
    let filterMatched = false;

    const char = rawSearchString.charAt(i);
    const isFilterStart = char === ":";
    const nextSpaceIdx = rawSearchString.indexOf(" ", i + 1);

    // Check if we have found a filter of the form `:{filter-name} {value}` and
    // try parsing its name.
    if (isFilterStart && nextSpaceIdx !== -1) {
      const foundName = rawSearchString.substring(i + 1, nextSpaceIdx);
      const matchingFilter = Object.entries(SEARCH_FILTERS).find(
        ([_, value]) => {
          return value === foundName;
        },
      );

      if (matchingFilter !== undefined) {
        // If there is already an active filter, parse its value and add it to
        // `filters`.
        if (filterKey !== undefined && filterValue !== undefined) {
          const value = parseFilterValue(filterKey, filterValue);
          if (value !== undefined) {
            filters[filterKey] = value as any;
          }
        }

        filterKey = matchingFilter[0] as keyof SearchFilters;
        filterValue = "";
        i = nextSpaceIdx; // skip over the filter name that was just scanned
        filterMatched = true;
      }
    }

    if (!filterMatched && filterValue !== undefined) {
      filterValue += char;
    }
  }

  // Add any remaining filter to `filters`.
  if (filterKey !== undefined && filterValue !== undefined) {
    const value = parseFilterValue(filterKey, filterValue);
    if (value !== undefined) {
      filters[filterKey] = value as any;
    }
  }

  return filters;
}

function parseFilterValue<K extends keyof SearchFilters>(
  filter: K,
  value: string,
): SearchFilters[K] {
  switch (filter) {
    case "source":
    case "destination":
    case "owner":
    case "driver": {
      return value.trim() as SearchFilters[K];
    }

    case "participants": {
      return value.split(",").map((str) => str.trim()) as SearchFilters[K];
    }

    case "dateAfter":
    case "dateBefore": {
      const ms = Date.parse(value.trim());
      if (Number.isNaN(ms)) {
        return undefined;
      }

      return new Date(ms) as SearchFilters[K];
    }
  }
}

export function recommendSearchFilters(partialSearchString: string): string[] {
  const lastFilterStartIdx = partialSearchString.lastIndexOf(":");
  if (lastFilterStartIdx === -1) {
    return [];
  }

  const filterStart = partialSearchString.substring(lastFilterStartIdx + 1);

  return Object.values(SEARCH_FILTERS).filter((value) =>
    value.startsWith(filterStart),
  );
}
