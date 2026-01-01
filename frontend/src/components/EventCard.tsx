"use client";

import { CalendarEvent } from "@/types";
import { Calendar, Clock, MapPin, Users, Video } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";

interface EventCardProps {
  event: CalendarEvent;
}

export function EventCard({ event }: EventCardProps) {
  const handleJoin = () => {
    if (event.meetingLink) {
      window.open(event.meetingLink, "_blank");
    }
  };

  return (
    <div className="glass-panel p-4 space-y-3 max-w-md animate-fade-in">
      {/* Header */}
      <div className="flex items-start gap-3">
        <div className="h-10 w-10 rounded-lg bg-teal-500/20 flex items-center justify-center">
          <Calendar size={20} className="text-teal-400" />
        </div>
        <div className="flex-1 min-w-0">
          <h3 className="font-semibold text-white truncate">{event.title}</h3>
          <div className="flex items-center gap-2 text-sm text-slate-400 mt-1">
            <Clock size={14} />
            <span>
              {event.startTime} - {event.endTime}
            </span>
          </div>
        </div>
      </div>

      {/* Location */}
      {event.location && (
        <div className="flex items-center gap-2 text-sm text-slate-400">
          <MapPin size={14} className="shrink-0" />
          <span className="truncate">{event.location}</span>
        </div>
      )}

      {/* Description */}
      {event.description && (
        <p className="text-sm text-slate-400 line-clamp-2">{event.description}</p>
      )}

      {/* Attendees */}
      {event.attendees.length > 0 && (
        <div className="flex items-center gap-2">
          <Users size={14} className="text-slate-500" />
          <div className="flex -space-x-2">
            {event.attendees.slice(0, 4).map((attendee) => (
              <Avatar
                key={attendee.id}
                className="h-6 w-6 border-2 border-slate-900"
              >
                <AvatarFallback className="text-[10px]">
                  {attendee.name.charAt(0)}
                </AvatarFallback>
              </Avatar>
            ))}
          </div>
          {event.attendees.length > 4 && (
            <span className="text-xs text-slate-500">
              +{event.attendees.length - 4} more
            </span>
          )}
        </div>
      )}

      {/* Join Button */}
      {event.meetingLink && (
        <Button
          onClick={handleJoin}
          className="w-full gap-2"
          variant="default"
          size="sm"
        >
          <Video size={16} />
          Join Meeting
        </Button>
      )}
    </div>
  );
}
