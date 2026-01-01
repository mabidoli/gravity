import { Router, Request, Response } from "express";
import { PriorityItem } from "../types/index.js";
import { mockStreamItems } from "../data/mockData.js";

export const streamRouter = Router();

// Get all stream items
streamRouter.get("/", (req: Request, res: Response) => {
  const { filter, priority } = req.query;

  let items = [...mockStreamItems];

  // Filter by priority
  if (priority && typeof priority === "string") {
    items = items.filter((item) => item.priority === priority);
  }

  // Filter by unread
  if (filter === "unread") {
    items = items.filter((item) => item.unread);
  }

  // Sort by priority (high first) then by timestamp
  items.sort((a, b) => {
    const priorityOrder = { high: 0, medium: 1, low: 2 };
    if (priorityOrder[a.priority] !== priorityOrder[b.priority]) {
      return priorityOrder[a.priority] - priorityOrder[b.priority];
    }
    return new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime();
  });

  res.json(items);
});

// Get single stream item by ID
streamRouter.get("/:id", (req: Request, res: Response) => {
  const { id } = req.params;
  const item = mockStreamItems.find((item) => item.id === id);

  if (!item) {
    res.status(404).json({ error: "Item not found" });
    return;
  }

  res.json(item);
});

// Mark item as read
streamRouter.patch("/:id/read", (req: Request, res: Response) => {
  const { id } = req.params;
  const itemIndex = mockStreamItems.findIndex((item) => item.id === id);

  if (itemIndex === -1) {
    res.status(404).json({ error: "Item not found" });
    return;
  }

  mockStreamItems[itemIndex].unread = false;
  res.json(mockStreamItems[itemIndex]);
});
