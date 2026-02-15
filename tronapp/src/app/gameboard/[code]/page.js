import Game from "@/app/components/game/Game";

export default async function Gameboard({ params }) {
  const p = await params;
  const code = p.code;
  return (
    <div className="h-screen flex justify-center">
      <main className="flex-1 flex flex-col gap-2 items-center justify-center w-full max-h-screen">
        <Game code={code} />
      </main>
    </div>
  );
}
