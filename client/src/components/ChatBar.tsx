import React, { useState } from "react";

export default function ChatBar() {
  const [message, setMessage] = useState("");

  const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setMessage(event.target.value);
  };

  const onSendBtnClicked = () => {};

  return (
    <div className={"sticky bottom-0 bg-gray-200 py-2 px-8 rounded-3xl"}>
      <input
        className={
          "px-4 py-2 border border-gray-300 rounded-l-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        }
        type="text"
        value={message}
        onChange={handleInputChange}
        placeholder="Type your message..."
      />
      <button
        className={
          message.trim().length === 0
            ? "px-4 py-2 bg-red-100 text-white rounded-r-md hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-opacity-50"
            : "px-4 py-2 bg-blue-500 text-white rounded-r-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50"
        }
        onSubmit={onSendBtnClicked}
        disabled={message.trim().length === 0}
      >
        Send
      </button>
    </div>
  );
}
