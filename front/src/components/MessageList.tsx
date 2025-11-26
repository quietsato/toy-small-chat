import { useEffect, useRef, useState } from "react";
import { Message, type MessageProps } from "./Message";
import "./MessageList.css";

export interface MessageListProps {
  messages: MessageProps[] | undefined;
}

export function MessageList(props: MessageListProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [shouldFollowLatestMessage, setShouldFollowLatestMessage] = useState(true);
  const prevMessagesRef = useRef<MessageProps[] | undefined>(undefined);

  const checkScrollPosition = () => {
    const container = containerRef.current;
    if (!container) return;

    const { scrollTop, scrollHeight, clientHeight } = container;
    const isNearBottom = scrollTop + clientHeight >= scrollHeight - 16;
    setShouldFollowLatestMessage(isNearBottom);
  };

  const scrollToLatestMessage = (animated: boolean = true) => {
    const container = containerRef.current;
    if (!container) return;

    container.scrollTo({
      top: container.scrollHeight,
      behavior: animated ? "smooth" : "instant",
    });
  };

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    checkScrollPosition();
    container.addEventListener("scroll", checkScrollPosition);

    return () => {
      container.removeEventListener("scroll", checkScrollPosition);
    };
  }, []);

  useEffect(() => {
    // 最新メッセージを追従している場合は自動スクロール
    if (shouldFollowLatestMessage && props.messages !== undefined) {
      // undefined -> array の遷移時はアニメーションなし、array -> array の場合はアニメーションあり
      const wasUndefined = prevMessagesRef.current === undefined;
      scrollToLatestMessage(!wasUndefined);
    }
    prevMessagesRef.current = props.messages;
    checkScrollPosition();
  }, [props.messages]);

  return (
    <div className="messageListWrapper">
      <div className="messageListContainer" ref={containerRef}>
        {props.messages === undefined ? (
          <div className="loadingContainer">
            <div className="loadingSpinner"></div>
            <p>メッセージを読み込んでいます...</p>
          </div>
        ) : (
          props.messages.map((message) => (
            <div key={message.id} className="message-container">
              <Message {...message} />
            </div>
          ))
        )}
      </div>
      {!shouldFollowLatestMessage && (
        <button className="scrollToBottomFab" onClick={() => scrollToLatestMessage()}>
          <svg
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M12 5L12 19M12 19L19 12M12 19L5 12"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
          </svg>
        </button>
      )}
    </div>
  );
}
