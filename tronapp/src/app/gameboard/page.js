import Game from "../components/game/Game";

export default function Gameboard() {
  return (
    <div className="items-center justify-items-center min-h-screen h-screen">
      <main className="flex flex-col gap-8 row-start-2 items-center sm:items-start w-full  max-w-screen-md">
        <Game/>
      </main>
    </div>
  );
}
