import { create } from "zustand";
import { PriorityItem, Message, FilterType, AIInsight } from "@/types";
import { mockStreamItems } from "@/data/mockData";

interface GravityState {
  // Stream state
  items: PriorityItem[];
  selectedItemId: string | null;
  filter: FilterType;
  isMobileViewingChat: boolean;

  // Modal state
  isContextModalOpen: boolean;
  contextModalMessage: Message | null;

  // Actions
  setFilter: (filter: FilterType) => void;
  selectItem: (itemId: string | null) => void;
  markAsRead: (itemId: string) => void;
  setMobileViewingChat: (viewing: boolean) => void;

  // Modal actions
  openContextModal: (message: Message) => void;
  closeContextModal: () => void;

  // Message actions
  sendMessage: (itemId: string, content: string) => void;
  updateEventTime: (itemId: string, messageId: string, newStartTime: string, newEndTime: string) => void;
  refineDraft: (itemId: string, messageId: string, insightId: string, refinement: string) => void;
  regenerateDraft: (itemId: string, messageId: string, insightId: string) => void;

  // Computed
  getFilteredItems: () => PriorityItem[];
  getSelectedItem: () => PriorityItem | undefined;
  getUnreadCount: () => number;
}

export const useGravityStore = create<GravityState>((set, get) => ({
  // Initial state
  items: mockStreamItems,
  selectedItemId: null,
  filter: "all",
  isMobileViewingChat: false,
  isContextModalOpen: false,
  contextModalMessage: null,

  // Filter actions
  setFilter: (filter) => set({ filter }),

  // Selection actions
  selectItem: (itemId) => {
    set({ selectedItemId: itemId, isMobileViewingChat: itemId !== null });
    if (itemId) {
      get().markAsRead(itemId);
    }
  },

  markAsRead: (itemId) =>
    set((state) => ({
      items: state.items.map((item) =>
        item.id === itemId ? { ...item, unread: false } : item
      ),
    })),

  setMobileViewingChat: (viewing) => set({ isMobileViewingChat: viewing }),

  // Modal actions
  openContextModal: (message) =>
    set({ isContextModalOpen: true, contextModalMessage: message }),

  closeContextModal: () =>
    set({ isContextModalOpen: false, contextModalMessage: null }),

  // Message actions
  sendMessage: (itemId, content) => {
    const lowerContent = content.toLowerCase();

    set((state) => {
      const item = state.items.find((i) => i.id === itemId);
      if (!item) return state;

      const userMessage: Message = {
        id: `msg-${Date.now()}`,
        sender: "user",
        content,
        timestamp: new Date().toISOString(),
        type: "text",
      };

      const newMessages = [...item.messages, userMessage];

      // Handle calendar rescheduling
      if (item.source === "calendar" && lowerContent.includes("reschedule")) {
        const timeMatch = content.match(/(\d{1,2})\s*(pm|am|:00)?/i);
        if (timeMatch) {
          const hour = parseInt(timeMatch[1]);
          const isPM = timeMatch[2]?.toLowerCase() === "pm" || hour >= 1 && hour <= 6;
          const formattedHour = isPM && hour < 12 ? hour : hour;
          const newStartTime = `${formattedHour}:00 ${isPM ? "PM" : "AM"}`;
          const newEndTime = `${formattedHour + 1}:00 ${isPM ? "PM" : "AM"}`;

          // Update the event details
          const updatedMessages = newMessages.map((msg) => {
            if (msg.type === "event" && msg.eventDetails) {
              return {
                ...msg,
                eventDetails: {
                  ...msg.eventDetails,
                  startTime: newStartTime,
                  endTime: newEndTime,
                },
              };
            }
            return msg;
          });

          const systemMessage: Message = {
            id: `msg-${Date.now() + 1}`,
            sender: "system",
            content: `✓ I've updated the meeting time to ${newStartTime}. Notifications have been sent to all attendees.`,
            timestamp: new Date().toISOString(),
            type: "text",
          };

          updatedMessages.push(systemMessage);

          return {
            items: state.items.map((i) =>
              i.id === itemId ? { ...i, messages: updatedMessages } : i
            ),
          };
        }
      }

      // Handle social content summarization
      if (
        (item.source === "youtube" ||
          item.source === "linkedin" ||
          item.source === "twitter") &&
        lowerContent.includes("summarize")
      ) {
        const summaryContent = getSummaryContent(item.source);
        const summaryMessage: Message = {
          id: `msg-${Date.now() + 1}`,
          sender: "system",
          content: summaryContent,
          timestamp: new Date().toISOString(),
          type: "text",
          aiInsights: [
            {
              id: `insight-${Date.now()}`,
              type: "analysis",
              label: "📊 Key Points",
              content: getDetailedSummary(item.source),
              isDraft: false,
            },
          ],
        };
        newMessages.push(summaryMessage);
      }

      return {
        items: state.items.map((i) =>
          i.id === itemId ? { ...i, messages: newMessages } : i
        ),
      };
    });
  },

  updateEventTime: (itemId, messageId, newStartTime, newEndTime) =>
    set((state) => ({
      items: state.items.map((item) =>
        item.id === itemId
          ? {
              ...item,
              messages: item.messages.map((msg) =>
                msg.id === messageId && msg.eventDetails
                  ? {
                      ...msg,
                      eventDetails: {
                        ...msg.eventDetails,
                        startTime: newStartTime,
                        endTime: newEndTime,
                      },
                    }
                  : msg
              ),
            }
          : item
      ),
    })),

  refineDraft: (itemId, messageId, insightId, refinement) =>
    set((state) => ({
      items: state.items.map((item) =>
        item.id === itemId
          ? {
              ...item,
              messages: item.messages.map((msg) =>
                msg.id === messageId
                  ? {
                      ...msg,
                      aiInsights: msg.aiInsights?.map((insight) =>
                        insight.id === insightId
                          ? {
                              ...insight,
                              content: getRefinedDraft(insight.content, refinement),
                            }
                          : insight
                      ),
                    }
                  : msg
              ),
            }
          : item
      ),
    })),

  regenerateDraft: (itemId, messageId, insightId) =>
    set((state) => ({
      items: state.items.map((item) =>
        item.id === itemId
          ? {
              ...item,
              messages: item.messages.map((msg) =>
                msg.id === messageId
                  ? {
                      ...msg,
                      aiInsights: msg.aiInsights?.map((insight) =>
                        insight.id === insightId
                          ? {
                              ...insight,
                              content: getRegeneratedDraft(),
                            }
                          : insight
                      ),
                    }
                  : msg
              ),
            }
          : item
      ),
    })),

  // Computed values
  getFilteredItems: () => {
    const { items, filter } = get();
    let filteredItems = [...items];

    if (filter === "high") {
      filteredItems = filteredItems.filter((item) => item.priority === "high");
    } else if (filter === "unread") {
      filteredItems = filteredItems.filter((item) => item.unread);
    }

    // Sort by priority then timestamp
    return filteredItems.sort((a, b) => {
      const priorityOrder = { high: 0, medium: 1, low: 2 };
      if (priorityOrder[a.priority] !== priorityOrder[b.priority]) {
        return priorityOrder[a.priority] - priorityOrder[b.priority];
      }
      return new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime();
    });
  },

  getSelectedItem: () => {
    const { items, selectedItemId } = get();
    return items.find((item) => item.id === selectedItemId);
  },

  getUnreadCount: () => {
    return get().items.filter((item) => item.unread).length;
  },
}));

