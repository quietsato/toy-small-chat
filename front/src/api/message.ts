import { authenticatedFetch } from "./index";

interface CreateMessageInput {
  content: string;
  roomId: string;
}

export async function createMessage(input: CreateMessageInput) {
  await authenticatedFetch(`/rooms/${input.roomId}/messages`, {
    body: JSON.stringify({
      content: input.content,
    }),
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
  });
}
export interface GetMessagesInput {
  roomId: string;
}

export interface GetMessagesOutput {
  messages: Message[];
}

export interface Message {
  id: string;
  content: string;
  author: string;
  createdAt: string;
}

export async function getMessages(
  input: GetMessagesInput
): Promise<GetMessagesOutput> {
  const res = await authenticatedFetch(
    `/rooms/${input.roomId}/messages`,
    {
      method: "GET",
    }
  );

  const resJson = await res.json();

  // TODO: 型チェック
  return resJson as GetMessagesOutput;
}
