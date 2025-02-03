import { expect, it, describe } from "vitest";
import {
  parseSearchString,
  recommendSearchFilters,
  SEARCH_FILTERS,
} from "./search";

describe("parseSearchString", () => {
  it("should not find any search filters", () => {
    expect(parseSearchString("")).toEqual({});
    expect(parseSearchString("some text")).toEqual({});
    expect(parseSearchString("fjlfj1421jkjr3c:fjsladf")).toEqual({});
  });

  it("should find one search filter", () => {
    expect(parseSearchString(":from Graz")).toEqual({
      source: "Graz",
    });

    expect(parseSearchString(":before 12.12.2024")).toEqual({
      dateBefore: new Date(Date.parse("12.12.2024")),
    });
  });

  it("should find mulitple search filters", () => {
    expect(
      parseSearchString(
        " :from Graz :to Kaindorf  :participants user1@example.com,user:2@example.com",
      ),
    ).toEqual({
      source: "Graz",
      destination: "Kaindorf",
      participants: ["user1@example.com", "user:2@example.com"],
    });

    expect(parseSearchString(":driver me@example.com :to Kaindorf")).toEqual({
      driver: "me@example.com",
      destination: "Kaindorf",
    });
  });

  it("should respect spaces in search filters", () => {
    expect(
      parseSearchString(":to Kaindorf an der Sulm :owner me@example.com"),
    ).toEqual({
      owner: "me@example.com",
      destination: "Kaindorf an der Sulm",
    });
  });
});

describe("recommendSearchFilters", () => {
  it("should recommend nothing", () => {
    expect(recommendSearchFilters("")).toEqual([]);
    expect(recommendSearchFilters("some text")).toEqual([]);
    expect(recommendSearchFilters(":from Graz")).toEqual([]);
    expect(recommendSearchFilters(":x")).toEqual([]);
  });

  it("should recommend filters", () => {
    expect(recommendSearchFilters(":")).toEqual(Object.values(SEARCH_FILTERS));
    expect(recommendSearchFilters(":fr")).toEqual(["from"]);
    expect(recommendSearchFilters(":to")).toEqual(["to"]);
    expect(recommendSearchFilters(":dri")).toEqual(["driver"]);
  });
});
