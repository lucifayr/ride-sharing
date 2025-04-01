import { useMemo, useRef, useState } from "react";

type SearchInputEntry<T> = { value: T; display: string; ordinal: string };

// I am sorry
export function SearchInput<T>({
  items,
  entryMap,
  onConfirm,
  extraProps,
}: {
  items: T[];
  entryMap: (item: T) => SearchInputEntry<T>;
  onConfirm: (item: T | undefined) => void;
  extraProps?: {
    inputField?: React.InputHTMLAttributes<HTMLInputElement>;
  };
}) {
  const containerRef = useRef<HTMLDivElement>(null);

  const [search, setSearch] = useState("");
  const [open, setOpen] = useState(false);
  const [confirmed, setConfirmed] = useState(false);

  const possibleItems = useMemo(() => {
    return items.map(entryMap).filter((entry) => {
      return entry.ordinal.toLowerCase().includes(search.toLowerCase());
    });
  }, [items, search]);

  return (
    <div
      className="relative"
      ref={containerRef}
    >
      <input
        data-testid={"input-test-id"}
        className="w-full border-b border-neutral-200 bg-transparent p-1 text-xl focus:border-cyan-500 focus:outline-none disabled:border-none dark:border-neutral-500"
        type="text"
        autoComplete="off"
        value={search}
        onFocus={() => {
          setOpen(true);
          setConfirmed(false);
        }}
        onBlur={(e) => {
          setOpen(containerRef.current?.contains(e.relatedTarget) ?? false);
        }}
        onChange={(e) => {
          setOpen(true);
          setConfirmed(false);
          setSearch(e.target.value);
        }}
        onKeyDown={(e) => {
          if (e.key !== "Enter" || possibleItems.length !== 1) {
            return;
          }
          setConfirmed(true);
          setSearch(possibleItems[0].display);
          onConfirm(possibleItems[0].value);
        }}
        {...extraProps?.inputField}
      />
      <div className="absolute flex w-full flex-col divide-y-[1px] divide-neutral-300 dark:divide-neutral-700">
        {open && !confirmed
          ? possibleItems
            .sort((a, b) => {
              return a.ordinal.localeCompare(b.ordinal);
            })
            .map((entry, idx) => {
              return (
                <button
                  key={idx}
                  className="bg-neutral-200 p-2 text-left text-lg font-semibold duration-150 hover:bg-neutral-300 focus:bg-neutral-300 focus:outline-none dark:bg-neutral-800 hover:dark:bg-neutral-700 focus:dark:bg-neutral-700"
                  onClick={() => {
                    setConfirmed(true);
                    setSearch(entry.display);
                    onConfirm(entry.value);
                  }}
                >
                  <span>{entry.display}</span>
                </button>
              );
            })
          : null}
      </div>
    </div>
  );
}
