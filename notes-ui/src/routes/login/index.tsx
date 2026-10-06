import { Switch } from '@headlessui/react';
import { useState } from 'react';
import { LoginForm } from './-LoginForm';
import { createFileRoute, redirect } from '@tanstack/react-router';
import { z } from 'zod';
export const Route = createFileRoute('/login/')({
  validateSearch: z.object({
    redirect: z.string().optional().catch('')
  }),
  beforeLoad: ({ context, search }) => {
    if (context.auth.Context.isLoggedIn) {
      throw redirect({
        to: search.redirect || '/home'
      })
    }
  },
  component: Login
})

function Login({ isAdmin = false }) {
  const [isAdminToggled, setIsAdminToggled] = useState(isAdmin); //just for pretty styling user and admin login goto the same URL
  const handleAdminToggle = () => {
    setIsAdminToggled(!isAdminToggled);
  };
  return (
    <div className="bg-zinc-950 min-h-screen py-2 px-2">
      <div className="container mx-auto rounded-md bg-zinc-900 flex flex-col">
        <div className="self-end px-2 py-2">
          <Switch.Group>
            <Switch.Label className="mr-4 text-zinc-200 tracking-tight text-base font-bold">
              {isAdminToggled ? 'Admin' : 'User'}
            </Switch.Label>
            <Switch
              checked={isAdminToggled}
              onChange={handleAdminToggle}
              className={`${isAdminToggled ? 'bg-amber-400' : 'bg-zinc-300'
                }  relative inline-flex h-6 w-11 items-center rounded-full`}
            >
              <span
                className={`${isAdminToggled
                  ? 'translate-x-6  bg-amber-700'
                  : 'translate-x-1  bg-zinc-500'
                  } inline-block h-4 w-4 transform rounded-full  transition-transform`}
              />
            </Switch>
          </Switch.Group>
        </div>
        <h2 className="font-sans text-zinc-300 tracking-tight text-3xl text-center font-bold ">
          Log in
        </h2>

        <LoginForm isAdmin={isAdminToggled} />
      </div>
    </div>
  );
}

