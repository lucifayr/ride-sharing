import { useState } from "react";

export const GroupBar = () => {
  const [barMode, setBarMode] = useState("Groups");
  return (
    <div className="flex h-full flex-col">
      <div className="m-2 flex flex-grow flex-col overflow-hidden rounded-xl border-2 border-gray-200 p-2 text-center">
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
              <h2
                className={`doto-h2 transition-colors duration-1000 ${barMode === "Groups" ? "text-white" : "text-gray-500"}`}
              >
                {"Groups"}
              </h2>
              <h2
                className={`doto-h2 transition-colors duration-1000 ${barMode === "Groups" ? "text-gray-500" : "text-white"}`}
              >
                {"Rides"}
              </h2>
            </div>
            <div
              className={`flex transition-all ${barMode === "Groups" ? "justify-start" : "justify-end"}`}
            >
              <div
                className={`mx-3 h-3 rounded transition-all duration-1000 ease-in-out ${barMode === "Groups" ? "w-1/2 bg-cyan-500" : "w-[43%] bg-orange-500"} `}
              ></div>
            </div>
          </div>
        </div>
        <div className="flex flex-grow flex-col overflow-scroll">
          <div className="flex flex-col gap-y-2 p-2">
            {groups.map((group, i) => {
              return (
                <GroupItem
                  groupName={group.groupName}
                  barMode={barMode}
                  key={i}
                />
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
};

const GroupItem = ({
  groupName,
  barMode,
}: {
  groupName: string;
  barMode: string;
}) => {
  return (
    <button
      className={`h-12 rounded-lg text-xl font-bold transition-all duration-1000 ${barMode === "Groups" ? "bg-cyan-500 hover:bg-cyan-800" : "bg-orange-500 hover:bg-orange-800"}`}
    >
      {groupName}
    </button>
  );
};
