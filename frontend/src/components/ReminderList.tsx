import React, { useState, useEffect } from "react";
import { useMutation } from "@tanstack/react-query";
import { Button } from "./ui";
import { deleteReminder, Reminder } from "../api/reminders";
import {
  format,
  formatDistanceToNow,
  differenceInMinutes,
  isBefore,
  addMinutes,
} from "date-fns";
import { cn } from "../utils/cn";

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

  const isReminderPast = (reminder: Reminder) => {
    if (reminder.type === "repeating") return false;
    return new Date(reminder.start_time) < new Date();
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

  if (isLoading)
    return <div className="text-center py-4">Loading reminders...</div>;
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
        (reminder) =>
          reminder.type === "repeating" || !isReminderPast(reminder),
      );

  const pastReminderCount =
    filteredReminders?.filter(isReminderPast).length || 0;

  // Sort reminders by next occurrence time for timeline
  const sortedReminders = [...(activeReminders || [])].sort(
    (a, b) => getNextOccurrence(a).getTime() - getNextOccurrence(b).getTime(),
  );

  // Find where "now" should be inserted in the timeline
  const nowIndex = sortedReminders.findIndex(
    (reminder) => getNextOccurrence(reminder) > currentTime,
  );
  const showNowMarker = sortedReminders.length > 0;

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
        <div className="relative">
          {/* Timeline line */}
          <div className="absolute left-8 top-0 bottom-0 w-0.5 bg-gray-200"></div>

          <div className="space-y-6">
            {sortedReminders.map((reminder, index) => {
              const reminderTime = getNextOccurrence(reminder);
              const timelineInfo = getTimelinePosition(reminderTime);
              const isUpcoming = reminderTime > currentTime;

              return (
                <div key={reminder.id}>
                  {/* Show NOW marker before the first upcoming reminder */}
                  {showNowMarker && index === nowIndex && (
                    <div className="relative flex items-center mb-6">
                      <div className="flex items-center">
                        <div className="w-4 h-4 bg-red-500 rounded-full border-2 border-white shadow-md z-10"></div>
                        <div className="ml-4 bg-red-50 border border-red-200 rounded-lg px-4 py-2">
                          <div className="flex items-center gap-2">
                            <span className="text-red-600 font-semibold text-sm">
                              NOW
                            </span>
                            <span className="text-red-500 text-xs">
                              {format(currentTime, "h:mm a")}
                            </span>
                          </div>
                        </div>
                      </div>
                    </div>
                  )}

                  {/* Reminder card */}
                  <div className="relative flex items-center">
                    {/* Timeline dot */}
                    <div
                      className={`w-4 h-4 rounded-full border-2 border-white shadow-md z-10 flex-shrink-0 ${
                        isUpcoming
                          ? reminder.type === "repeating"
                            ? "bg-blue-500"
                            : "bg-green-500"
                          : "bg-gray-400"
                      }`}
                    ></div>

                    {/* Reminder card */}
                    <div
                      className={`ml-4 flex-1 p-4 rounded-lg shadow border ${
                        isUpcoming
                          ? reminder.type === "repeating"
                            ? "bg-white border-blue-200"
                            : "bg-white border-green-200"
                          : "bg-gray-50 border-gray-300"
                      }`}
                    >
                      <div className="flex justify-between items-start">
                        <div className="flex-1">
                          <div className="flex items-center gap-2 mb-2">
                            <p
                              className={cn(
                                "font-medium",
                                !reminder.message && "text-gray-600",
                              )}
                            >
                              {reminder.message || "-"}
                            </p>
                            {reminder.type === "repeating" && (
                              <span
                                className={`px-2 py-1 text-xs font-medium rounded-full ${
                                  isUpcoming
                                    ? "bg-blue-100 text-blue-800"
                                    : "bg-gray-200 text-gray-600"
                                }`}
                              >
                                Repeats every {reminder.period_minutes}m
                              </span>
                            )}
                          </div>

                          <div className="text-sm space-y-1">
                            <div
                              className={`flex items-center gap-2 ${
                                isUpcoming ? "text-gray-700" : "text-gray-500"
                              }`}
                            >
                              {reminder.type === "repeating" && (
                                <span className="text-xs text-blue-600">
                                  Next occurrence in {timelineInfo.distance}
                                </span>
                              )}
                            </div>
                            <div
                              className={`text-xs ${
                                isUpcoming ? "text-gray-500" : "text-gray-400"
                              }`}
                            >
                              {format(reminderTime, "MMM d, yyyy h:mm a")}
                            </div>
                          </div>
                        </div>

                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => deleteMutation.mutate(reminder.id)}
                          disabled={deleteMutation.isPending}
                          className="ml-4 text-red-600 hover:text-red-800"
                        >
                          🗑️
                        </Button>
                      </div>
                    </div>
                  </div>
                </div>
              );
            })}

            {/* Show NOW marker at the end if all reminders are in the past */}
            {showNowMarker && nowIndex === -1 && (
              <div className="relative flex items-center">
                <div className="flex items-center">
                  <div className="w-4 h-4 bg-red-500 rounded-full border-2 border-white shadow-md z-10"></div>
                  <div className="ml-4 bg-red-50 border border-red-200 rounded-lg px-4 py-2">
                    <div className="flex items-center gap-2">
                      <span className="text-red-600 font-semibold text-sm">
                        NOW
                      </span>
                      <span className="text-red-500 text-xs">
                        {format(currentTime, "h:mm a")}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default ReminderList;
