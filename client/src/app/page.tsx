"use client";
import ChatBar from "@/components/ChatBar";
import ChatBubble from "@/components/ChatBubble";
import Image from "next/image";
import { useState } from "react";

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col  bg-blue-100 px-4">
      <div className={"sticky top-0 bg-gray-200 py-2 px-8"}></div>
      <ChatBubble message={"test"} isReceived={true} />
      <ChatBubble message={"test test est"} isReceived={false} />
      <ChatBubble message={"test"} isReceived={true} />
      <ChatBubble message={"test test est"} isReceived={false} />
      <ChatBubble message={"test"} isReceived={true} />
      <ChatBubble message={"test test est"} isReceived={false} />
      <ChatBubble message={"test"} isReceived={true} />
      <ChatBubble message={"test test est"} isReceived={false} />
      <ChatBubble message={"test"} isReceived={true} />
      <ChatBubble message={"test test est"} isReceived={false} />
      <ChatBubble message={"test"} isReceived={true} />
      <ChatBubble message={"test test est"} isReceived={false} />
      <ChatBubble message={"test"} isReceived={true} />
      <ChatBubble message={"test"} isReceived={true} />
      <ChatBubble message={"test test est"} isReceived={false} />
      {/* <div className="flex min-h-screen flex-col z-10 w-full max-w-5xl items-center justify-between font-mono text-sm lg:flex">
        <p className="fixed left-0 top-0 flex w-full justify-center border-b border-gray-300 bg-gradient-to-b from-zinc-200 pb-6 pt-8 backdrop-blur-2xl dark:border-neutral-800 dark:bg-zinc-800/30 dark:from-inherit lg:static lg:w-auto  lg:rounded-xl lg:border lg:bg-gray-200 lg:p-4 lg:dark:bg-zinc-800/30">
          Get started by editing&nbsp;
          <code className="font-mono font-bold">src/app/page.tsx</code>
        </p>
      </div> */}
      <ChatBar />
    </main>
  );
}
