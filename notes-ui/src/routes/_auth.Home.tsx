import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export const Route = createFileRoute('/_auth/Home')({
  component: HomeComponent,
})

function HomeComponent() {
  return (
    <div className="container mx-auto p-8 text-center relative z-10">
      <Card className="bg-zinc-900 text-zinc-300 cursor-pointer rounded-lg 
    hover:outline-2 hover:outline-blue-500 transition-colors">
        <CardHeader className="gap-4">
          <CardTitle className="text-3xl font-bold">Bun + React</CardTitle>
          <CardDescription className="text-zinc-400">
            Edit <code className="rounded bg-amber-100 px-[0.3rem] py-[0.2rem] font-mono">src/routes/index.tsx</code> and save to test HMR
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="flex items-center justify-center gap-8">
            <div className="text-center">
              <div className="text-5xl font-bold text-amber-400">⚡</div>
              <p className="mt-2 text-sm text-zinc-400">Lightning Fast</p>
            </div>
            <div className="text-center">
              <div className="text-5xl font-bold text-blue-400">⚛️</div>
              <p className="mt-2 text-sm text-zinc-400">Powered by React</p>
            </div>
            <div className="text-center">
              <div className="text-5xl font-bold text-purple-400">🍞</div>
              <p className="mt-2 text-sm text-zinc-400">Built with Bun</p>
            </div>
          </div>

          <div className="border-t border-zinc-800 pt-6">
            <h3 className="text-lg font-semibold mb-3">Quick Start</h3>
            <div className="space-y-2 text-left">
              <div className="bg-zinc-800 rounded p-3 font-mono text-sm">
                <span className="text-zinc-500">$</span> <span className="text-amber-400">bun</span> install
              </div>
              <div className="bg-zinc-800 rounded p-3 font-mono text-sm">
                <span className="text-zinc-500">$</span> <span className="text-amber-400">bun</span> dev
              </div>
            </div>
          </div>

          <p className="text-zinc-500 text-sm">
            Experience the blazing speed of Bun combined with the power of React. 
            Your development workflow just got a serious upgrade.
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
