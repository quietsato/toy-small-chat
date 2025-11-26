import { useEffect, useRef, useState } from "react";
import { createMessage, getMessages, type Message } from "./api/message";
import { getRooms, createRoom, type Room } from "./api/room";
import "./App.css";
import { MessageList } from "./components/MessageList";
import { RoomList } from "./components/RoomList";
import { Login } from "./components/Login";

interface AppState {
  roomId: string | undefined;
  rooms: Room[];
  messages: Message[] | undefined;
  token: string | null;
  username: string | null;
}

function App() {
  const [appState, setAppState] = useState<AppState>({
    roomId: undefined,
    rooms: [],
    messages: undefined,
    token: localStorage.getItem("token"),
    username: localStorage.getItem("username"),
  });

  const [leftPaneWidth, setLeftPaneWidth] = useState(300);
  const isResizingRef = useRef(false);
  const draftTextInputRef = useRef<HTMLTextAreaElement>(null);

  const handleMouseDown = () => {
    isResizingRef.current = true;
    document.body.style.cursor = 'col-resize';
  };

  const handleMouseMove = (e: MouseEvent) => {
    if (!isResizingRef.current) return;
    const newWidth = e.clientX;
    if (newWidth >= 200 && newWidth <= 600) {
      setLeftPaneWidth(newWidth);
    }
  };

  const handleMouseUp = () => {
    isResizingRef.current = false;
    document.body.style.cursor = 'default';
  };

  useEffect(() => {
    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);
    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };
  }, []);

  function handleLoginSuccess(token: string, username: string) {
    localStorage.setItem("token", token);
    localStorage.setItem("username", username);
    setAppState((s) => ({ ...s, token, username }));
  }

  function handleLogout() {
    if (window.confirm("ログアウトしますか？")) {
      localStorage.removeItem("token");
      localStorage.removeItem("username");
      setAppState((s) => ({ ...s, token: null, username: null, roomId: undefined, rooms: [], messages: undefined }));
    }
  }

  function isSendButtonEnabled() {
    return appState.roomId !== null;
  }

  function onSendButtonClick() {
    const roomId = appState.roomId;
    if (roomId === undefined) {
      return;
    }

    const currentDraftTextInput = draftTextInputRef.current!; // TODO: 直す
    const content = currentDraftTextInput.value;

    createMessage({
      content: content,
      roomId: roomId,
    }).then(updateMessageList);

    currentDraftTextInput.value = "";
  }

  function handleKeyDown(e: React.KeyboardEvent<HTMLTextAreaElement>) {
    // Shift+Enterで送信、Enterキーのみで改行
    if (e.key === "Enter" && e.shiftKey) {
      e.preventDefault();
      onSendButtonClick();
    }
  }

  // ルームが選択された時の副作用
  useEffect(() => {
    if (appState.roomId !== undefined) {
      setAppState((s) => ({ ...s, messages: undefined }));
      updateMessageList();
    }
  }, [appState.roomId]);

  function onRoomSelected(roomId: string) {
    setAppState((s) => ({ ...s, roomId, messages: undefined }));
  }

  async function onCreateRoom(roomName: string) {
    await createRoom({ name: roomName });
    await updateRoomList();
  }

  async function updateMessageList() {
    const currentRoomId = appState.roomId;
    if (currentRoomId === undefined) {
      return;
    }

    const res = await getMessages({ roomId: currentRoomId });
    setAppState((s) => {
      if (s.roomId !== currentRoomId) {
        return s; // 現在選択中のルームが変わっていた場合はステートを更新しない
      }
      return {
        ...s,
        roomId: currentRoomId,
        messages: res.messages,
      };
    });
  }
  async function updateRoomList() {
    // const currentRoomId = appState.roomId;
    const res = await getRooms();
    setAppState((s) => ({
      ...s,
      roomId: s.roomId === undefined ? res.rooms.at(0)?.id : s.roomId,
      rooms: res.rooms,
    }));
  }

  // 5秒おきにメッセージリスト・ルームリストを更新する
  useEffect(() => {
    // トークンがない場合は自動fetchを実行しない
    if (!appState.token) {
      return;
    }

    // initial run
    updateRoomList();

    const interval = setInterval(() => {
      updateMessageList();
      updateRoomList();
    }, 5000);

    return () => clearInterval(interval);
  }, [appState.token]);

  if (!appState.token) {
    return <Login onLoginSuccess={handleLoginSuccess} />;
  }

  return (
    <>
      <div className="appContainer">
        <div className="appLeftPane" style={{ width: `${leftPaneWidth}px` }}>
          <RoomList
            rooms={appState.rooms}
            selectedRoomId={appState.roomId}
            onRoomSelected={onRoomSelected}
            onCreateRoom={onCreateRoom}
          />
          <div className="userInfo">
            <span className="username">{appState.username}</span>
          </div>
          <button onClick={handleLogout} className="logoutButton">
            ログアウト
          </button>
        </div>
        <div className="resizeHandle" onMouseDown={handleMouseDown}></div>
        <div className="appRightPane">
          <MessageList messages={appState.messages} />
          <div className="messageInputContainer">
            <textarea
              ref={draftTextInputRef}
              className="messageInput"
              placeholder="メッセージを入力... (Shift+Enterで送信)"
              onKeyDown={handleKeyDown}
              rows={3}
            />
            <button
              disabled={!isSendButtonEnabled()}
              onClick={onSendButtonClick}
              className="sendButton"
            >
              送信
            </button>
          </div>
        </div>
      </div>
    </>
  );
}

export default App;
