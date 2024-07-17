const BACKEND_URL = process.env.NEXT_PUBLIC_BACKEND_URL;

export async function getOptionsForUser(email: string): Promise<any> {
  const res = await fetch(`${BACKEND_URL}/account/register-request`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email }),
  });
  const data = await res.json();
  return data.options.publicKey;
}

export async function registerUser(email: string, credential: any): Promise<any> {
  const res = await fetch(`${BACKEND_URL}/account/register/${email}`, {
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
