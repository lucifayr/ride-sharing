import { useState } from "react";
import { useScheduledRideGroupsStore } from "../stores";

export const GroupBar = () => {
  //const { groups } = useScheduledRideGroupsStore();
  const groups = [
    {
      groupName: "Group1",
    },
  ];
  const [barMode, setBarMode] = useState("Groups");
  return (
    <div className="flex h-full flex-col">
      <div className="m-2 flex w-1/4 flex-grow flex-col overflow-hidden rounded-xl border-2 border-gray-200 p-2 text-center">
        <div
          onClick={() => {
            if (barMode === "Groups") {
              setBarMode("Rides");
            } else {
              setBarMode("Groups");
            }
          }}
        >
          <div className="flex flex-col">
            <div className="flex flex-row justify-center gap-4">
              <h2 className="doto-h2">{"Groups"}</h2>
              <h2 className="doto-h2">{"Rides"}</h2>
            </div>
            <div
              className={`flex w-full flex-row justify-${barMode === "Groups" ? "start" : "end"}`}
            >
              <div className="groupBarSelect"></div>
            </div>
          </div>
        </div>
        <div className="flex flex-grow flex-col overflow-scroll">
          <div className="flex flex-col gap-y-2 p-2">
            {groups.map((group, i) => {
              if (barMode === "Groups") {
                return (
                  <GroupItem
                    groupName={group.groupName}
                    key={i}
                  />
                );
              } else {
                return (
                  <RideItem
                    groupName={group.groupName}
                    key={i}
                  />
                );
              }
            })}
          </div>
        </div>
      </div>
    </div>
  );
};

const GroupItem = ({ groupName }: { groupName: string }) => {
  return (
    <button className="h-12 rounded-lg bg-cyan-500 text-xl font-bold hover:bg-cyan-800">
      {groupName}
    </button>
  );
};

const RideItem = ({ groupName }: { groupName: string }) => {
  return (
    <button className="h-12 rounded-lg bg-orange-500 text-xl font-bold hover:bg-orange-800">
      {groupName}
    </button>
  );
};
