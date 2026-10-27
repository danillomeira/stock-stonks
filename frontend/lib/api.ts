const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export type AssetType = "STOCK" | "FII" | "ETF" | "BDR";

export interface Asset {
  ID: number;
  Ticker: string;
  Name: string;
  Type: AssetType;
  CreatedAt: string;
  UpdatedAt: string;
}

export type OperationType = "BUY" | "SELL";

export interface Operation {
  ID: number;
  AssetID: number;
  Asset?: Asset;
  Date: string;
  Type: OperationType;
  Quantity: string;
  UnitPrice: string;
  Fees: string;
  Broker: string;
}

export interface Position {
  AssetID: number;
  Quantity: string;
  AverageCost: string;
  RealizedProfit: string;
}

class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = "ApiError";
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_URL}${path}`, {
    ...init,
    headers: { "Content-Type": "application/json", ...init?.headers },
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new ApiError(res.status, body.error ?? res.statusText);
  }
  if (res.status === 204) {
    return undefined as T;
  }
  return res.json() as Promise<T>;
}

export const api = {
  listAssets: () => request<Asset[]>("/api/v1/assets"),
  createAsset: (input: { ticker: string; name: string; type: AssetType }) =>
    request<Asset>("/api/v1/assets", { method: "POST", body: JSON.stringify(input) }),
  deleteAsset: (id: number) => request<void>(`/api/v1/assets/${id}`, { method: "DELETE" }),

  listOperations: (params?: { assetId?: number }) => {
    const qs = params?.assetId ? `?asset_id=${params.assetId}` : "";
    return request<Operation[]>(`/api/v1/operations${qs}`);
  },
  createOperation: (input: {
    asset_id: number;
    date: string;
    type: OperationType;
    quantity: string;
    unit_price: string;
    fees: string;
    broker: string;
  }) => request<Operation>("/api/v1/operations", { method: "POST", body: JSON.stringify(input) }),
  deleteOperation: (id: number) => request<void>(`/api/v1/operations/${id}`, { method: "DELETE" }),

  listPositions: () => request<Position[]>("/api/v1/positions"),
};
