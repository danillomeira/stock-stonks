"use client";

import { useEffect, useState } from "react";
import { api, type Asset, type AssetType } from "@/lib/api";

const ASSET_TYPES: AssetType[] = ["STOCK", "FII", "ETF", "BDR"];

export default function AssetsPage() {
  const [assets, setAssets] = useState<Asset[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [ticker, setTicker] = useState("");
  const [name, setName] = useState("");
  const [type, setType] = useState<AssetType>("STOCK");
  const [submitting, setSubmitting] = useState(false);

  async function load() {
    setLoading(true);
    try {
      setAssets(await api.listAssets());
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Falha ao carregar ativos");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, []);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setSubmitting(true);
    try {
      await api.createAsset({ ticker: ticker.toUpperCase(), name, type });
      setTicker("");
      setName("");
      setType("STOCK");
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Falha ao criar ativo");
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(id: number) {
    try {
      await api.deleteAsset(id);
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Falha ao remover ativo");
    }
  }

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-2xl font-semibold">Ativos</h1>

      <form onSubmit={handleSubmit} className="flex flex-wrap items-end gap-3 rounded-lg border border-neutral-200 bg-white p-4">
        <label className="flex flex-col text-sm">
          Ticker
          <input
            required
            value={ticker}
            onChange={(e) => setTicker(e.target.value)}
            placeholder="PETR4"
            className="mt-1 rounded border border-neutral-300 px-2 py-1"
          />
        </label>
        <label className="flex flex-col text-sm">
          Nome
          <input
            required
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Petrobras PN"
            className="mt-1 rounded border border-neutral-300 px-2 py-1"
          />
        </label>
        <label className="flex flex-col text-sm">
          Tipo
          <select
            value={type}
            onChange={(e) => setType(e.target.value as AssetType)}
            className="mt-1 rounded border border-neutral-300 px-2 py-1"
          >
            {ASSET_TYPES.map((t) => (
              <option key={t} value={t}>
                {t}
              </option>
            ))}
          </select>
        </label>
        <button
          type="submit"
          disabled={submitting}
          className="rounded bg-neutral-900 px-3 py-1.5 text-sm text-white disabled:opacity-50"
        >
          Adicionar
        </button>
      </form>

      {error && <p className="text-sm text-red-600">{error}</p>}

      {loading ? (
        <p className="text-sm text-neutral-500">Carregando...</p>
      ) : assets.length === 0 ? (
        <p className="text-sm text-neutral-500">Nenhum ativo cadastrado ainda.</p>
      ) : (
        <table className="w-full text-left text-sm">
          <thead>
            <tr className="border-b border-neutral-200 text-neutral-500">
              <th className="py-2">Ticker</th>
              <th className="py-2">Nome</th>
              <th className="py-2">Tipo</th>
              <th className="py-2" />
            </tr>
          </thead>
          <tbody>
            {assets.map((asset) => (
              <tr key={asset.ID} className="border-b border-neutral-100">
                <td className="py-2 font-medium">{asset.Ticker}</td>
                <td className="py-2">{asset.Name}</td>
                <td className="py-2">{asset.Type}</td>
                <td className="py-2 text-right">
                  <button
                    onClick={() => handleDelete(asset.ID)}
                    className="text-red-600 hover:underline"
                  >
                    Remover
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
