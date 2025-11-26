import { useState } from "react";
import "./RoomList.css";

interface RoomListItem {
  name: string;
  id: string;
}

interface RoomListProps {
  rooms: RoomListItem[];
  selectedRoomId: string | undefined;
  onRoomSelected: (roomId: string) => void;
  onCreateRoom: (roomName: string) => void;
}

export function RoomList(props: RoomListProps) {
  const [isCreating, setIsCreating] = useState(false);
  const [roomName, setRoomName] = useState("");

  function roomListItemClassName(roomId: string) {
    const classNames = ["roomListItem"];
    if (props.selectedRoomId == roomId) {
      classNames.push("selectedRoomListItem");
    }

    return classNames.join(" ");
  }

  function handleCreateRoom() {
    if (roomName.trim()) {
      props.onCreateRoom(roomName);
      setRoomName("");
      setIsCreating(false);
    }
  }

  function handleKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key === "Enter" && e.shiftKey) {
      e.preventDefault();
      handleCreateRoom();
    }
  }

  return (
    <div className="roomListContainer">
      <div className="roomListHeader">
        {isCreating ? (
          <div className="createRoomForm">
            <input
              type="text"
              placeholder="ルーム名"
              value={roomName}
              onChange={(e) => setRoomName(e.target.value)}
              onKeyDown={handleKeyDown}
              autoFocus
            />
            <div className="createRoomButtons">
              <button onClick={handleCreateRoom}>作成</button>
              <button onClick={() => setIsCreating(false)}>キャンセル</button>
            </div>
          </div>
        ) : (
          <button className="createRoomButton" onClick={() => setIsCreating(true)}>
            + ルーム作成
          </button>
        )}
      </div>
      <div className="roomListItems">
        {props.rooms.map((room) => (
          <div
            onClick={() => {
              if (props.selectedRoomId !== room.id) {
                props.onRoomSelected(room.id);
              }
            }}
            key={room.id}
            className={roomListItemClassName(room.id)}
          >
            <p>{room.name}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
