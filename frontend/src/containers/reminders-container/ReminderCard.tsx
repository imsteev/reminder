import React from "react";
import { format } from "date-fns";
import { Button, Card, CardContent } from "../../components/ui";
import { type Reminder } from "../../api/reminders";
import { cn } from "../../utils/cn";

interface ReminderCardProps {
  reminder: Reminder;
  reminderTime: Date;
  isUpcoming: boolean;
  onDelete: (id: number) => void;
  onEdit?: (reminder: Reminder) => void;
  isDeleting: boolean;
}

export default function ReminderCard({
  reminder,
  reminderTime,
  isUpcoming,
  onDelete,
  onEdit,
  isDeleting,
}: ReminderCardProps) {
  return (
    <Card
      className={cn(
        "flex-1",
        isUpcoming && reminder.is_repeating && "border-blue-200",
        isUpcoming && !reminder.is_repeating && "border-green-200"
      )}
    >
      <CardContent>
        <div className="flex justify-between items-center">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-2">
              <span className="text-xs text-gray-500">#{reminder.id}</span>
              <p className={cn("font-medium")}>
                {reminder.body?.slice(0, 30) +
                  (reminder.body?.length > 30 ? "..." : "") || "-"}
              </p>
              {reminder.is_repeating && (
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
                {reminder.is_repeating && (
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

          <div className="flex items-center gap-2 ml-4">
            {!reminder.is_repeating && (
              <Button
                variant="ghost"
                size="sm"
                onClick={() => onEdit?.(reminder)}
                disabled={isDeleting}
              >
                Edit
              </Button>
            )}
            <Button
              variant="ghost"
              size="sm"
              onClick={() => onDelete(reminder.id)}
              disabled={isDeleting}
              className="text-red-600 hover:text-red-800"
            >
              Delete
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
