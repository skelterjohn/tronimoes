import Game from "../components/game/Game";

export default function Gameboard() {
  return (
    <div className="h-screen flex">
      <main className="flex-1 flex flex-col gap-2 items-center sm:items-start w-full max-h-screen">
        <Game/>
      </main>
    </div>
  );
}
