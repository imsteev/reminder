import React from "react";
import { format } from "date-fns";
import { Button } from "../../components/ui";
import { Reminder } from "../../api/reminders";
import { cn } from "../../utils/cn";

interface ReminderCardProps {
  reminder: Reminder;
  reminderTime: Date;
  timelineInfo: { distance: string };
  isUpcoming: boolean;
  onDelete: (id: number) => void;
  isDeleting: boolean;
}

export default function ReminderCard({
  reminder,
  reminderTime,
  timelineInfo,
  isUpcoming,
  onDelete,
  isDeleting,
}: ReminderCardProps) {
  return (
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
                  !reminder.message && "text-gray-600"
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
                <span className="font-medium">{timelineInfo.distance}</span>
                {reminder.type === "repeating" && (
                  <span className="text-xs text-blue-600">
                    • Next occurrence
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
            onClick={() => onDelete(reminder.id)}
            disabled={isDeleting}
            className="ml-4 text-red-600 hover:text-red-800"
          >
            🗑️
          </Button>
        </div>
      </div>
    </div>
  );
}
