import React from "react";

export default function ChatBubble({
  message,
  isReceived,
}: {
  message: string;
  isReceived: boolean;
}) {
  const bubbleStyles = isReceived
    ? "bg-gray-200 text-gray-700 rounded-tl-full rounded-br-full rounded-tr-full px-4 py-2"
    : "bg-blue-500 text-white rounded-tr-full rounded-bl-full rounded-tl-full px-4 py-2";

  return (
    <div className={`flex p-1 ${isReceived ? "justify-start" : "justify-end"}`}>
      <div className={`max-w-xs ${bubbleStyles}`}>{message}</div>
    </div>
  );
}
