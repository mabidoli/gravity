// Package model defines the core domain models for the Gravity BFF.
// These structs represent the data structures used throughout the application.
package model

import (
	"time"
)

// SourceType represents the origin platform of a priority item.
type SourceType string

const (
	SourceEmail    SourceType = "email"
	SourceWhatsApp SourceType = "whatsapp"
	SourceSlack    SourceType = "slack"
	SourceTeams    SourceType = "teams"
	SourceCalendar SourceType = "calendar"
	SourceTask     SourceType = "task"
	SourceYouTube  SourceType = "youtube"
	SourceLinkedIn SourceType = "linkedin"
	SourceTwitter  SourceType = "twitter"
)

// Priority represents the urgency level of a priority item.
type Priority string

const (
	PriorityHigh   Priority = "high"
	PriorityMedium Priority = "medium"
	PriorityLow    Priority = "low"
)

// SenderType represents who sent a message.
type SenderType string

const (
	SenderUser   SenderType = "user"
	SenderOther  SenderType = "other"
	SenderSystem SenderType = "system"
)

// ContentType represents the type of message content.
type ContentType string

const (
	ContentText   ContentType = "text"
	ContentEvent  ContentType = "event"
	ContentSocial ContentType = "social"
)

// InsightType represents the type of AI insight.
type InsightType string

const (
	InsightSuggestion InsightType = "suggestion"
	InsightAnalysis   InsightType = "analysis"
	InsightDraft      InsightType = "draft"
)

// SocialPlatform represents supported social media platforms.
type SocialPlatform string

const (
	PlatformYouTube  SocialPlatform = "youtube"
	PlatformLinkedIn SocialPlatform = "linkedin"
	PlatformTwitter  SocialPlatform = "twitter"
)
