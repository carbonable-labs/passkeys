import { BACKEND_URL } from "./index";

export async function getLoginOptionsForUser(email: string): Promise<any> {
  const res = await fetch(`${BACKEND_URL}/account/login-request`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email }),
  });
  const data = await res.json();
  return data.options.publicKey;
}

export async function loginUser(email: string, credential: any): Promise<any> {
  const res = await fetch(`${BACKEND_URL}/account/login?email=${email}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      ...credential,
    }),
  });

  console.log(res)
}
