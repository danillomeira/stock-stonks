const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export interface User {
  ID: number;
  Email: string;
  Name: string;
  AvatarURL: string;
}

// Full-page navigation, not fetch: the browser needs to actually land on
// Google's consent screen and come back with a session cookie set.
export const googleLoginURL = `${API_URL}/api/v1/auth/google/login`;

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

export class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = "ApiError";
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_URL}${path}`, {
    ...init,
    // The session lives in an HttpOnly cookie set by the backend; without
    // this the browser won't send it cross-port (localhost:3000 -> :8080).
    credentials: "include",
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
  getMe: () => request<User>("/api/v1/auth/me"),
  logout: () => request<void>("/api/v1/auth/logout", { method: "POST" }),

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
