"use client";

import { Message } from "@/types";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { EventCard } from "./EventCard";
import { ContentCard } from "./ContentCard";
import { InsightPanel } from "./InsightPanel";
import { formatTimestamp, formatFileSize } from "@/lib/utils";
import { cn } from "@/lib/utils";
import { Paperclip, FileText, Image as ImageIcon, File } from "lucide-react";

interface MessageBubbleProps {
  message: Message;
  itemId: string;
  onOpenContext: () => void;
}

export function MessageBubble({ message, itemId, onOpenContext }: MessageBubbleProps) {
  const isUser = message.sender === "user";
  const isSystem = message.sender === "system";

  const getFileIcon = (type: string) => {
    if (type.includes("pdf")) return <FileText size={16} className="text-red-400" />;
    if (type.includes("image")) return <ImageIcon size={16} className="text-blue-400" />;
    return <File size={16} className="text-slate-400" />;
  };

  // Render Event Card for calendar events
  if (message.type === "event" && message.eventDetails) {
    return (
      <div className="flex gap-3 animate-fade-in">
        <div className="h-8 w-8 rounded-full bg-teal-500/20 flex items-center justify-center shrink-0">
          <span className="text-xs text-teal-400">SYS</span>
        </div>
        <div className="flex-1">
          <div className="text-xs text-slate-500 mb-2">
            {formatTimestamp(message.timestamp)}
          </div>
          <EventCard event={message.eventDetails} />
        </div>
      </div>
    );
  }

  // Render Content Card for social media
  if (message.type === "social" && message.socialContent) {
    return (
      <div className="flex gap-3 animate-fade-in">
        <div className="h-8 w-8 rounded-full bg-slate-700 flex items-center justify-center shrink-0">
          <span className="text-xs text-slate-400">
            {message.socialContent.platform.charAt(0).toUpperCase()}
          </span>
        </div>
        <div className="flex-1">
          <div className="text-xs text-slate-500 mb-2">
            {formatTimestamp(message.timestamp)}
          </div>
          <ContentCard content={message.socialContent} />
          {/* AI Insights */}
          {message.aiInsights?.map((insight) => (
            <InsightPanel
              key={insight.id}
              insight={insight}
              itemId={itemId}
              messageId={message.id}
            />
          ))}
        </div>
      </div>
    );
  }

  // Regular text message
  return (
    <div
      className={cn(
        "flex gap-3 animate-fade-in",
        isUser && "flex-row-reverse"
      )}
    >
      {/* Avatar */}
      {!isUser && (
        <Avatar className="h-8 w-8 shrink-0">
          <AvatarFallback className={cn(
            isSystem ? "bg-teal-500/20 text-teal-400" : "bg-slate-700 text-slate-300"
          )}>
            {isSystem ? "AI" : message.senderInfo?.name.charAt(0) || "?"}
          </AvatarFallback>
        </Avatar>
      )}

      <div className={cn("flex-1 max-w-[80%]", isUser && "flex flex-col items-end")}>
        {/* Sender Info & Timestamp */}
        <div className={cn(
          "flex items-center gap-2 mb-1",
          isUser && "flex-row-reverse"
        )}>
          {!isUser && message.senderInfo && (
            <span className="text-xs font-medium text-slate-300">
              {message.senderInfo.name}
            </span>
          )}
          <span className="text-xs text-slate-500">
            {formatTimestamp(message.timestamp)}
          </span>
        </div>

        {/* Message Content */}
        <div
          onClick={message.fullContent ? onOpenContext : undefined}
          className={cn(
            "px-4 py-3 rounded-2xl text-sm cursor-pointer transition-all duration-200",
            isUser
              ? "bg-teal-600 text-white rounded-tr-md hover:bg-teal-500"
              : isSystem
              ? "bg-slate-800/80 text-slate-200 rounded-tl-md border border-white/5"
              : "bg-slate-800/50 text-slate-200 rounded-tl-md hover:bg-slate-800/70"
          )}
        >
          <div className="message-content whitespace-pre-wrap">
            {message.content}
          </div>
        </div>

        {/* Attachments */}
        {message.attachments && message.attachments.length > 0 && (
          <div className="mt-2 space-y-1">
            {message.attachments.map((attachment) => (
              <div
                key={attachment.id}
                className="flex items-center gap-2 px-3 py-2 rounded-lg bg-slate-800/50 border border-white/5 hover:bg-slate-800/70 transition-colors cursor-pointer"
              >
                {getFileIcon(attachment.type)}
                <div className="flex-1 min-w-0">
                  <div className="text-xs text-slate-300 truncate">
                    {attachment.name}
                  </div>
                  <div className="text-[10px] text-slate-500">
                    {formatFileSize(attachment.size)}
                  </div>
                </div>
                <Paperclip size={12} className="text-slate-500" />
              </div>
            ))}
          </div>
        )}

        {/* AI Insights */}
        {message.aiInsights?.map((insight) => (
          <InsightPanel
            key={insight.id}
            insight={insight}
            itemId={itemId}
            messageId={message.id}
          />
        ))}
      </div>

      {/* User Avatar Placeholder */}
      {isUser && (
        <Avatar className="h-8 w-8 shrink-0">
          <AvatarFallback className="bg-teal-600 text-white">
            U
          </AvatarFallback>
        </Avatar>
      )}
    </div>
  );
}
