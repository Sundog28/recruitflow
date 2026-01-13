import { useEffect, useMemo, useState } from "react";
import { api } from "../lib/api";
import type { Job } from "../lib/api";
import { clearToken, getToken } from "../lib/auth";
import { Button, Card, Input } from "../components/Card";
import { useNavigate } from "react-router-dom";

export default function JobsPage() {
  const nav = useNavigate();
  const token = useMemo(() => getToken(), []);
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState<string | null>(null);

  const [company, setCompany] = useState("");
  const [title, setTitle] = useState("");
  const [status, setStatus] = useState("Saved");

  async function load() {
    if (!token) {
      nav("/");
      return;
    }
    setErr(null);
    setLoading(true);
    try {
      const data = await api.listJobs(token);
      setJobs(data);
    } catch (e: any) {
      setErr(e?.message ?? "Failed to load jobs");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function create() {
    if (!token) return;
    setErr(null);
    try {
      const created = await api.createJob(token, { company, title, status });
      setJobs((prev) => [created, ...prev]);
      setCompany("");
      setTitle("");
      setStatus("Saved");
    } catch (e: any) {
      setErr(e?.message ?? "Failed to create job");
    }
  }

  function logout() {
    clearToken();
    nav("/");
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="mx-auto max-w-4xl space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold">Jobs</h1>
          <div className="flex gap-2">
            <Button variant="ghost" onClick={load}>Refresh</Button>
            <Button variant="ghost" onClick={logout}>Logout</Button>
          </div>
        </div>

        <Card>
          <h2 className="text-lg font-semibold mb-3">Add job</h2>
          <div className="grid gap-3 sm:grid-cols-3">
            <div className="sm:col-span-1">
              <label className="text-sm text-gray-600">Company</label>
              <Input value={company} onChange={(e) => setCompany(e.target.value)} />
            </div>
            <div className="sm:col-span-1">
              <label className="text-sm text-gray-600">Title</label>
              <Input value={title} onChange={(e) => setTitle(e.target.value)} />
            </div>
            <div className="sm:col-span-1">
              <label className="text-sm text-gray-600">Status</label>
              <select
                className="w-full rounded-xl border px-3 py-2"
                value={status}
                onChange={(e) => setStatus(e.target.value)}
              >
                <option>Saved</option>
                <option>Applied</option>
                <option>Interview</option>
                <option>Offer</option>
                <option>Rejected</option>
              </select>
            </div>
          </div>
          <div className="mt-4">
            <Button onClick={create} disabled={!company || !title}>
              Create
            </Button>
          </div>
          {err && (
            <div className="mt-4 rounded-xl border border-red-200 bg-red-50 p-3 text-sm text-red-700">
              {err}
            </div>
          )}
        </Card>

        <Card>
          <h2 className="text-lg font-semibold mb-3">Your jobs</h2>
          {loading ? (
            <div className="text-gray-600">Loading...</div>
          ) : jobs.length === 0 ? (
            <div className="text-gray-600">No jobs yet. Add one above.</div>
          ) : (
            <div className="space-y-3">
              {jobs.map((j) => (
                <div key={j.id} className="rounded-xl border p-4 bg-white">
                  <div className="flex items-start justify-between gap-4">
                    <div>
                      <div className="font-semibold">
                        {j.company} â€” {j.title}
                      </div>
                      <div className="text-sm text-gray-600">
                        Status: <span className="font-medium">{j.status}</span>
                      </div>
                    </div>
                    <div className="text-xs text-gray-500">#{j.id}</div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </Card>
      </div>
    </div>
  );
}


