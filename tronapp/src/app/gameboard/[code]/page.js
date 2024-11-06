import Game from "../../components/game/Game";

export default function Gameboard({params}) {
  return (
    <div className="h-screen flex justify-center">
      <main className="flex-1 flex flex-col gap-2 items-center justify-center w-full max-h-screen">
        <Game code={params.code} />
      </main>
    </div>
  );
}
