"use client";

import { useGravityStore } from "@/store/useGravityStore";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { formatTimestamp, formatFileSize } from "@/lib/utils";
import {
  ExternalLink,
  Download,
  FileText,
  Image as ImageIcon,
  File,
  Paperclip,
} from "lucide-react";

export function FullContextModal() {
  const { isContextModalOpen, contextModalMessage, closeContextModal } =
    useGravityStore();

  if (!contextModalMessage) return null;

  const getFileIcon = (type: string) => {
    if (type.includes("pdf"))
      return <FileText size={20} className="text-red-400" />;
    if (type.includes("image"))
      return <ImageIcon size={20} className="text-blue-400" />;
    return <File size={20} className="text-slate-400" />;
  };

  const hasFullContent = contextModalMessage.fullContent;
  const hasAttachments =
    contextModalMessage.attachments &&
    contextModalMessage.attachments.length > 0;

  return (
    <Dialog open={isContextModalOpen} onOpenChange={closeContextModal}>
      <DialogContent className="max-w-3xl max-h-[85vh] flex flex-col">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <span>Full Context</span>
            {contextModalMessage.senderInfo && (
              <span className="text-sm font-normal text-slate-400">
                from {contextModalMessage.senderInfo.name}
              </span>
            )}
          </DialogTitle>
          <p className="text-xs text-slate-500">
            {formatTimestamp(contextModalMessage.timestamp)}
          </p>
        </DialogHeader>

        <ScrollArea className="flex-1 -mx-6 px-6">
          <div className="space-y-6 pb-4">
            {/* Full HTML Content */}
            {hasFullContent ? (
              <div className="prose prose-invert prose-sm max-w-none">
                <div
                  className="bg-slate-800/30 rounded-lg p-4 border border-white/5"
                  dangerouslySetInnerHTML={{
                    __html: contextModalMessage.fullContent!,
                  }}
                />
              </div>
            ) : (
              <div className="text-sm text-slate-300 whitespace-pre-wrap bg-slate-800/30 rounded-lg p-4 border border-white/5">
                {contextModalMessage.content}
              </div>
            )}

            {/* Event Details */}
            {contextModalMessage.type === "event" &&
              contextModalMessage.eventDetails && (
                <div className="space-y-4">
                  <h3 className="text-sm font-medium text-white">
                    Event Details
                  </h3>
                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div className="space-y-1">
                      <span className="text-slate-500">Time</span>
                      <p className="text-white">
                        {contextModalMessage.eventDetails.startTime} -{" "}
                        {contextModalMessage.eventDetails.endTime}
                      </p>
                    </div>
                    {contextModalMessage.eventDetails.location && (
                      <div className="space-y-1">
                        <span className="text-slate-500">Location</span>
                        <p className="text-white">
                          {contextModalMessage.eventDetails.location}
                        </p>
                      </div>
                    )}
                    <div className="space-y-1 col-span-2">
                      <span className="text-slate-500">Attendees</span>
                      <div className="flex flex-wrap gap-2">
                        {contextModalMessage.eventDetails.attendees.map(
                          (attendee) => (
                            <span
                              key={attendee.id}
                              className="px-2 py-1 bg-slate-800 rounded-full text-xs text-slate-300"
                            >
                              {attendee.name}
                            </span>
                          )
                        )}
                      </div>
                    </div>
                    {contextModalMessage.eventDetails.description && (
                      <div className="space-y-1 col-span-2">
                        <span className="text-slate-500">Description</span>
                        <p className="text-slate-300">
                          {contextModalMessage.eventDetails.description}
                        </p>
                      </div>
                    )}
                  </div>
                </div>
              )}

            {/* Social Content Details */}
            {contextModalMessage.type === "social" &&
              contextModalMessage.socialContent && (
                <div className="space-y-4">
                  <h3 className="text-sm font-medium text-white">
                    Content Details
                  </h3>
                  <div className="space-y-2 text-sm">
                    <p className="text-slate-300">
                      Platform:{" "}
                      <span className="text-white capitalize">
                        {contextModalMessage.socialContent.platform}
                      </span>
                    </p>
                    <p className="text-slate-300">
                      Author:{" "}
                      <span className="text-white">
                        {contextModalMessage.socialContent.author}
                      </span>
                    </p>
                    {contextModalMessage.socialContent.description && (
                      <p className="text-slate-400 mt-2">
                        {contextModalMessage.socialContent.description}
                      </p>
                    )}
                  </div>
                </div>
              )}

            {/* Attachments */}
            {hasAttachments && (
              <div className="space-y-3">
                <h3 className="text-sm font-medium text-white flex items-center gap-2">
                  <Paperclip size={14} />
                  Attachments ({contextModalMessage.attachments!.length})
                </h3>
                <div className="space-y-2">
                  {contextModalMessage.attachments!.map((attachment) => (
                    <div
                      key={attachment.id}
                      className="flex items-center gap-3 p-3 rounded-lg bg-slate-800/50 border border-white/5 hover:bg-slate-800/70 transition-colors group cursor-pointer"
                    >
                      <div className="h-10 w-10 rounded-lg bg-slate-700/50 flex items-center justify-center">
                        {getFileIcon(attachment.type)}
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="text-sm text-white truncate">
                          {attachment.name}
                        </div>
                        <div className="text-xs text-slate-500">
                          {formatFileSize(attachment.size)}
                        </div>
                      </div>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="opacity-0 group-hover:opacity-100 transition-opacity"
                      >
                        <Download size={16} className="text-slate-400" />
                      </Button>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </ScrollArea>

        {/* Footer Actions */}
        <div className="flex items-center justify-between pt-4 border-t border-white/5">
          <Button variant="ghost" size="sm" onClick={closeContextModal}>
            Close
          </Button>
          <Button variant="outline" size="sm" className="gap-2">
            <ExternalLink size={14} />
            Open in Source
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
