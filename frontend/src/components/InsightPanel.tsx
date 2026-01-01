"use client";

import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { AIInsight } from "@/types";
import { Button } from "@/components/ui/button";
import {
  Send,
  Wand2,
  Copy,
  RefreshCw,
  Check,
  ChevronDown,
  ChevronUp,
} from "lucide-react";
import { useGravityStore } from "@/store/useGravityStore";

interface InsightPanelProps {
  insight: AIInsight;
  itemId: string;
  messageId: string;
}

export function InsightPanel({ insight, itemId, messageId }: InsightPanelProps) {
  const [isExpanded, setIsExpanded] = useState(false);
  const [isRefining, setIsRefining] = useState(false);
  const [refinementInput, setRefinementInput] = useState("");
  const [copied, setCopied] = useState(false);

  const { refineDraft, regenerateDraft } = useGravityStore();

  const handleCopy = async () => {
    await navigator.clipboard.writeText(insight.content);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleRefine = () => {
    if (refinementInput.trim()) {
      refineDraft(itemId, messageId, insight.id, refinementInput);
      setRefinementInput("");
      setIsRefining(false);
    }
  };

  const handleRegenerate = () => {
    regenerateDraft(itemId, messageId, insight.id);
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleRefine();
    }
  };

  return (
    <div className="mt-2">
      {/* Smart Pill Trigger */}
      <button
        onClick={() => setIsExpanded(!isExpanded)}
        className="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded-full transition-all duration-200 bg-teal-500/10 text-teal-400 border border-teal-500/20 hover:bg-teal-500/20 hover:border-teal-500/30"
      >
        <span>{insight.label}</span>
        {isExpanded ? <ChevronUp size={12} /> : <ChevronDown size={12} />}
      </button>

      {/* Expanded Panel */}
      <AnimatePresence>
        {isExpanded && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: "auto" }}
            exit={{ opacity: 0, height: 0 }}
            transition={{ duration: 0.2 }}
            className="overflow-hidden"
          >
            <div className="mt-3 p-4 rounded-xl bg-slate-800/50 border border-white/5 space-y-3">
              {/* Insight Content */}
              <div className="text-sm text-slate-300 whitespace-pre-wrap">
                {insight.content}
              </div>

              {/* Refinement Input */}
              <AnimatePresence>
                {isRefining && (
                  <motion.div
                    initial={{ opacity: 0, y: -10 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -10 }}
                    className="flex gap-2"
                  >
                    <input
                      type="text"
                      value={refinementInput}
                      onChange={(e) => setRefinementInput(e.target.value)}
                      onKeyDown={handleKeyDown}
                      placeholder="Make it more formal, shorter, friendlier..."
                      className="flex-1 px-3 py-2 text-sm bg-slate-900/50 border border-white/10 rounded-lg text-white placeholder:text-slate-500 focus:outline-none focus:border-teal-500/50"
                      autoFocus
                    />
                    <Button size="sm" onClick={handleRefine}>
                      Apply
                    </Button>
                  </motion.div>
                )}
              </AnimatePresence>

              {/* Actions */}
              <div className="flex flex-wrap gap-2">
                {insight.isDraft && (
                  <Button size="sm" className="gap-1.5">
                    <Send size={14} />
                    Send
                  </Button>
                )}
                <Button
                  size="sm"
                  variant="secondary"
                  className="gap-1.5"
                  onClick={() => setIsRefining(!isRefining)}
                >
                  <Wand2 size={14} />
                  Refine
                </Button>
                <Button
                  size="sm"
                  variant="ghost"
                  className="gap-1.5"
                  onClick={handleCopy}
                >
                  {copied ? <Check size={14} /> : <Copy size={14} />}
                  {copied ? "Copied" : "Copy"}
                </Button>
                <Button
                  size="sm"
                  variant="ghost"
                  className="gap-1.5"
                  onClick={handleRegenerate}
                >
                  <RefreshCw size={14} />
                  Regenerate
                </Button>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
