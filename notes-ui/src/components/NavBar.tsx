
import type { ReactNode } from 'react';
import { Link, Outlet, useNavigate, useRouter } from '@tanstack/react-router';
import { useAuth } from '../context/authContext';
import type { ClassValue } from 'clsx';
import { classMerge } from './utility';

const AnchorLinks = ({ className, to, children, isTitle, ...props }: { className?: ClassValue, to: string, children?: ReactNode, isTitle: boolean }) => {
    return (
        <Link
            to={to}
            className={classMerge(
                // Base styles
                'font-sans py-2.5 rounded-lg text-zinc-200 text-center font-bold',
                'hover:outline-sky-300 hover:outline-2',
                // Conditional styles based on isTitle
                isTitle
                    ? 'tracking-tighter text-4xl px-3'
                    : 'text-base px-5',
                // Custom className passed as prop
                className
            )}
            activeProps={{ className: 'bg-zinc-800' }}
            {...props}
        >
            {children}
        </Link>
    );
};

const LogoutButton = ({ logOutHandler }: { logOutHandler: () => Promise<boolean> }) => {
    const router = useRouter()
    const navigate = useNavigate()
    const handleLogout = () => {
        if (window.confirm('Are you sure you want to logout'))
            logOutHandler().then(() => {
                router.invalidate().finally(() => {
                    navigate({ to: "/" })
                })
            })
    }
    return (
        <button
            className="font-sans px-3 py-2.5 rounded-lg text-base text-amber-200 text-center font-bold 
      hover:text-amber-500
      hover:outline-amber-500 hover:outline-4"
            onClick={handleLogout}
        >
            Logout
        </button>
    );
};

const Navbar = ({ children }: { children?: ReactNode }) => {
    const { Context: appContext, LogOut } = useAuth();
    if (!LogOut) {
        console.error("Error: LogOut function not found")
        return
    }
    let isAdmin = false
    const homeLink = appContext.isLoggedIn ? '/home' : '/'
    if (appContext.isLoggedIn) {
        isAdmin = appContext.claims.type === 'admin'
    }
    return (
        <div className=" bg-zinc-950 mx-auto py-2 border-2 border-zinc-500">
            <nav className=" bg-zinc-950 rounded-lg border-2 border-zinc-950 text-gray-200 container mx-auto flex flex-wrap items-center justify-between">
                <AnchorLinks to={homeLink} isTitle={true}>
                    Nōto
                </AnchorLinks>
                <div className="px-4">
                    {appContext.isLoggedIn ? (
                        <>
                            <AnchorLinks to="/book/tables" isTitle={false}>
                                Book Tables
                            </AnchorLinks>
                            {isAdmin && (
                                <AnchorLinks to="/users" isTitle={false}>
                                    Users
                                </AnchorLinks>
                            )}
                            <LogoutButton logOutHandler={LogOut} />
                        </>
                    ) : (
                        <>
                            <AnchorLinks className="px-6" to="/login" isTitle={false}>
                                Login
                            </AnchorLinks>
                            <AnchorLinks className="px-4" to="/signup" isTitle={false}>
                                SignUp
                            </AnchorLinks>
                        </>
                    )}
                </div>
            </nav>
            <Outlet />
            {children}
        </div>
    );
};
export default Navbar;
