interface LoginInput {
  username: string;
  password: string;
}

interface LoginOutput {
  username: string;
  token: string;
}

export async function login(input: LoginInput): Promise<LoginOutput> {
  const res = await fetch(`http://localhost:18081/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(input),
  });

  if (!res.ok) {
    throw new Error("Login failed");
  }

  const resJson = await res.json();
  return resJson as LoginOutput;
}

interface CreateAccountInput {
  username: string;
  password: string;
}

interface CreateAccountOutput {
  username: string;
  token: string;
}

export async function createAccount(
  input: CreateAccountInput
): Promise<CreateAccountOutput> {
  const res = await fetch(`http://localhost:18081/accounts`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(input),
  });

  if (!res.ok) {
    throw new Error("Account creation failed");
  }

  const resJson = await res.json();
  return resJson as CreateAccountOutput;
}
