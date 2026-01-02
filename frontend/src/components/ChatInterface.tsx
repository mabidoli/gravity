"use client";

import { useState, useRef, useEffect } from "react";
import { motion } from "framer-motion";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Button } from "@/components/ui/button";
import { MessageBubble } from "./MessageBubble";
import { SourceIcon, getSourceLabel } from "./SourceIcon";
import { useGravityStore } from "@/store/useGravityStore";
import { useStreamItem, useOptimisticUpdate } from "@/hooks/useStream";
import { Send, ArrowLeft, MoreVertical, ExternalLink, Loader2 } from "lucide-react";
import { cn } from "@/lib/utils";

export function ChatInterface() {
  const [inputValue, setInputValue] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const {
    selectedItemId,
    sendMessage,
    setMobileViewingChat,
    openContextModal,
    setLocalItem,
    getLocalItem,
  } = useGravityStore();

  // Fetch item details from API
  const { data: apiItem, isLoading } = useStreamItem(selectedItemId);
  const { markAsRead } = useOptimisticUpdate();

  // Get item from local cache (for local mutations) or API
  const localItem = selectedItemId ? getLocalItem(selectedItemId) : undefined;
  const selectedItem = localItem || apiItem;

  // Sync API data to local cache when it arrives
  useEffect(() => {
    if (apiItem && !localItem) {
      setLocalItem(apiItem);
      // Mark as read when viewing
      markAsRead(apiItem.id);
    }
  }, [apiItem, localItem, setLocalItem, markAsRead]);

  // Auto-scroll to bottom when messages change
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [selectedItem?.messages]);

  const handleSend = () => {
    if (inputValue.trim() && selectedItem) {
      sendMessage(selectedItem.id, inputValue);
      setInputValue("");
      inputRef.current?.focus();
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  // Empty state - no item selected
  if (!selectedItemId) {
    return (
      <div className="h-full flex flex-col items-center justify-center text-slate-500">
        <div className="h-20 w-20 rounded-full bg-slate-800/50 flex items-center justify-center mb-4">
          <Send size={32} className="opacity-50" />
        </div>
        <h2 className="text-lg font-medium text-slate-400 mb-2">
          Select a conversation
        </h2>
        <p className="text-sm">
          Choose an item from the stream to view the conversation
        </p>
      </div>
    );
  }

  // Loading state
  if (isLoading && !selectedItem) {
    return (
      <div className="h-full flex flex-col items-center justify-center text-slate-500">
        <Loader2 size={32} className="mb-2 animate-spin" />
        <p className="text-sm">Loading conversation...</p>
      </div>
    );
  }

  // No item found
  if (!selectedItem) {
    return (
      <div className="h-full flex flex-col items-center justify-center text-slate-500">
        <div className="h-20 w-20 rounded-full bg-slate-800/50 flex items-center justify-center mb-4">
          <Send size={32} className="opacity-50" />
        </div>
        <h2 className="text-lg font-medium text-slate-400 mb-2">
          Conversation not found
        </h2>
        <p className="text-sm">
          This item may have been removed or is unavailable
        </p>
      </div>
    );
  }

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="h-full flex flex-col"
    >
      {/* Header */}
      <div className="flex items-center gap-3 p-4 border-b border-white/5 bg-slate-900/30">
        {/* Back Button (Mobile) */}
        <button
          onClick={() => setMobileViewingChat(false)}
          className="lg:hidden p-2 -ml-2 rounded-lg hover:bg-white/5 transition-colors"
        >
          <ArrowLeft size={20} className="text-slate-400" />
        </button>

        {/* Source Icon */}
        <SourceIcon source={selectedItem.source} />

        {/* Title & Meta */}
        <div className="flex-1 min-w-0">
          <h2 className="font-medium text-white truncate">
            {selectedItem.title}
          </h2>
          <div className="flex items-center gap-2 text-xs text-slate-500">
            <span>{getSourceLabel(selectedItem.source)}</span>
            <span>•</span>
            <span>
              {selectedItem.participants.length} participant
              {selectedItem.participants.length !== 1 ? "s" : ""}
            </span>
          </div>
        </div>

        {/* Actions */}
        <Button variant="ghost" size="icon">
          <ExternalLink size={18} className="text-slate-400" />
        </Button>
        <Button variant="ghost" size="icon">
          <MoreVertical size={18} className="text-slate-400" />
        </Button>
      </div>

      {/* Messages */}
      <ScrollArea className="flex-1 p-4">
        <div className="space-y-6 max-w-3xl mx-auto">
          {selectedItem.messages.map((message) => (
            <MessageBubble
              key={message.id}
              message={message}
              itemId={selectedItem.id}
              onOpenContext={() => openContextModal(message)}
            />
          ))}
          <div ref={messagesEndRef} />
        </div>
      </ScrollArea>

      {/* Input Area */}
      <div className="p-4 border-t border-white/5 bg-slate-900/30">
        <div className="max-w-3xl mx-auto">
          <div className="flex items-center gap-3">
            <div className="flex-1 relative">
              <input
                ref={inputRef}
                type="text"
                value={inputValue}
                onChange={(e) => setInputValue(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder={getPlaceholder(selectedItem.source)}
                className={cn(
                  "w-full px-4 py-3 rounded-xl text-sm",
                  "bg-slate-800/50 border border-white/10",
                  "text-white placeholder:text-slate-500",
                  "focus:outline-none focus:border-teal-500/50 focus:bg-slate-800/70",
                  "transition-all duration-200"
                )}
              />
            </div>
            <Button
              onClick={handleSend}
              disabled={!inputValue.trim()}
              className="h-12 w-12 rounded-xl"
            >
              <Send size={18} />
            </Button>
          </div>

          {/* Hint Text */}
          <p className="text-xs text-slate-500 mt-2 text-center">
            {getHintText(selectedItem.source)}
          </p>
        </div>
      </div>
    </motion.div>
  );
}

function getPlaceholder(source: string): string {
  switch (source) {
    case "calendar":
      return "Type a message or command (e.g., 'Reschedule to 4 PM')...";
    case "youtube":
    case "linkedin":
    case "twitter":
      return "Ask about this content (e.g., 'Summarize this')...";
    default:
      return "Type your message...";
  }
}

function getHintText(source: string): string {
  switch (source) {
    case "calendar":
      return "You can reschedule or modify this event using natural language";
    case "youtube":
      return "Ask me to summarize, analyze, or find key points in this video";
    case "linkedin":
    case "twitter":
      return "I can help you understand and engage with this content";
    default:
      return "Press Enter to send, Shift+Enter for new line";
  }
}
