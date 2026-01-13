const API_BASE = (import.meta.env.VITE_API_BASE ?? "http://localhost:8080") + "/v1";

export type AuthResponse = { token: string };

export type Job = {
  id: number;
  user_id: number;
  company: string;
  title: string;
  link: string;
  status: string;
  salary: string;
  notes: string;
  follow_up_date: string | null;
  created_at: string;
  updated_at: string;
};

export type CreateJobRequest = {
  company: string;
  title: string;
  link?: string;
  status?: string;
  salary?: string;
  notes?: string;
  follow_up_date?: string | null;
};

async function request<T>(
  path: string,
  opts: RequestInit = {},
  token?: string | null
): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(opts.headers as Record<string, string> | undefined),
  };

  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(`${API_BASE}${path}`, { ...opts, headers });

  // Try to read JSON if possible; fall back to text
  const text = await res.text();
  const data = text ? safeJson(text) : null;

  if (!res.ok) {
    const msg =
      (data && (data.error || data.message)) ||
      text ||
      `${res.status} ${res.statusText}`;
    throw new Error(msg);
  }

  return data as T;
}

function safeJson(s: string) {
  try {
    return JSON.parse(s);
  } catch {
    return s;
  }
}

export const api = {
  register(email: string, password: string) {
    return request<AuthResponse>("/auth/register", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
  },
  login(email: string, password: string) {
    return request<AuthResponse>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
  },
  me(token: string) {
    return request<{ user_id: number }>("/me", { method: "GET" }, token);
  },
  listJobs(token: string) {
    return request<Job[]>("/jobs", { method: "GET" }, token);
  },
  createJob(token: string, payload: CreateJobRequest) {
    return request<Job>(
      "/jobs",
      { method: "POST", body: JSON.stringify(payload) },
      token
    );
  },
};


