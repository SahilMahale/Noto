import Navbar from '@/components/NavBar'
import type { AppContext } from '@/context/authContext'
import {  createRootRouteWithContext, Link, Outlet } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/router-devtools'

const isProd = process.env.NODE_ENV === "PROD"
interface RouterAppContext {
    auth: AppContext
}
export const Route = createRootRouteWithContext<RouterAppContext>()({
  component: () => (
    <>
      <Navbar />
      {!isProd && (
        <TanStackRouterDevtools position="bottom-right" initialIsOpen={false} />
      )}
    </>
  ),
  notFoundComponent: () => { return (<p className="text-red-500">Root component not Sound</p>) }
})
