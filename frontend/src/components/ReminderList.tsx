import React, { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { Button } from "./ui";
import { deleteReminder, Reminder } from "../api/reminders";
import { format } from "date-fns";

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

  return (
    <div className="space-y-4 w-full flex flex-col items-center justify-center">
      {pastReminderCount > 0 && includeShowCompletedButton && (
        <Button
          variant="outline"
          size="sm"
          onClick={() => setShowPast(!showPast)}
          className="text-gray-500"
        >
          {showPast ? "Hide" : "Show"} completed reminders ({pastReminderCount})
        </Button>
      )}
      {activeReminders?.map((reminder) => {
        return (
          <div
            key={reminder.id}
            className={`w-full p-4 rounded-lg shadow border ${
              reminder.type === "repeating"
                ? "bg-white border-blue-200"
                : "bg-white border-gray-200"
            }`}
          >
            <div className="flex justify-between items-start">
              <div>
                <div className="flex items-center gap-2 mb-2">
                  <p className={`font-medium ${"text-gray-900"}`}>
                    {reminder.message}
                  </p>
                  <span
                    className={`px-2 py-1 text-xs font-medium rounded-full ${
                      reminder.type === "repeating"
                        ? "bg-blue-100 text-blue-800"
                        : "bg-gray-100 text-gray-800"
                    }`}
                  >
                    {reminder.type === "repeating" && "Repeats"}
                  </span>
                </div>
                <div className="text-sm text-gray-600 space-y-1">
                  {reminder.type === "one-time" && (
                    <Button variant="outline" size="sm">
                      {format(
                        new Date(reminder.start_time),
                        "MMM d, yyyy h:mm a"
                      )}
                    </Button>
                  )}

                  {reminder.type === "repeating" && (
                    <>
                      <p>Every {reminder.period_minutes} minutes</p>
                      <p>
                        <Button variant="outline" size="sm">
                          {format(
                            new Date(reminder.start_time),
                            "MMM d, yyyy h:mm a"
                          )}
                        </Button>
                      </p>
                    </>
                  )}
                </div>
              </div>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => deleteMutation.mutate(reminder.id)}
                disabled={deleteMutation.isPending}
                className="ml-4"
              >
                X
              </Button>
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default ReminderList;
