"use client";
import { useCallback, useState } from "react";
import { Button } from "../components/Button";
import { getOptionsForUser, registerUser } from "../actions/register";
import { startRegistration } from "@simplewebauthn/browser";

export default function Home() {
  const [email, setEmail] = useState("");

  const handleRegister = useCallback(async () => {
    try {

      const opts = await getOptionsForUser(email);
      const credentials = await startRegistration(opts);
      try {
        await registerUser(email, credentials);
        setEmail("");
      } catch (e) {
        // Hanlde error here 
      }
    } catch (e) {
      // Handle opts registration error here
      // and navigator.credentials.create error here
    }
  }, [email, setEmail]);

  const handleLogin = useCallback(async () => {
  }, []);

  const handleInputChange = useCallback(async (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value)
  }, [setEmail]);

  return (
    <main className="">
      <div className="py-24">
        <div className="flex my-6 items-center justify-center font-mono text-sm">
          <input type="text" placeholder="Username" className="w-250px justify-center p-3 border rounded-md bg-transparent text-white" onChange={handleInputChange} />
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
