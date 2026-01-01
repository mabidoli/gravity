"use client";

import { motion, AnimatePresence } from "framer-motion";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ScrollArea } from "@/components/ui/scroll-area";
import { PriorityItemCard } from "./PriorityItemCard";
import { useGravityStore } from "@/store/useGravityStore";
import { FilterType } from "@/types";
import { Inbox, AlertCircle, Bell } from "lucide-react";

export function StreamSidebar() {
  const {
    filter,
    setFilter,
    selectedItemId,
    selectItem,
    getFilteredItems,
    getUnreadCount,
  } = useGravityStore();

  const items = getFilteredItems();
  const unreadCount = getUnreadCount();

  const filterOptions: { value: FilterType; label: string; icon: React.ReactNode }[] = [
    { value: "all", label: "All", icon: <Inbox size={14} /> },
    { value: "high", label: "High Priority", icon: <AlertCircle size={14} /> },
    { value: "unread", label: "Unread", icon: <Bell size={14} /> },
  ];

  return (
    <div className="h-full flex flex-col bg-slate-900/30 border-r border-white/5">
      {/* Header */}
      <div className="p-4 border-b border-white/5">
        <div className="flex items-center justify-between mb-4">
          <h1 className="text-xl font-bold text-gradient">Gravity</h1>
          {unreadCount > 0 && (
            <span className="px-2 py-0.5 text-xs font-medium bg-teal-500/20 text-teal-400 rounded-full">
              {unreadCount} new
            </span>
          )}
        </div>

        {/* Filter Tabs */}
        <Tabs value={filter} onValueChange={(v) => setFilter(v as FilterType)}>
          <TabsList className="w-full grid grid-cols-3">
            {filterOptions.map((option) => (
              <TabsTrigger
                key={option.value}
                value={option.value}
                className="flex items-center gap-1.5 text-xs"
              >
                {option.icon}
                <span className="hidden sm:inline">{option.label}</span>
              </TabsTrigger>
            ))}
          </TabsList>
        </Tabs>
      </div>

      {/* Stream List */}
      <ScrollArea className="flex-1">
        <div className="p-3 space-y-2">
          <AnimatePresence mode="popLayout">
            {items.map((item) => (
              <motion.div
                key={item.id}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, scale: 0.95 }}
                transition={{ duration: 0.2 }}
              >
                <PriorityItemCard
                  item={item}
                  isSelected={selectedItemId === item.id}
                  onClick={() => selectItem(item.id)}
                />
              </motion.div>
            ))}
          </AnimatePresence>

          {items.length === 0 && (
            <div className="flex flex-col items-center justify-center py-12 text-slate-500">
              <Inbox size={32} className="mb-2 opacity-50" />
              <p className="text-sm">No items to show</p>
            </div>
          )}
        </div>
      </ScrollArea>
    </div>
  );
}
