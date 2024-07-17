"use client";
import { useCallback } from "react";
import { Button } from "../components/Button";

export default function Home() {
  const handleRegister = useCallback(async () => {
  }, []);
  const handleLogin = useCallback(async () => {
  }, []);
  return (
    <main className="">
      <div className="py-24">
        <div className="flex my-6 items-center justify-center font-mono text-sm">
          <input type="text" placeholder="Username" className="w-250px justify-center p-3 border rounded-md bg-transparent text-white" />
        </div>
        <div className="flex my-6 items-center justify-center font-mono text-sm">
          <Button onClick={handleRegister}>Register</Button>
          <Button onClick={handleLogin}>Login</Button>
        </div>
      </div>
      <div className="min-h-screen flex"></div>
    </main>
  );
}
