import Game from "../components/game/Game";

export default function Gameboard() {
  return (
    <div className="min-h-screen max-h-screen flex">
      <main className="flex-1 flex flex-col gap-8 items-center sm:items-start w-full">
        <Game/>
      </main>
    </div>
  );
}
