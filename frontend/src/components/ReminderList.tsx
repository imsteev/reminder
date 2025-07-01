import React, { useState, useEffect } from "react";
import { useMutation } from "@tanstack/react-query";
import { Button } from "./ui";
import { deleteReminder, Reminder } from "../api/reminders";
import {
  format,
  formatDistanceToNow,
  differenceInMinutes,
  isBefore,
} from "date-fns";

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
    return showPast || !isReminderPast(reminder);
  });

  const activeReminders = showPast
    ? filteredReminders
    : filteredReminders?.filter((reminder) => !isReminderPast(reminder));

  const pastReminderCount =
    filteredReminders?.filter(isReminderPast).length || 0;

  // Sort reminders by start time for timeline
  const sortedReminders = [...(activeReminders || [])].sort(
    (a, b) =>
      new Date(a.start_time).getTime() - new Date(b.start_time).getTime(),
  );

  // Find where "now" should be inserted in the timeline
  const nowIndex = sortedReminders.findIndex(
    (reminder) => new Date(reminder.start_time) > currentTime,
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
              const reminderTime = new Date(reminder.start_time);
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
                  <div className="relative flex items-start">
                    {/* Timeline dot */}
                    <div
                      className={`w-4 h-4 rounded-full border-2 border-white shadow-md z-10 ${
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
                              className={`font-medium ${isUpcoming ? "text-gray-900" : "text-gray-600"}`}
                            >
                              {reminder.message}
                            </p>
                            <span
                              className={`px-2 py-1 text-xs font-medium rounded-full ${
                                reminder.type === "repeating"
                                  ? isUpcoming
                                    ? "bg-blue-100 text-blue-800"
                                    : "bg-gray-200 text-gray-600"
                                  : isUpcoming
                                    ? "bg-green-100 text-green-800"
                                    : "bg-gray-200 text-gray-600"
                              }`}
                            >
                              {reminder.type === "repeating"
                                ? "Repeats"
                                : "One-time"}
                            </span>
                          </div>

                          <div className="text-sm space-y-1">
                            <div
                              className={`flex items-center gap-2 ${isUpcoming ? "text-gray-700" : "text-gray-500"}`}
                            >
                              <span className="font-medium">
                                {timelineInfo.distance}
                              </span>
                              {reminder.type === "repeating" && (
                                <span className="text-xs">
                                  • Every {reminder.period_minutes} min
                                </span>
                              )}
                            </div>
                            <div
                              className={`text-xs ${isUpcoming ? "text-gray-500" : "text-gray-400"}`}
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
