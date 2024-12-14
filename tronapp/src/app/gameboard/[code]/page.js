import Game from "@/app/components/game/Game";
import { TipProvider } from "@/app/components/tutorial/InnerTip";
export default async function Gameboard({params}) {
  let p = await params;
  const code = p.code;
  return (
    <div className="h-screen flex justify-center">
      <main className="flex-1 flex flex-col gap-2 items-center justify-center w-full max-h-screen">
	  <TipProvider>
        <Game code={code} />
	  </TipProvider>
      </main>
    </div>
  );
}
