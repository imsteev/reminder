import { differenceInMinutes, isBefore, formatDistanceToNow } from "date-fns";
import { type Reminder } from "../../../api/reminders";

export const getNextOccurrence = (reminder: Reminder, currentTime: Date) => {
  const startTime = new Date(reminder.start_time);

  if (!reminder.is_repeating) {
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
