import { useState } from "react";
import { api } from "../lib/api";
import { setToken } from "../lib/auth";
import { Button, Card, Input } from "../components/Card";
import { useNavigate } from "react-router-dom";

export default function AuthPage() {
  const nav = useNavigate();
  const [mode, setMode] = useState<"login" | "register">("login");
  const [email, setEmail] = useState("you@example.com");
  const [password, setPassword] = useState("password123");
  const [err, setErr] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  async function submit() {
    setErr(null);
    setBusy(true);
    try {
      const res =
        mode === "register"
          ? await api.register(email, password)
          : await api.login(email, password);

      setToken(res.token);
      nav("/jobs");
    } catch (e: any) {
      setErr(e?.message ?? "Something went wrong");
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center p-6">
      <div className="w-full max-w-md space-y-4">
        <h1 className="text-3xl font-bold">RecruitFlow</h1>
        <p className="text-gray-600">
          Track applications, interviews, and follow-ups.
        </p>

        <Card>
          <div className="flex gap-2 mb-4">
            <Button
              variant={mode === "login" ? "primary" : "ghost"}
              onClick={() => setMode("login")}
              type="button"
            >
              Login
            </Button>
            <Button
              variant={mode === "register" ? "primary" : "ghost"}
              onClick={() => setMode("register")}
              type="button"
            >
              Register
            </Button>
          </div>

          <div className="space-y-3">
            <div>
              <label className="text-sm text-gray-600">Email</label>
              <Input value={email} onChange={(e) => setEmail(e.target.value)} />
            </div>
            <div>
              <label className="text-sm text-gray-600">Password</label>
              <Input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>

            {err && (
              <div className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm text-red-700">
                {err}
              </div>
            )}

            <Button onClick={submit} disabled={busy} className="w-full">
              {busy ? "Working..." : mode === "register" ? "Create account" : "Sign in"}
            </Button>
          </div>
        </Card>

        <div className="text-xs text-gray-500">
          API: http://localhost:8080 • Web: http://localhost:5173
        </div>
      </div>
    </div>
  );
}
