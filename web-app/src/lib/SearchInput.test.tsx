// @vitest-environment jsdom
import { describe, expect, it } from "vitest";
import { SearchInput } from "./components/SearchInput";
import { cleanup, fireEvent, render } from "@testing-library/react";

type MockPerson = {
  firstname: string,
  lastname: string,
  birthdate: Date,
  inGroups: number
}

const mockPeople: MockPerson[] = [
  {
    firstname: "Tom",
    lastname: "Peeping",
    birthdate: new Date(Date.parse("10.10.2005")),
    inGroups: 10
  }
]

describe("Test SearchInput confirm functionality", () => {
  it("should call confirm", async () => {
    let resItem: any = undefined;
    const comp = render(<SearchInput
      items={mockPeople}
      onConfirm={(item) => {
        resItem = item
      }}
      entryMap={(person) => ({
        value: person,
        ordinal: person.firstname,
        display: person.firstname,
      })} />)

    const inputItem = await comp.findByTestId("input-test-id");
    inputItem.focus();
    fireEvent.change(inputItem, { target: { value: "Tom" } })
    fireEvent.keyDown(inputItem, { key: "Enter" })
    expect(resItem.firstname).toEqual("Tom")
    cleanup();
  })
  it("should not call confirm", async () => {
    let resItem: any = undefined;
    const comp = render(<SearchInput
      items={mockPeople}
      onConfirm={(item) => {
        resItem = item
      }}
      entryMap={(person) => ({
        value: person,
        ordinal: person.firstname,
        display: person.firstname,
      })} />)
    const inputItem = await comp.findByTestId("input-test-id");
    inputItem.focus();
    fireEvent.change(inputItem, { target: { value: "Tim" } })
    fireEvent.keyDown(inputItem, { key: "Enter" })
    expect(resItem).toBeUndefined()
    cleanup();
  })
})


