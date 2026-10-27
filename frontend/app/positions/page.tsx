"use client";

import { useEffect, useState } from "react";
import { api, type Asset, type Position } from "@/lib/api";

export default function PositionsPage() {
  const [positions, setPositions] = useState<Position[]>([]);
  const [assets, setAssets] = useState<Asset[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    Promise.all([api.listPositions(), api.listAssets()])
      .then(([pos, assetList]) => {
        setPositions(pos);
        setAssets(assetList);
      })
      .catch((err) => setError(err instanceof Error ? err.message : "Falha ao carregar posições"))
      .finally(() => setLoading(false));
  }, []);

  function assetLabel(id: number) {
    const asset = assets.find((a) => a.ID === id);
    return asset ? `${asset.Ticker} — ${asset.Name}` : `#${id}`;
  }

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-2xl font-semibold">Posições</h1>
      <p className="text-sm text-neutral-600">
        Quantidade e preço médio calculados a partir do histórico completo de
        operações (custo médio ponderado), conforme regra da Receita Federal.
      </p>

      {error && <p className="text-sm text-red-600">{error}</p>}

      {loading ? (
        <p className="text-sm text-neutral-500">Carregando...</p>
      ) : positions.length === 0 ? (
        <p className="text-sm text-neutral-500">Nenhuma posição ainda — registre operações.</p>
      ) : (
        <table className="w-full text-left text-sm">
          <thead>
            <tr className="border-b border-neutral-200 text-neutral-500">
              <th className="py-2">Ativo</th>
              <th className="py-2">Quantidade</th>
              <th className="py-2">Preço médio</th>
              <th className="py-2">Lucro realizado</th>
            </tr>
          </thead>
          <tbody>
            {positions.map((pos) => (
              <tr key={pos.AssetID} className="border-b border-neutral-100">
                <td className="py-2 font-medium">{assetLabel(pos.AssetID)}</td>
                <td className="py-2">{pos.Quantity}</td>
                <td className="py-2">{pos.AverageCost}</td>
                <td className={`py-2 ${Number(pos.RealizedProfit) < 0 ? "text-red-600" : "text-green-700"}`}>
                  {pos.RealizedProfit}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
