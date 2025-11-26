import { Mutex } from "async-mutex";

const API_BASE_URL = "http://localhost:18081";

interface FetchOptions extends RequestInit {
  headers?: Record<string, string>;
}

// Mutex to ensure logout process runs only once
const logoutMutex = new Mutex();
let isLoggingOut = false;

/**
 * Fetch wrapper that automatically attaches JWT token
 */
export async function authenticatedFetch(
  path: string,
  options: FetchOptions = {}
): Promise<Response> {
  const token = localStorage.getItem("token");

  const headers: Record<string, string> = {
    ...options.headers,
  };

  // Add Authorization header if token exists
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers,
  });

  // Handle 401 Unauthorized error
  if (response.status === 401) {
    await logoutMutex.runExclusive(async () => {
      // If already logging out, skip
      if (isLoggingOut) {
        return;
      }

      // Mark as logging out
      isLoggingOut = true;

      // Clear authentication data
      localStorage.removeItem("token");
      localStorage.removeItem("username");

      // Show message to user
      alert("ログイン状態が失効しました。再度ログインしてください。");

      // Reload page to return to login screen
      window.location.reload();
    });

    // This will never execute if reload happened, but TypeScript needs it
    throw new Error("Unauthorized");
  }

  return response;
}
