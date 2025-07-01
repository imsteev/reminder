import React, { useState, useEffect } from "react";
import { useMutation } from "@tanstack/react-query";
import { Button } from "./ui";
import { deleteReminder, Reminder } from "../api/reminders";
import { differenceInMinutes, isBefore } from "date-fns";
import Timeline from "./Timeline";

interface ReminderListProps {
  reminders?: Reminder[];
  isLoading: boolean;
  error: any;
  onDelete: () => void;
  completedOnly?: boolean;
  includeShowCompletedButton?: boolean;
}

const ReminderList: React.FC<ReminderListProps> = ({
  reminders,
  isLoading,
  error,
  onDelete,
  completedOnly,
  includeShowCompletedButton,
}) => {
  const [showPast, setShowPast] = useState(false);
  const [currentTime, setCurrentTime] = useState(new Date());

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date());
    }, 30000); // Update every 30 seconds
    return () => clearInterval(timer);
  }, []);

  const deleteMutation = useMutation({
    mutationFn: deleteReminder,
    onSuccess: () => {
      onDelete();
    },
  });

  if (isLoading) {
    return <div className="text-center py-4">Loading reminders...</div>;
  }

  if (error)
    return (
      <div className="text-red-500 text-center py-4">
        Error loading reminders
      </div>
    );

  const filteredReminders = reminders?.filter((reminder) => {
    if (completedOnly) {
      return reminder.type === "one-time" && isReminderPast(reminder);
    }
    // Always show repeating reminders, only filter one-time reminders based on showPast
    if (reminder.type === "repeating") {
      return true;
    }
    return showPast || !isReminderPast(reminder);
  });

  const activeReminders = showPast
    ? filteredReminders
    : filteredReminders?.filter(
        (reminder) => reminder.type === "repeating" || !isReminderPast(reminder)
      );

  const pastReminderCount =
    filteredReminders?.filter(isReminderPast).length || 0;

  // Sort reminders by next occurrence time for timeline
  const sortedReminders = [...(activeReminders || [])].sort(
    (a, b) =>
      getNextOccurrence(a, currentTime).getTime() -
      getNextOccurrence(b, currentTime).getTime()
  );

  return (
    <div className="w-full">
      {pastReminderCount > 0 && includeShowCompletedButton && (
        <div className="mb-4">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setShowPast(!showPast)}
            className="text-gray-500"
          >
            {showPast ? "Hide" : "Show"} completed reminders (
            {pastReminderCount})
          </Button>
        </div>
      )}

      {sortedReminders.length === 0 ? (
        <p className="text-gray-500 text-center py-8">
          No reminders yet. Create your first one!
        </p>
      ) : (
        <Timeline
          reminders={sortedReminders}
          currentTime={currentTime}
          onDelete={(id) => deleteMutation.mutate(id)}
          isDeleting={deleteMutation.isPending}
        />
      )}
    </div>
  );
};

const isReminderPast = (reminder: Reminder) => {
  if (reminder.type === "repeating") return false;
  return new Date(reminder.start_time) < new Date();
};

const getNextOccurrence = (reminder: Reminder, currentTime: Date) => {
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

export default ReminderList;
