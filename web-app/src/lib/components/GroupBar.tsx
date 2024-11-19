const mockNames = [
  "Group1",
  "Group2",
  "Group3",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
  "Group4",
];

export const GroupBar = () => {
  return (
    <div className="flex h-full flex-col">
      <div className="m-2 flex w-1/4 flex-grow flex-col overflow-hidden rounded-xl border-2 border-gray-200 p-2 text-center">
        <h2 className="doto-h2">Groups</h2>
        <div className="flex flex-grow flex-col overflow-scroll">
          <div className="flex flex-col gap-y-2 p-2">
            {mockNames.map((mockName, i) => {
              return (
                <GroupItem
                  groupName={mockName}
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

const GroupItem = ({ groupName }: { groupName: string }) => {
  return (
    <button className="h-12 rounded-lg bg-cyan-500 text-xl font-bold hover:bg-cyan-800">
      {groupName}
    </button>
  );
};
