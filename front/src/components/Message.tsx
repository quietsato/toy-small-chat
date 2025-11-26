import "./Message.css";

export interface MessageProps {
  id: string;
  content: string;
  author: string;
  createdAt: string;
}

export function Message(props: MessageProps) {
  return (
    <div className="container">
      <div className="meta">
        <p className="author">{props.author}</p>
        <p className="createdAt">{props.createdAt}</p>
      </div>
      <p className="content">{props.content}</p>
    </div>
  );
}
