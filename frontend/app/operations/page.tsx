"use client";

import { useEffect, useState } from "react";
import { api, type Asset, type Operation, type OperationType } from "@/lib/api";

export default function OperationsPage() {
  const [operations, setOperations] = useState<Operation[]>([]);
  const [assets, setAssets] = useState<Asset[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [assetId, setAssetId] = useState<number | "">("");
  const [date, setDate] = useState("");
  const [type, setType] = useState<OperationType>("BUY");
  const [quantity, setQuantity] = useState("");
  const [unitPrice, setUnitPrice] = useState("");
  const [fees, setFees] = useState("0");
  const [broker, setBroker] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function load() {
    setLoading(true);
    try {
      const [ops, assetList] = await Promise.all([api.listOperations(), api.listAssets()]);
      setOperations(ops);
      setAssets(assetList);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Falha ao carregar operações");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, []);

  function assetLabel(id: number) {
    return assets.find((a) => a.ID === id)?.Ticker ?? `#${id}`;
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (assetId === "") return;
    setSubmitting(true);
    try {
      await api.createOperation({
        asset_id: assetId,
        date: new Date(date).toISOString(),
        type,
        quantity,
        unit_price: unitPrice,
        fees,
        broker,
      });
      setQuantity("");
      setUnitPrice("");
      setFees("0");
      setBroker("");
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Falha ao registrar operação");
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(id: number) {
    try {
      await api.deleteOperation(id);
      await load();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Falha ao remover operação");
    }
  }

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-2xl font-semibold">Operações</h1>

      {assets.length === 0 && !loading ? (
        <p className="text-sm text-neutral-500">
          Cadastre um ativo primeiro na página{" "}
          <a href="/assets" className="underline">
            Ativos
          </a>
          .
        </p>
      ) : (
        <form
          onSubmit={handleSubmit}
          className="flex flex-wrap items-end gap-3 rounded-lg border border-neutral-200 bg-white p-4"
        >
          <label className="flex flex-col text-sm">
            Ativo
            <select
              required
              value={assetId}
              onChange={(e) => setAssetId(Number(e.target.value))}
              className="mt-1 rounded border border-neutral-300 px-2 py-1"
            >
              <option value="" disabled>
                Selecione
              </option>
              {assets.map((a) => (
                <option key={a.ID} value={a.ID}>
                  {a.Ticker}
                </option>
              ))}
            </select>
          </label>
          <label className="flex flex-col text-sm">
            Data
            <input
              required
              type="date"
              value={date}
              onChange={(e) => setDate(e.target.value)}
              className="mt-1 rounded border border-neutral-300 px-2 py-1"
            />
          </label>
          <label className="flex flex-col text-sm">
            Tipo
            <select
              value={type}
              onChange={(e) => setType(e.target.value as OperationType)}
              className="mt-1 rounded border border-neutral-300 px-2 py-1"
            >
              <option value="BUY">Compra</option>
              <option value="SELL">Venda</option>
            </select>
          </label>
          <label className="flex flex-col text-sm">
            Quantidade
            <input
              required
              value={quantity}
              onChange={(e) => setQuantity(e.target.value)}
              placeholder="100"
              className="mt-1 w-24 rounded border border-neutral-300 px-2 py-1"
            />
          </label>
          <label className="flex flex-col text-sm">
            Preço unit.
            <input
              required
              value={unitPrice}
              onChange={(e) => setUnitPrice(e.target.value)}
              placeholder="30.00"
              className="mt-1 w-28 rounded border border-neutral-300 px-2 py-1"
            />
          </label>
          <label className="flex flex-col text-sm">
            Taxas
            <input
              value={fees}
              onChange={(e) => setFees(e.target.value)}
              placeholder="5.00"
              className="mt-1 w-24 rounded border border-neutral-300 px-2 py-1"
            />
          </label>
          <label className="flex flex-col text-sm">
            Corretora
            <input
              value={broker}
              onChange={(e) => setBroker(e.target.value)}
              placeholder="XP"
              className="mt-1 w-28 rounded border border-neutral-300 px-2 py-1"
            />
          </label>
          <button
            type="submit"
            disabled={submitting}
            className="rounded bg-neutral-900 px-3 py-1.5 text-sm text-white disabled:opacity-50"
          >
            Registrar
          </button>
        </form>
      )}

      {error && <p className="text-sm text-red-600">{error}</p>}

      {loading ? (
        <p className="text-sm text-neutral-500">Carregando...</p>
      ) : operations.length === 0 ? (
        <p className="text-sm text-neutral-500">Nenhuma operação registrada ainda.</p>
      ) : (
        <table className="w-full text-left text-sm">
          <thead>
            <tr className="border-b border-neutral-200 text-neutral-500">
              <th className="py-2">Data</th>
              <th className="py-2">Ativo</th>
              <th className="py-2">Tipo</th>
              <th className="py-2">Qtd.</th>
              <th className="py-2">Preço</th>
              <th className="py-2">Taxas</th>
              <th className="py-2">Corretora</th>
              <th className="py-2" />
            </tr>
          </thead>
          <tbody>
            {operations.map((op) => (
              <tr key={op.ID} className="border-b border-neutral-100">
                <td className="py-2">{new Date(op.Date).toLocaleDateString("pt-BR")}</td>
                <td className="py-2">{assetLabel(op.AssetID)}</td>
                <td className="py-2">{op.Type === "BUY" ? "Compra" : "Venda"}</td>
                <td className="py-2">{op.Quantity}</td>
                <td className="py-2">{op.UnitPrice}</td>
                <td className="py-2">{op.Fees}</td>
                <td className="py-2">{op.Broker}</td>
                <td className="py-2 text-right">
                  <button onClick={() => handleDelete(op.ID)} className="text-red-600 hover:underline">
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
