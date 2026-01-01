"use client";

import {
  Mail,
  MessageCircle,
  Hash,
  Users,
  Calendar,
  CheckSquare,
  Youtube,
  Linkedin,
  Twitter,
  LucideIcon,
} from "lucide-react";
import { SourceType } from "@/types";
import { cn } from "@/lib/utils";

const sourceConfig: Record<
  SourceType,
  { icon: LucideIcon; color: string; bgColor: string }
> = {
  email: {
    icon: Mail,
    color: "text-blue-400",
    bgColor: "bg-blue-500/10",
  },
  whatsapp: {
    icon: MessageCircle,
    color: "text-green-400",
    bgColor: "bg-green-500/10",
  },
  slack: {
    icon: Hash,
    color: "text-purple-400",
    bgColor: "bg-purple-500/10",
  },
  teams: {
    icon: Users,
    color: "text-indigo-400",
    bgColor: "bg-indigo-500/10",
  },
  calendar: {
    icon: Calendar,
    color: "text-teal-400",
    bgColor: "bg-teal-500/10",
  },
  task: {
    icon: CheckSquare,
    color: "text-orange-400",
    bgColor: "bg-orange-500/10",
  },
  youtube: {
    icon: Youtube,
    color: "text-red-400",
    bgColor: "bg-red-500/10",
  },
  linkedin: {
    icon: Linkedin,
    color: "text-sky-400",
    bgColor: "bg-sky-500/10",
  },
  twitter: {
    icon: Twitter,
    color: "text-slate-300",
    bgColor: "bg-slate-500/10",
  },
};

interface SourceIconProps {
  source: SourceType;
  size?: "sm" | "md" | "lg";
  className?: string;
}

export function SourceIcon({ source, size = "md", className }: SourceIconProps) {
  const config = sourceConfig[source];
  const Icon = config.icon;

  const sizeClasses = {
    sm: "h-6 w-6",
    md: "h-8 w-8",
    lg: "h-10 w-10",
  };

  const iconSizes = {
    sm: 12,
    md: 16,
    lg: 20,
  };

  return (
    <div
      className={cn(
        "flex items-center justify-center rounded-lg transition-colors duration-200",
        sizeClasses[size],
        config.bgColor,
        className
      )}
    >
      <Icon size={iconSizes[size]} className={config.color} />
    </div>
  );
}

export function getSourceLabel(source: SourceType): string {
  const labels: Record<SourceType, string> = {
    email: "Email",
    whatsapp: "WhatsApp",
    slack: "Slack",
    teams: "Teams",
    calendar: "Calendar",
    task: "Task",
    youtube: "YouTube",
    linkedin: "LinkedIn",
    twitter: "X",
  };
  return labels[source];
}
