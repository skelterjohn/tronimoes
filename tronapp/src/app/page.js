import Board from "./components/board/Board";

export default function Home() {
  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex flex-col gap-8 row-start-2 items-center sm:items-start w-full  max-w-screen-md">
        Tronimes
        <div className="w-full aspect-square">
          <Board width={10} height={11} />
        </div>
      </main>
      <footer className="row-start-3 flex gap-6 flex-wrap items-center justify-center">
        JT
      </footer>
    </div>
  );
}
