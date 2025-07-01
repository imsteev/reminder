import { differenceInMinutes, isBefore, formatDistanceToNow } from "date-fns";
import { Reminder } from "../../../api/reminders";

export const getNextOccurrence = (reminder: Reminder, currentTime: Date) => {
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
    startTime.getTime() + (periodsPassed + 1) * periodMs
  );

  return nextOccurrence;
};

export const getTimelinePosition = (reminderTime: Date, currentTime: Date) => {
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
