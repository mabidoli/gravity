"use client";

import { motion } from "framer-motion";
import { PriorityItem } from "@/types";
import { SourceIcon } from "./SourceIcon";
import { formatTimestamp } from "@/lib/utils";
import { cn } from "@/lib/utils";

interface PriorityItemCardProps {
  item: PriorityItem;
  isSelected: boolean;
  onClick: () => void;
}

export function PriorityItemCard({
  item,
  isSelected,
  onClick,
}: PriorityItemCardProps) {
  const priorityColors = {
    high: "bg-red-500",
    medium: "bg-amber-500",
    low: "bg-green-500",
  };

  return (
    <motion.div
      layoutId={`item-${item.id}`}
      onClick={onClick}
      className={cn(
        "relative p-4 rounded-xl cursor-pointer transition-all duration-200",
        "border border-white/5 hover:border-white/10",
        "bg-slate-900/40 hover:bg-slate-800/50",
        isSelected && "glass-active border-teal-500/30 bg-teal-500/5"
      )}
      whileHover={{ scale: 1.01 }}
      whileTap={{ scale: 0.99 }}
    >
      {/* Priority Indicator */}
      <div
        className={cn(
          "absolute left-0 top-1/2 -translate-y-1/2 w-1 h-8 rounded-r-full",
          priorityColors[item.priority]
        )}
      />

      {/* Unread Indicator */}
      {item.unread && (
        <div className="absolute right-3 top-3">
          <div className="h-2 w-2 rounded-full bg-teal-400 animate-pulse" />
        </div>
      )}

      <div className="flex gap-3 pl-2">
        {/* Source Icon */}
        <SourceIcon source={item.source} />

        {/* Content */}
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between gap-2">
            <h3
              className={cn(
                "text-sm font-medium truncate",
                item.unread ? "text-white" : "text-slate-300"
              )}
            >
              {item.title}
            </h3>
            <span className="text-xs text-slate-500 whitespace-nowrap">
              {formatTimestamp(item.timestamp)}
            </span>
          </div>

          <p className="text-xs text-slate-400 mt-1 line-clamp-2">
            {item.snippet}
          </p>

          {/* Participants */}
          {item.participants.length > 0 && (
            <div className="flex items-center gap-1 mt-2">
              <div className="flex -space-x-2">
                {item.participants.slice(0, 3).map((participant, idx) => (
                  <div
                    key={participant.id}
                    className="h-5 w-5 rounded-full bg-slate-700 border border-slate-600 flex items-center justify-center"
                    title={participant.name}
                  >
                    <span className="text-[10px] text-slate-300">
                      {participant.name.charAt(0)}
                    </span>
                  </div>
                ))}
              </div>
              {item.participants.length > 3 && (
                <span className="text-xs text-slate-500">
                  +{item.participants.length - 3}
                </span>
              )}
            </div>
          )}
        </div>
      </div>
    </motion.div>
  );
}
