"use client";

import type { ReactNode } from "react";
import { googleLoginURL } from "@/lib/api";
import { useAuth } from "@/lib/auth";

// Gates page content behind a logged-in session. Wrap any page/section that
// needs real auth (every data page does); the nav/header itself stays
// visible regardless (see app/layout.tsx), so a logged-out visitor always
// has a way to log in.
export function AuthGate({ children }: { children: ReactNode }) {
  const { user, loading } = useAuth();

  if (loading) {
    return <p className="text-sm text-neutral-500">Carregando...</p>;
  }

  if (!user) {
    return (
      <div className="flex flex-col items-center gap-4 rounded-lg border border-neutral-200 bg-white px-6 py-16 text-center">
        <h1 className="text-2xl font-semibold">Stock-Stonks</h1>
        <p className="max-w-sm text-neutral-600">
          Entre com sua conta Google para cadastrar ativos, registrar operações e
          acompanhar sua posição.
        </p>
        <a
          href={googleLoginURL}
          className="rounded bg-neutral-900 px-4 py-2 text-sm text-white hover:bg-neutral-800"
        >
          Entrar com Google
        </a>
      </div>
    );
  }

  return <>{children}</>;
}
