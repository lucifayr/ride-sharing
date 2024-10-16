export const SideBard = () => {
  const mockArray = ["Gamer", "Shit", "Ass", "Fuck"];
  return (
    <div className="h-screen w-52 border border-neutral-500">
      <div>
        {mockArray.map((d, i) => {
          return (
            <div
              className="m-2 h-10 rounded bg-gray-300 p-2 hover:bg-gray-500"
              key={i}
            >
              {d}
            </div>
          );
        })}
      </div>
    </div>
  );
};
