export type AuthResponse = {
  access_token: string;
  refresh_token: string;
  expires_at: string;
};

export type RegisterRequest = {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
};

export type LoginRequest = {
  email: string;
  password: string;
};

async function request<T>(
  path: string,
  options: RequestInit,
): Promise<T> {
  const response = await fetch(`/api/user${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
  });

  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || "Ошибка запроса");
  }

  return response.json() as Promise<T>;
}

export function register(data: RegisterRequest) {
  return request<AuthResponse>("/auth/register", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export function login(data: LoginRequest) {
  return request<AuthResponse>("/auth/login", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export function refresh(refreshToken: string) {
  return request<AuthResponse>("/auth/refresh", {
    method: "POST",
    body: JSON.stringify({
      refresh_token: refreshToken,
    }),
  });
}