"use client";

import { SocialContent } from "@/types";
import { Eye, Heart, MessageCircle, Share2, ExternalLink } from "lucide-react";
import { formatNumber } from "@/lib/utils";
import { Button } from "@/components/ui/button";

interface ContentCardProps {
  content: SocialContent;
}

export function ContentCard({ content }: ContentCardProps) {
  const handleOpen = () => {
    window.open(content.url, "_blank");
  };

  const platformColors = {
    youtube: "text-red-400",
    linkedin: "text-sky-400",
    twitter: "text-slate-300",
  };

  return (
    <div className="glass-panel overflow-hidden max-w-md animate-fade-in">
      {/* Thumbnail */}
      {content.thumbnail && (
        <div className="relative aspect-video bg-slate-800">
          <img
            src={content.thumbnail}
            alt={content.title || "Content thumbnail"}
            className="w-full h-full object-cover"
          />
          {content.platform === "youtube" && (
            <div className="absolute inset-0 flex items-center justify-center">
              <div className="h-14 w-14 rounded-full bg-red-600/90 flex items-center justify-center">
                <div className="w-0 h-0 border-t-8 border-b-8 border-l-12 border-transparent border-l-white ml-1" />
              </div>
            </div>
          )}
        </div>
      )}

      {/* Content */}
      <div className="p-4 space-y-3">
        {/* Author */}
        <div className="flex items-center gap-2">
          <div className="h-8 w-8 rounded-full bg-slate-700 flex items-center justify-center">
            <span className="text-sm text-slate-300">
              {content.author.charAt(0)}
            </span>
          </div>
          <span className={`text-sm font-medium ${platformColors[content.platform]}`}>
            {content.author}
          </span>
        </div>

        {/* Title */}
        {content.title && (
          <h3 className="font-medium text-white line-clamp-2">{content.title}</h3>
        )}

        {/* Description */}
        {content.description && (
          <p className="text-sm text-slate-400 line-clamp-3">
            {content.description}
          </p>
        )}

        {/* Stats */}
        <div className="flex items-center gap-4 text-xs text-slate-500">
          {content.stats.views !== undefined && (
            <div className="flex items-center gap-1">
              <Eye size={12} />
              <span>{formatNumber(content.stats.views)}</span>
            </div>
          )}
          {content.stats.likes !== undefined && (
            <div className="flex items-center gap-1">
              <Heart size={12} />
              <span>{formatNumber(content.stats.likes)}</span>
            </div>
          )}
          {content.stats.comments !== undefined && (
            <div className="flex items-center gap-1">
              <MessageCircle size={12} />
              <span>{formatNumber(content.stats.comments)}</span>
            </div>
          )}
          {content.stats.shares !== undefined && (
            <div className="flex items-center gap-1">
              <Share2 size={12} />
              <span>{formatNumber(content.stats.shares)}</span>
            </div>
          )}
        </div>

        {/* Open Button */}
        <Button
          onClick={handleOpen}
          variant="outline"
          size="sm"
          className="w-full gap-2"
        >
          <ExternalLink size={14} />
          Open in {content.platform === "youtube" ? "YouTube" : content.platform === "linkedin" ? "LinkedIn" : "X"}
        </Button>
      </div>
    </div>
  );
}
