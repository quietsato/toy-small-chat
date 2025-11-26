import { useState } from "react";
import { login, createAccount } from "../api/account";
import "./Login.css";

interface LoginProps {
  onLoginSuccess: (token: string, username: string) => void;
}

export function Login(props: LoginProps) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isCreatingAccount, setIsCreatingAccount] = useState(false);

  async function handleLogin() {
    setError(null);
    try {
      const result = await login({ username, password });
      props.onLoginSuccess(result.token, result.username);
    } catch (err) {
      setError("ログインに失敗しました");
    }
  }

  async function handleCreateAccount() {
    setError(null);
    try {
      const result = await createAccount({ username, password });
      props.onLoginSuccess(result.token, result.username);
    } catch (err) {
      setError("アカウント作成に失敗しました");
    }
  }

  return (
    <div className="loginContainer">
      <div className="loginBox">
        <h2>{isCreatingAccount ? "アカウント作成" : "ログイン"}</h2>
        <div className="loginForm">
          <input
            type="text"
            placeholder="ユーザー名"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <input
            type="password"
            placeholder="パスワード"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          {error && <div className="loginError">{error}</div>}
          {isCreatingAccount ? (
            <div className="loginButtons">
              <button onClick={handleCreateAccount}>作成</button>
              <button onClick={() => setIsCreatingAccount(false)}>
                戻る
              </button>
            </div>
          ) : (
            <div className="loginButtons">
              <button onClick={handleLogin}>ログイン</button>
              <button onClick={() => setIsCreatingAccount(true)}>
                アカウント作成
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
