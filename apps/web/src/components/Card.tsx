import React from "react";

export function Card({ children }: { children: React.ReactNode }) {
  return (
    <div className="rounded-2xl border bg-white p-6 shadow-sm">{children}</div>
  );
}

export function Input(props: React.InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      {...props}
      className={
        "w-full rounded-xl border px-3 py-2 outline-none focus:ring " +
        (props.className ?? "")
      }
    />
  );
}

export function Button(
  props: React.ButtonHTMLAttributes<HTMLButtonElement> & { variant?: "primary" | "ghost" }
) {
  const v = props.variant ?? "primary";
  const base = "rounded-xl px-4 py-2 font-medium";
  const cls =
    v === "primary"
      ? base + " bg-black text-white hover:opacity-90"
      : base + " border bg-white hover:bg-gray-50";

  return <button {...props} className={cls + " " + (props.className ?? "")} />;
}
