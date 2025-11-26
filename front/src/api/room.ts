import { authenticatedFetch } from "./index";

interface CreateRoomInput {
  name: string;
}

export async function createRoom(input: CreateRoomInput) {
  await authenticatedFetch(`/rooms`, {
    body: JSON.stringify({
      name: input.name,
    }),
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  });
}

export interface GetRoomsOutput {
  rooms: Room[];
}

export interface Room {
  id: string;
  name: string;
}

export async function getRooms(): Promise<GetRoomsOutput> {
  const res = await authenticatedFetch(`/rooms`, {
    method: "GET",
  });

  const resJson = await res.json();

  // TODO: 型チェック
  return resJson as GetRoomsOutput;
}
