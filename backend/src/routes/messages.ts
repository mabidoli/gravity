import { Router, Request, Response } from "express";
import { Message } from "../types/index.js";
import { mockStreamItems } from "../data/mockData.js";

export const messagesRouter = Router();

// Send a message to a stream item
messagesRouter.post("/:itemId", (req: Request, res: Response) => {
  const { itemId } = req.params;
  const { content } = req.body;

  const itemIndex = mockStreamItems.findIndex((item) => item.id === itemId);

  if (itemIndex === -1) {
    res.status(404).json({ error: "Item not found" });
    return;
  }

  const item = mockStreamItems[itemIndex];

  // Create user message
  const userMessage: Message = {
    id: `msg-${Date.now()}`,
    sender: "user",
    content,
    timestamp: new Date().toISOString(),
    type: "text",
  };

  item.messages.push(userMessage);

  // Handle system interactions (Calendar, AI responses)
  const lowerContent = content.toLowerCase();

  // Calendar rescheduling
  if (item.source === "calendar" && lowerContent.includes("reschedule")) {
    const timeMatch = content.match(/(\d{1,2})\s*(pm|am|:00)/i);
    if (timeMatch) {
      const newTime = timeMatch[1] + (timeMatch[2].toLowerCase() === "pm" ? ":00 PM" : ":00 AM");

      // Find the event message and update it
      const eventMessage = item.messages.find((m) => m.type === "event" && m.eventDetails);
      if (eventMessage?.eventDetails) {
        eventMessage.eventDetails.startTime = newTime;
        eventMessage.eventDetails.endTime = `${parseInt(timeMatch[1]) + 1}:00 ${timeMatch[2].toUpperCase()}`;
      }

      const systemMessage: Message = {
        id: `msg-${Date.now() + 1}`,
        sender: "system",
        content: `✓ I've updated the meeting time to ${newTime}. Notifications have been sent to all attendees.`,
        timestamp: new Date().toISOString(),
        type: "text",
      };
      item.messages.push(systemMessage);
    }
  }

  // Social content summarization
  if (
    (item.source === "youtube" || item.source === "linkedin" || item.source === "twitter") &&
    lowerContent.includes("summarize")
  ) {
    const summaryMessage: Message = {
      id: `msg-${Date.now() + 1}`,
      sender: "system",
      content: generateSummary(item.source),
      timestamp: new Date().toISOString(),
      type: "text",
      aiInsights: [
        {
          id: `insight-${Date.now()}`,
          type: "analysis",
          label: "📊 Key Points",
          content: generateDetailedSummary(item.source),
          isDraft: false,
        },
      ],
    };
    item.messages.push(summaryMessage);
  }

  res.json(item);
});

// Generate AI draft for a message
messagesRouter.post("/:itemId/draft", (req: Request, res: Response) => {
  const { itemId } = req.params;
  const { messageId, refinement } = req.body;

  const item = mockStreamItems.find((item) => item.id === itemId);

  if (!item) {
    res.status(404).json({ error: "Item not found" });
    return;
  }

  // Mock AI draft generation/refinement
  let draftContent = "Thank you for reaching out. I'll review this and get back to you shortly.";

  if (refinement) {
    if (refinement.toLowerCase().includes("formal") || refinement.toLowerCase().includes("professional")) {
      draftContent =
        "Dear colleague,\n\nThank you for bringing this to my attention. I will thoroughly review the matter and provide you with a comprehensive response at my earliest convenience.\n\nBest regards";
    } else if (refinement.toLowerCase().includes("friendly") || refinement.toLowerCase().includes("casual")) {
      draftContent = "Hey! Thanks for the heads up 👍 I'll take a look and get back to you soon!";
    } else if (refinement.toLowerCase().includes("short") || refinement.toLowerCase().includes("brief")) {
      draftContent = "Got it, will review. Thanks!";
    }
  }

  res.json({
    draft: draftContent,
    timestamp: new Date().toISOString(),
  });
});

function generateSummary(source: string): string {
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

function generateDetailedSummary(source: string): string {
  return `Based on my analysis, here are the most important points:\n\n1. The main argument centers around emerging technology trends\n2. Several actionable insights are provided\n3. The conclusion suggests further reading on related topics\n\nWould you like me to extract any specific information?`;
}
