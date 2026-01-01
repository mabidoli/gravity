import { PriorityItem } from "../types/index.js";

export const mockStreamItems: PriorityItem[] = [
  // High Priority Email
  {
    id: "item-1",
    title: "Q4 Revenue Report - Action Required",
    source: "email",
    priority: "high",
    unread: true,
    snippet: "Please review the attached Q4 projections before tomorrow's board meeting...",
    timestamp: new Date(Date.now() - 30 * 60 * 1000).toISOString(),
    participants: [
      { id: "user-1", name: "Sarah Chen", email: "sarah.chen@company.com", avatar: "/avatars/sarah.jpg" },
    ],
    messages: [
      {
        id: "msg-1-1",
        sender: "other",
        senderInfo: { id: "user-1", name: "Sarah Chen", email: "sarah.chen@company.com" },
        content:
          "Hi,\n\nPlease review the attached Q4 projections before tomorrow's board meeting. We need your sign-off on the marketing budget allocation.\n\nThe key highlights:\n• Revenue up 23% YoY\n• Customer acquisition cost down 15%\n• New market expansion on track\n\nLet me know if you have any questions.\n\nBest,\nSarah",
        timestamp: new Date(Date.now() - 30 * 60 * 1000).toISOString(),
        type: "text",
        attachments: [
          { id: "att-1", name: "Q4_Revenue_Report.pdf", type: "application/pdf", size: 2456000, url: "/files/q4-report.pdf" },
          { id: "att-2", name: "Budget_Allocation.xlsx", type: "application/xlsx", size: 156000, url: "/files/budget.xlsx" },
        ],
        aiInsights: [
          {
            id: "insight-1",
            type: "draft",
            label: "✨ Draft Available",
            content:
              "Hi Sarah,\n\nThank you for the comprehensive report. I've reviewed the Q4 projections and the numbers look solid. I approve the marketing budget allocation as proposed.\n\nA few notes:\n• Great work on reducing CAC\n• Let's discuss the new market expansion timeline in tomorrow's meeting\n\nBest regards",
            isDraft: true,
          },
        ],
        fullContent:
          '<!DOCTYPE html><html><body style="font-family: Arial, sans-serif;"><p>Hi,</p><p>Please review the attached Q4 projections before tomorrow\'s board meeting. We need your sign-off on the marketing budget allocation.</p><p><strong>The key highlights:</strong></p><ul><li>Revenue up 23% YoY</li><li>Customer acquisition cost down 15%</li><li>New market expansion on track</li></ul><p>Let me know if you have any questions.</p><p>Best,<br>Sarah</p></body></html>',
      },
    ],
  },

  // High Priority Calendar
  {
    id: "item-2",
    title: "Product Sync - Starting in 15 minutes",
    source: "calendar",
    priority: "high",
    unread: true,
    snippet: "Weekly product team sync with engineering leads",
    timestamp: new Date(Date.now() + 15 * 60 * 1000).toISOString(),
    participants: [
      { id: "user-2", name: "Mike Johnson", avatar: "/avatars/mike.jpg" },
      { id: "user-3", name: "Emily Davis", avatar: "/avatars/emily.jpg" },
      { id: "user-4", name: "Alex Kim", avatar: "/avatars/alex.jpg" },
    ],
    messages: [
      {
        id: "msg-2-1",
        sender: "system",
        content: "Upcoming meeting",
        timestamp: new Date(Date.now() + 15 * 60 * 1000).toISOString(),
        type: "event",
        eventDetails: {
          id: "event-1",
          title: "Product Sync",
          startTime: "2:00 PM",
          endTime: "3:00 PM",
          attendees: [
            { id: "user-2", name: "Mike Johnson" },
            { id: "user-3", name: "Emily Davis" },
            { id: "user-4", name: "Alex Kim" },
          ],
          location: "Conference Room A / Zoom",
          meetingLink: "https://zoom.us/j/123456789",
          description: "Weekly sync to discuss product roadmap, current sprint progress, and blockers.",
        },
      },
    ],
  },

  // Slack Message
  {
    id: "item-3",
    title: "#engineering - Deployment Issue",
    source: "slack",
    priority: "high",
    unread: true,
    snippet: "@channel Production deployment failed. Rolling back...",
    timestamp: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
    participants: [{ id: "user-5", name: "DevOps Bot" }],
    messages: [
      {
        id: "msg-3-1",
        sender: "other",
        senderInfo: { id: "user-5", name: "DevOps Bot" },
        content:
          "🚨 @channel Production deployment failed.\n\nError: Database migration timeout\nBuild: #4521\nCommit: abc123f\n\nRolling back to previous version...",
        timestamp: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
        type: "text",
        aiInsights: [
          {
            id: "insight-2",
            type: "analysis",
            label: "⚠️ Risk Detected",
            content:
              "This deployment failure appears to be related to the new user authentication migration. The timeout suggests the migration script may need optimization for the production dataset size.\n\nRecommended actions:\n1. Check migration logs in CloudWatch\n2. Review the authentication table indexes\n3. Consider batched migration approach",
            isDraft: false,
          },
        ],
      },
      {
        id: "msg-3-2",
        sender: "other",
        senderInfo: { id: "user-6", name: "John Developer" },
        content: "On it! Checking the logs now.",
        timestamp: new Date(Date.now() - 3 * 60 * 1000).toISOString(),
        type: "text",
      },
    ],
  },

  // WhatsApp
  {
    id: "item-4",
    title: "Family Group",
    source: "whatsapp",
    priority: "low",
    unread: false,
    snippet: "Mom: Don't forget dinner on Sunday!",
    timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
    participants: [
      { id: "user-7", name: "Mom" },
      { id: "user-8", name: "Dad" },
      { id: "user-9", name: "Sister" },
    ],
    messages: [
      {
        id: "msg-4-1",
        sender: "other",
        senderInfo: { id: "user-7", name: "Mom" },
        content: "Don't forget dinner on Sunday! We're having your favorite 🍝",
        timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
        type: "text",
      },
      {
        id: "msg-4-2",
        sender: "other",
        senderInfo: { id: "user-9", name: "Sister" },
        content: "I'll bring dessert!",
        timestamp: new Date(Date.now() - 1.5 * 60 * 60 * 1000).toISOString(),
        type: "text",
      },
    ],
  },

  // YouTube Content
  {
    id: "item-5",
    title: "New AI Trends 2025 - Must Watch",
    source: "youtube",
    priority: "medium",
    unread: true,
    snippet: "The Future of AI: What's Coming Next",
    timestamp: new Date(Date.now() - 4 * 60 * 60 * 1000).toISOString(),
    participants: [{ id: "channel-1", name: "TechInsider" }],
    messages: [
      {
        id: "msg-5-1",
        sender: "system",
        content: "New video from a channel you follow",
        timestamp: new Date(Date.now() - 4 * 60 * 60 * 1000).toISOString(),
        type: "social",
        socialContent: {
          id: "yt-1",
          platform: "youtube",
          author: "TechInsider",
          authorAvatar: "/avatars/techinsider.jpg",
          thumbnail: "https://picsum.photos/seed/aitrends/640/360",
          title: "The Future of AI: What's Coming in 2025",
          description:
            "In this video, we explore the emerging trends in artificial intelligence, from large language models to autonomous systems. What does the future hold?",
          stats: {
            views: 245000,
            likes: 12400,
            comments: 892,
          },
          url: "https://youtube.com/watch?v=example123",
        },
        aiInsights: [
          {
            id: "insight-3",
            type: "suggestion",
            label: "💡 Quick Summary Available",
            content:
              "This 18-minute video covers:\n\n1. **LLM Evolution** - How models are becoming more efficient\n2. **Edge AI** - On-device processing trends\n3. **AI Regulation** - Upcoming policy changes\n4. **Industry Impact** - Which sectors will transform first\n\nKey takeaway: The presenter predicts a shift toward specialized AI agents by mid-2025.",
            isDraft: false,
          },
        ],
      },
    ],
  },

  // LinkedIn
  {
    id: "item-6",
    title: "Connection Request - Jane Smith",
    source: "linkedin",
    priority: "medium",
    unread: true,
    snippet: "VP of Engineering at TechCorp wants to connect",
    timestamp: new Date(Date.now() - 6 * 60 * 60 * 1000).toISOString(),
    participants: [{ id: "li-1", name: "Jane Smith", email: "jane.smith@techcorp.com" }],
    messages: [
      {
        id: "msg-6-1",
        sender: "system",
        content: "New connection request",
        timestamp: new Date(Date.now() - 6 * 60 * 60 * 1000).toISOString(),
        type: "social",
        socialContent: {
          id: "li-post-1",
          platform: "linkedin",
          author: "Jane Smith",
          authorAvatar: "/avatars/jane.jpg",
          title: "VP of Engineering at TechCorp",
          description:
            "Hi! I came across your profile and was impressed by your work on distributed systems. I'd love to connect and discuss potential collaboration opportunities.",
          stats: {
            views: 0,
          },
          url: "https://linkedin.com/in/janesmith",
        },
        aiInsights: [
          {
            id: "insight-4",
            type: "draft",
            label: "✨ Draft Response",
            content:
              "Hi Jane,\n\nThank you for reaching out! I'd be happy to connect. Your work at TechCorp on cloud infrastructure looks fascinating.\n\nI'd love to learn more about your team's approach to distributed systems. Would you be open to a brief call next week?\n\nBest regards",
            isDraft: true,
          },
        ],
      },
    ],
  },

  // Twitter/X
  {
    id: "item-7",
    title: "Trending in Tech",
    source: "twitter",
    priority: "low",
    unread: false,
    snippet: "Thread about startup funding trends going viral",
    timestamp: new Date(Date.now() - 8 * 60 * 60 * 1000).toISOString(),
    participants: [{ id: "tw-1", name: "@venturecap" }],
    messages: [
      {
        id: "msg-7-1",
        sender: "system",
        content: "Trending thread from someone you follow",
        timestamp: new Date(Date.now() - 8 * 60 * 60 * 1000).toISOString(),
        type: "social",
        socialContent: {
          id: "tw-post-1",
          platform: "twitter",
          author: "@venturecap",
          authorAvatar: "/avatars/venturecap.jpg",
          title: "The state of startup funding in 2025",
          description:
            "🧵 Thread: After analyzing 500+ Series A rounds this year, here's what I've learned about the current funding landscape...\n\n1/ Valuations are finally normalizing after the 2021 peak...",
          stats: {
            likes: 4500,
            comments: 234,
            shares: 1200,
          },
          url: "https://twitter.com/venturecap/status/123456789",
        },
      },
    ],
  },

  // Teams Message
  {
    id: "item-8",
    title: "Design Review - New Dashboard",
    source: "teams",
    priority: "medium",
    unread: true,
    snippet: "The new dashboard mockups are ready for review",
    timestamp: new Date(Date.now() - 45 * 60 * 1000).toISOString(),
    participants: [{ id: "user-10", name: "Lisa Designer" }],
    messages: [
      {
        id: "msg-8-1",
        sender: "other",
        senderInfo: { id: "user-10", name: "Lisa Designer" },
        content:
          "Hey! The new dashboard mockups are ready for review. I've incorporated the feedback from last week's session.\n\nMain changes:\n• Simplified navigation\n• New color scheme for better accessibility\n• Added quick action buttons\n\nLet me know your thoughts!",
        timestamp: new Date(Date.now() - 45 * 60 * 1000).toISOString(),
        type: "text",
        attachments: [
          { id: "att-3", name: "Dashboard_v2_Mockups.fig", type: "application/figma", size: 8900000, url: "/files/dashboard.fig" },
        ],
        aiInsights: [
          {
            id: "insight-5",
            type: "draft",
            label: "✨ Draft Available",
            content:
              "Hi Lisa,\n\nThese look great! The simplified navigation is a big improvement. I especially like the new color scheme.\n\nA few thoughts:\n• Can we make the quick action buttons more prominent?\n• The spacing on the sidebar looks a bit tight\n\nOverall, excellent work! Let's sync tomorrow to discuss.\n\nThanks!",
            isDraft: true,
          },
        ],
      },
    ],
  },

  // Task/Jira
  {
    id: "item-9",
    title: "JIRA-1234: Fix authentication bug",
    source: "task",
    priority: "high",
    unread: true,
    snippet: "Critical: Users unable to login with SSO",
    timestamp: new Date(Date.now() - 1 * 60 * 60 * 1000).toISOString(),
    participants: [{ id: "user-11", name: "QA Team" }],
    messages: [
      {
        id: "msg-9-1",
        sender: "other",
        senderInfo: { id: "user-11", name: "QA Team" },
        content:
          "**Bug Report: JIRA-1234**\n\n**Priority:** Critical\n**Status:** In Progress\n**Assignee:** You\n\n**Description:**\nUsers are unable to login using SSO. The authentication flow fails at the callback step.\n\n**Steps to Reproduce:**\n1. Click 'Login with SSO'\n2. Enter credentials\n3. Observe error on callback\n\n**Expected:** Successful login\n**Actual:** 500 error on callback",
        timestamp: new Date(Date.now() - 1 * 60 * 60 * 1000).toISOString(),
        type: "text",
        aiInsights: [
          {
            id: "insight-6",
            type: "analysis",
            label: "🔍 Analysis",
            content:
              "Based on similar issues in the codebase, this could be related to:\n\n1. **Token expiration** - The SSO token might be expiring before callback\n2. **Redirect URL mismatch** - Check the OAuth config\n3. **Recent changes** - Last commit to auth module was 2 days ago\n\nRecommended: Check the auth middleware logs first.",
            isDraft: false,
          },
        ],
        fullContent:
          '<div class="jira-ticket"><h1>JIRA-1234: Fix authentication bug</h1><div class="field"><label>Priority:</label><span class="critical">Critical</span></div><div class="field"><label>Status:</label><span>In Progress</span></div><div class="description"><h2>Description</h2><p>Users are unable to login using SSO. The authentication flow fails at the callback step with a 500 error.</p></div></div>',
      },
    ],
  },

  // Medium Priority Email
  {
    id: "item-10",
    title: "Weekly Newsletter - Tech Digest",
    source: "email",
    priority: "low",
    unread: false,
    snippet: "This week: AI breakthroughs, startup news, and more...",
    timestamp: new Date(Date.now() - 12 * 60 * 60 * 1000).toISOString(),
    participants: [{ id: "newsletter-1", name: "Tech Digest" }],
    messages: [
      {
        id: "msg-10-1",
        sender: "other",
        senderInfo: { id: "newsletter-1", name: "Tech Digest" },
        content:
          "📰 **This Week in Tech**\n\n• OpenAI announces GPT-5 preview\n• Apple's new AR glasses revealed\n• Bitcoin hits new all-time high\n• SpaceX Starship successful landing\n\nRead the full digest →",
        timestamp: new Date(Date.now() - 12 * 60 * 60 * 1000).toISOString(),
        type: "text",
      },
    ],
  },
];
