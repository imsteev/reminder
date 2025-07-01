import React from "react";
import { Reminder } from "../api/reminders";
import { formatDistanceToNow, differenceInMinutes, isBefore } from "date-fns";
import TimelineMarker from "./TimelineMarker";
import ReminderCard from "./ReminderCard";

interface TimelineProps {
  reminders: Reminder[];
  currentTime: Date;
  onDelete: (id: number) => void;
  isDeleting: boolean;
}

export default function Timeline({
  reminders,
  currentTime,
  onDelete,
  isDeleting,
}: TimelineProps) {
  const getNextOccurrence = (reminder: Reminder) => {
    const startTime = new Date(reminder.start_time);

    if (reminder.type === "one-time") {
      return startTime;
    }

    // For repeating reminders, calculate the next occurrence
    const now = currentTime;
    const periodMs = reminder.period_minutes * 60 * 1000;

    if (startTime > now) {
      // If start time is in the future, return start time
      return startTime;
    }

    // Calculate how many periods have passed since start time
    const timeSinceStart = now.getTime() - startTime.getTime();
    const periodsPassed = Math.floor(timeSinceStart / periodMs);

    // Calculate next occurrence
    const nextOccurrence = new Date(
      startTime.getTime() + (periodsPassed + 1) * periodMs,
    );

    return nextOccurrence;
  };

  const getTimelinePosition = (reminderTime: Date) => {
    const minutesToReminder = differenceInMinutes(reminderTime, currentTime);
    const isPast = isBefore(reminderTime, currentTime);

    if (isPast) return { position: -1, distance: "Past" };
    if (minutesToReminder === 0) return { position: 0, distance: "Now" };
    if (minutesToReminder <= 30)
      return { position: 1, distance: formatDistanceToNow(reminderTime) };
    if (minutesToReminder <= 180)
      return { position: 2, distance: formatDistanceToNow(reminderTime) };
    if (minutesToReminder <= 1440)
      return { position: 3, distance: formatDistanceToNow(reminderTime) };
    return { position: 4, distance: formatDistanceToNow(reminderTime) };
  };

  // Find where "now" should be inserted in the timeline
  const nowIndex = reminders.findIndex(
    (reminder) => getNextOccurrence(reminder) > currentTime,
  );
  const showNowMarker = reminders.length > 0;

  return (
    <div className="relative">
      {/* Timeline line */}
      <div className="absolute left-8 top-0 bottom-0 w-0.5 bg-gray-200"></div>

      <div className="space-y-6">
        {reminders.map((reminder, index) => {
          const reminderTime = getNextOccurrence(reminder);
          const timelineInfo = getTimelinePosition(reminderTime);
          const isUpcoming = reminderTime > currentTime;

          return (
            <div key={reminder.id}>
              {/* Show NOW marker before the first upcoming reminder */}
              {showNowMarker && index === nowIndex && (
                <TimelineMarker currentTime={currentTime} />
              )}

              <ReminderCard
                reminder={reminder}
                reminderTime={reminderTime}
                timelineInfo={timelineInfo}
                isUpcoming={isUpcoming}
                onDelete={onDelete}
                isDeleting={isDeleting}
              />
            </div>
          );
        })}

        {/* Show NOW marker at the end if all reminders are in the past */}
        {showNowMarker && nowIndex === -1 && (
          <TimelineMarker currentTime={currentTime} />
        )}
      </div>
    </div>
  );
}
