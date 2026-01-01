"use client";

import { motion, AnimatePresence } from "framer-motion";
import { StreamSidebar } from "./StreamSidebar";
import { ChatInterface } from "./ChatInterface";
import { FullContextModal } from "./FullContextModal";
import { useGravityStore } from "@/store/useGravityStore";
import { cn } from "@/lib/utils";

export function MainLayout() {
  const { isMobileViewingChat, selectedItemId } = useGravityStore();

  return (
    <div className="h-screen w-full flex overflow-hidden">
      {/* Sidebar - Stream List */}
      <div
        className={cn(
          "w-full lg:w-[350px] xl:w-[400px] shrink-0 h-full",
          "lg:block",
          isMobileViewingChat && "hidden"
        )}
      >
        <StreamSidebar />
      </div>

      {/* Main Content - Chat Interface */}
      <div
        className={cn(
          "flex-1 h-full border-l border-white/5",
          "lg:block",
          !isMobileViewingChat && !selectedItemId && "hidden lg:block"
        )}
      >
        <AnimatePresence mode="wait">
          <motion.div
            key={selectedItemId || "empty"}
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -20 }}
            transition={{ duration: 0.2 }}
            className="h-full"
          >
            <ChatInterface />
          </motion.div>
        </AnimatePresence>
      </div>

      {/* Full Context Modal */}
      <FullContextModal />
    </div>
  );
}
