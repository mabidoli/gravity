import express from "express";
import cors from "cors";
import dotenv from "dotenv";
import { streamRouter } from "./routes/stream.js";
import { messagesRouter } from "./routes/messages.js";

dotenv.config();

const app = express();
const PORT = process.env.PORT || 3001;

// Middleware
app.use(cors());
app.use(express.json());

// Routes
app.use("/api/stream", streamRouter);
app.use("/api/messages", messagesRouter);

// Health check
app.get("/api/health", (req, res) => {
  res.json({ status: "ok", timestamp: new Date().toISOString() });
});

app.listen(PORT, () => {
  console.log(`🚀 Gravity Backend running on http://localhost:${PORT}`);
});