// Helper functions for mock AI responses
function getSummaryContent(source: string): string {
  const summaries: Record<string, string> = {
    youtube:
      "📺 **Video Summary**\n\nThis video discusses the latest trends in AI development, focusing on three key areas:\n\n• Large language models and their applications\n• Edge computing and on-device AI\n• Ethical considerations in AI deployment",
    linkedin:
      "💼 **Post Summary**\n\nThe author shares insights on professional growth:\n\n• Importance of continuous learning\n• Building meaningful connections\n• Embracing change in the workplace",
    twitter:
      "🐦 **Thread Summary**\n\nKey takeaways from this thread:\n\n• Market trends for Q4\n• Predictions for the upcoming quarter\n• Industry expert opinions",
  };
  return summaries[source] || "Summary not available for this content type.";
}

function getDetailedSummary(source: string): string {
  return `Based on my analysis, here are the most important points:\n\n1. The main argument centers around emerging technology trends\n2. Several actionable insights are provided\n3. The conclusion suggests further reading on related topics\n\nWould you like me to extract any specific information?`;
}

function getRefinedDraft(currentDraft: string, refinement: string): string {
  const lowerRefinement = refinement.toLowerCase();

  if (lowerRefinement.includes("formal") || lowerRefinement.includes("professional")) {
    return "Dear colleague,\n\nThank you for bringing this to my attention. I will thoroughly review the matter and provide you with a comprehensive response at my earliest convenience.\n\nBest regards";
  }

  if (lowerRefinement.includes("friendly") || lowerRefinement.includes("casual")) {
    return "Hey! Thanks for the heads up 👍 I'll take a look and get back to you soon!";
  }

  if (lowerRefinement.includes("short") || lowerRefinement.includes("brief")) {
    return "Got it, will review. Thanks!";
  }

  if (lowerRefinement.includes("polite") || lowerRefinement.includes("4 pm") || lowerRefinement.includes("4pm")) {
    return "I'm afraid I can't make it at the original time. Would 4 PM work for you instead? I'd really appreciate your flexibility.";
  }

  return `${currentDraft}\n\n[Updated based on your feedback: "${refinement}"]`;
}

function getRegeneratedDraft(): string {
  const drafts = [
    "Thank you for your message. I've reviewed the details and wanted to follow up with my thoughts. Let me know if you'd like to discuss further.",
    "I appreciate you reaching out. After careful consideration, I'd like to propose we schedule a quick call to align on next steps.",
    "Thanks for sharing this with me. I have a few questions before proceeding - would you have time for a brief sync this week?",
  ];
  return drafts[Math.floor(Math.random() * drafts.length)];
}
