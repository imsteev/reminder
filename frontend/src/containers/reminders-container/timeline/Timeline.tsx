import React from "react";
import { Reminder } from "../../../api/reminders";
import NowMarker from "./NowMarker";
import ReminderCard from "../ReminderCard";
import { getNextOccurrence, getTimelinePosition } from "./utils";
import TimelineDot from "./TimelineDot";

interface Props {
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
}: Props) {
  return (
    <div className="relative">
      <div className="space-y-6">
        <NowMarker currentTime={currentTime} />
        {reminders.map((reminder) => {
          const reminderTime = getNextOccurrence(reminder, currentTime);
          const timelineInfo = getTimelinePosition(reminderTime, currentTime);
          const isUpcoming = reminderTime > currentTime;
          return (
            <div key={reminder.id} className="flex items-center gap-2">
              <TimelineDot isUpcoming={isUpcoming} type={reminder.type} />
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
      </div>
    </div>
  );
}
