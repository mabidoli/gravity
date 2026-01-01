export type SourceType =
  | "email"
  | "whatsapp"
  | "slack"
  | "teams"
  | "calendar"
  | "task"
  | "youtube"
  | "linkedin"
  | "twitter";

export interface User {
  id: string;
  name: string;
  email?: string;
  avatar?: string;
}

export interface Attachment {
  id: string;
  name: string;
  type: string;
  size: number;
  url: string;
}

export interface CalendarEvent {
  id: string;
  title: string;
  startTime: string;
  endTime: string;
  attendees: User[];
  location?: string;
  meetingLink?: string;
  description?: string;
}

export interface SocialContent {
  id: string;
  platform: "youtube" | "linkedin" | "twitter";
  author: string;
  authorAvatar?: string;
  thumbnail?: string;
  title?: string;
  description?: string;
  stats: {
    views?: number;
    likes?: number;
    comments?: number;
    shares?: number;
  };
  url: string;
}

export interface AIInsight {
  id: string;
  type: "suggestion" | "analysis" | "draft";
  label: string;
  content: string;
  isDraft: boolean;
}

export interface Message {
  id: string;
  sender: "user" | "other" | "system";
  senderInfo?: User;
  content: string;
  timestamp: string;
  type: "text" | "event" | "social";
  eventDetails?: CalendarEvent;
  socialContent?: SocialContent;
  aiInsights?: AIInsight[];
  attachments?: Attachment[];
  fullContent?: string;
}

export interface PriorityItem {
  id: string;
  title: string;
  source: SourceType;
  priority: "high" | "medium" | "low";
  unread: boolean;
  snippet: string;
  timestamp: string;
  participants: User[];
  messages: Message[];
}

export type FilterType = "all" | "high" | "unread";
