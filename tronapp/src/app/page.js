import Board from "./components/board/Board";
import Game from "./components/game/Game";

export default function Home() {
  return (
    <div className="items-center justify-items-center min-h-screen h-screen">
      <main className="flex flex-col gap-8 row-start-2 items-center sm:items-start w-full  max-w-screen-md">
        <Game/>
      </main>
      <footer className="row-start-3 flex gap-6 flex-wrap items-center justify-center">
        JT
      </footer>
    </div>
  );
}
