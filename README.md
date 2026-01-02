# Gravity V2

A "Personal Infrastructure" interface that unifies communication, calendar, and social streams into a single conversational workspace.

## Overview

Gravity rejects the siloed app model in favor of a conversation-centric paradigm where every input—whether a human message, a calendar event, or a social media post—is treated as a "Chat Item" that can be interacted with.

## Core Philosophy

1. **Unified Stream**: All inputs (Email, Slack, WhatsApp, Calendar, YouTube, LinkedIn) coexist in a single prioritized list
2. **Everything is a Chat**: The primary interaction model is conversation
3. **Proactive but Unobtrusive AI**: AI provides context-aware value via non-blocking UI elements ("Smart Pills")
4. **Deep Context on Demand**: Full fidelity content is always one click away

## Tech Stack

### Frontend
- **Framework**: Next.js 14+ (App Router)
- **Language**: TypeScript (Strict Mode)
- **Styling**: Tailwind CSS 4 + Framer Motion
- **UI Primitives**: Radix UI
- **State Management**: Zustand
- **Icons**: Lucide React

### Backend
- **Runtime**: Node.js with TypeScript
- **Framework**: Express.js
- **API**: RESTful endpoints

## Repository Organization

This repository follows a **component-based organization** where all files related to a specific part of the system are contained within that component's directory:

- **Root Level**: Only overall project documentation, monorepo configuration, and cross-cutting concerns
- **`frontend/`**: All frontend-related code, documentation, and planning files
- **`backend/`**: All backend-related code, documentation, and planning files

> **Important**: When adding documentation, tasks, or planning files, always place them in the appropriate component directory (frontend/ or backend/), not at the root level.

## Project Structure

```
gravity/
├── frontend/              # Next.js frontend application
│   ├── src/
│   │   ├── app/           # Next.js App Router pages
│   │   ├── components/    # React components
│   │   ├── lib/           # Utility functions
│   │   ├── store/         # Zustand state management
│   │   ├── types/         # TypeScript interfaces
│   │   └── data/          # Mock data
│   └── package.json
├── backend/               # Backend API server
│   ├── src/
│   │   ├── routes/        # API routes
│   │   ├── types/         # TypeScript interfaces
│   │   └── data/          # Mock data
│   ├── plans/             # Backend planning documentation
│   ├── tasks.md           # Backend implementation tasks
│   ├── README-BFF.md      # Backend-for-Frontend documentation
│   └── package.json
└── package.json           # Root monorepo config
```

## Getting Started

### Prerequisites
- Node.js 18+
- npm

### Installation

```bash
# Install all dependencies
npm run install:all

# Or install individually
cd frontend && npm install
cd ../backend && npm install
```

### Authentication Setup (Clerk)

Gravity uses [Clerk](https://clerk.com) for authentication. Follow these steps to configure it:

1. **Create a Clerk Application**
   - Go to [Clerk Dashboard](https://dashboard.clerk.com) and sign in
   - Create a new application or select an existing one
   - Note your API keys from the dashboard

2. **Configure Frontend Environment**

   Create or update `frontend/.env.local`:
   ```bash
   # Required - Get from Clerk Dashboard
   NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_your_key_here
   CLERK_SECRET_KEY=sk_test_your_key_here

   # Optional - Auth redirect URLs (defaults shown)
   NEXT_PUBLIC_CLERK_SIGN_IN_URL=/sign-in
   NEXT_PUBLIC_CLERK_SIGN_UP_URL=/sign-up

   # API Configuration
   NEXT_PUBLIC_API_URL=http://localhost:8080
   ```

3. **Configure Backend Environment**

   Create or update `backend/gravity-bff/.env`:
   ```bash
   # Required for production - Get from Clerk Dashboard
   CLERK_SECRET_KEY=sk_test_your_key_here
   ```

4. **Development Mode (Optional)**

   For local development without Clerk, simply leave `CLERK_SECRET_KEY` unset in the backend. This enables development mode where:
   - Backend allows unauthenticated requests
   - Frontend still requires Clerk keys for UI components

5. **Clerk Dashboard Settings**

   In your Clerk dashboard, configure:
   - **Allowed origins**: Add `http://localhost:3000` for local development
   - **Sign-in/Sign-up URLs**: Match your environment variables

### Development

```bash
# Run both frontend and backend
npm run dev

# Run frontend only
npm run dev:frontend

# Run backend only
npm run dev:backend
```

The frontend will be available at `http://localhost:3000`
The backend API will be available at `http://localhost:3001`

### Build

```bash
npm run build
```

## Features

### Unified Stream (Left Panel)
- Aggregates Email, WhatsApp, Slack, Calendar, YouTube, LinkedIn, Twitter
- Priority-based sorting (High/Medium/Low)
- Filter tabs: All, High Priority, Unread

### Conversational Context (Right Panel)
- Message bubbles for human communication
- Rich Event Cards for calendar items
- Content Cards for social media
- AI-powered Smart Pills for insights and drafts

### AI Integration
- **Smart Contextual Pills**: "✨ Draft Available", "⚠️ Risk Detected"
- **Insight Panel**: View AI analysis, drafts, and suggestions
- **Refinement**: Chat with AI to refine drafts

### System Interactions
- **Calendar**: "Reschedule to 4 PM" updates events
- **Social**: "Summarize this video" generates summaries

## Design System: Glassmorphic Zen

- Deep teal/dark backgrounds
- Frosted glass panels with `backdrop-blur`
- Subtle inner borders and ambient glow effects
- Smooth hover transitions (200ms)

## License

MIT
