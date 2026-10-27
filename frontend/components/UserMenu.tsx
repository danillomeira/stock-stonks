"use client";

import { googleLoginURL } from "@/lib/api";
import { useAuth } from "@/lib/auth";

// Nav-bar auth control: shows a login link when logged out, or the user's
// name/avatar + logout button when logged in. Reusable wherever the app
// shell needs to show auth state (currently just the root layout's header).
export function UserMenu() {
  const { user, loading, logout } = useAuth();

  if (loading) {
    return null;
  }

  if (!user) {
    return (
      <a href={googleLoginURL} className="text-sm font-medium text-neutral-700 hover:text-neutral-900">
        Entrar com Google
      </a>
    );
  }

  return (
    <div className="flex items-center gap-3 text-sm">
      {user.AvatarURL && (
        // eslint-disable-next-line @next/next/no-img-element -- external, unoptimized Google avatar URL
        <img src={user.AvatarURL} alt="" className="h-6 w-6 rounded-full" referrerPolicy="no-referrer" />
      )}
      <span className="text-neutral-700">{user.Name}</span>
      <button onClick={() => logout()} className="text-neutral-500 hover:text-neutral-900 hover:underline">
        Sair
      </button>
    </div>
  );
}
