import Game from "@/app/components/game/Game";
import { TipProvider } from "@/app/components/tutorial/Tip";
export default function Gameboard({params}) {
  return (
    <div className="h-screen flex justify-center">
      <main className="flex-1 flex flex-col gap-2 items-center justify-center w-full max-h-screen">
	  <TipProvider>
        <Game code={params.code} />
	  </TipProvider>
      </main>
    </div>
  );
}
