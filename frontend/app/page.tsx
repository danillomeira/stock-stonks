import Link from "next/link";

const CARDS = [
  {
    href: "/assets",
    title: "Ativos",
    description: "Cadastre os tickers que você negocia (ações, FIIs, ETFs, BDRs).",
  },
  {
    href: "/operations",
    title: "Operações",
    description: "Registre compras e vendas com data, quantidade, preço e taxas.",
  },
  {
    href: "/positions",
    title: "Posições",
    description: "Veja a posição atual, preço médio e lucro realizado por ativo.",
  },
];

export default function HomePage() {
  return (
    <div className="flex flex-col gap-8">
      <div>
        <h1 className="text-2xl font-semibold">Stock-Stonks</h1>
        <p className="mt-2 text-neutral-600">
          Acompanhamento de operações em ações e exportação de dados para o Imposto de
          Renda. MVP inicial: cadastro de ativos, operações e cálculo de posição por
          preço médio.
        </p>
      </div>
      <div className="grid gap-4 sm:grid-cols-3">
        {CARDS.map((card) => (
          <Link
            key={card.href}
            href={card.href}
            className="rounded-lg border border-neutral-200 bg-white p-4 transition hover:border-neutral-400"
          >
            <h2 className="font-medium">{card.title}</h2>
            <p className="mt-1 text-sm text-neutral-600">{card.description}</p>
          </Link>
        ))}
      </div>
    </div>
  );
}
