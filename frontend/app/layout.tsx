import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import Link from "next/link";
import { AuthGate } from "@/components/AuthGate";
import { UserMenu } from "@/components/UserMenu";
import { AuthProvider } from "@/lib/auth";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Stock-Stonks",
  description: "Registro de operações de ações e acompanhamento para Imposto de Renda",
};

const NAV_LINKS = [
  { href: "/", label: "Início" },
  { href: "/assets", label: "Ativos" },
  { href: "/operations", label: "Operações" },
  { href: "/positions", label: "Posições" },
];

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html
      lang="pt-BR"
      className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}
    >
      <body className="min-h-full flex flex-col bg-neutral-50 text-neutral-900">
        <AuthProvider>
          <header className="border-b border-neutral-200 bg-white">
            <nav className="mx-auto flex max-w-5xl items-center gap-6 px-4 py-3">
              <span className="font-semibold">Stock-Stonks</span>
              <ul className="flex flex-1 gap-4 text-sm text-neutral-600">
                {NAV_LINKS.map((link) => (
                  <li key={link.href}>
                    <Link href={link.href} className="hover:text-neutral-900">
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
              <UserMenu />
            </nav>
          </header>
          <main className="mx-auto w-full max-w-5xl flex-1 px-4 py-8">
            <AuthGate>{children}</AuthGate>
          </main>
        </AuthProvider>
      </body>
    </html>
  );
}
