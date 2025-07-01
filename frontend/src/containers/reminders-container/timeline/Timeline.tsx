import React from "react";
import { Reminder } from "../../../api/reminders";
import { formatDistanceToNow, differenceInMinutes, isBefore } from "date-fns";
import TimelineMarker from "./TimelineMarker";
import ReminderCard from "../ReminderCard";
import { getNextOccurrence, getTimelinePosition } from "./utils";

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
  const nowIndex = reminders.findIndex(
    (reminder) => getNextOccurrence(reminder, currentTime) > currentTime
  );
  const showNowMarker = reminders.length > 0;

  return (
    <div className="relative">
      {/* Timeline line */}
      <div className="absolute left-8 top-0 bottom-0 w-0.5 bg-gray-200"></div>

      <div className="space-y-6">
        {reminders.map((reminder, index) => {
          const reminderTime = getNextOccurrence(reminder, currentTime);
          const timelineInfo = getTimelinePosition(reminderTime, currentTime);
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
