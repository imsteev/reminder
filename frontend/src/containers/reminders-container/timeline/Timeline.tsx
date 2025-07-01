import React, { useState } from "react";
import { Reminder } from "../../../api/reminders";
import ReminderCard from "../ReminderCard";
import { getNextOccurrence, getTimelinePosition } from "./utils";
import {
  Dialog,
  DialogHeader,
  DialogTitle,
  DialogBody,
} from "../../../components/ui";
import ReminderForm from "../ReminderForm";
import { getCurrentDateTimeString } from "../../../utils/datetime";

interface Props {
  reminders: Reminder[];
  currentTime: Date;
  onDelete: (id: number) => void;
  onCreateSuccess: () => void;
  isDeleting: boolean;
}

export default function Timeline({
  reminders,
  currentTime,
  onDelete,
  onCreateSuccess,
  isDeleting,
}: Props) {
  const [showForm, setShowForm] = useState(false);
  const [rescheduleData, setRescheduleData] = useState<Reminder | null>(null);

  const handleReschedule = (reminder: Reminder) => {
    setRescheduleData(reminder);
    setShowForm(true);
  };

  const convertReminderToFormData = (reminder: Reminder) => {
    const totalMinutes = reminder.period_minutes;
    const days = Math.floor(totalMinutes / (24 * 60));
    const hours = Math.floor((totalMinutes % (24 * 60)) / 60);
    const minutes = totalMinutes % 60;

    return {
      name: reminder.name || "",
      message: reminder.message || "",
      reminderType: reminder.type as "one-time" | "repeating",
      deliveryType: reminder.delivery_type as "sms" | "email",
      intervalDays: days,
      intervalHours: hours,
      intervalMinutes: minutes,
      startTime: getCurrentDateTimeString(), // Set to now for reschedule
      contactValue: "", // Will be filled from contact methods if available
    };
  };
  return (
    <div className="relative">
      <div className="space-y-6">
        {reminders.map((reminder) => {
          const reminderTime = getNextOccurrence(reminder, currentTime);
          const timelineInfo = getTimelinePosition(reminderTime, currentTime);
          const isUpcoming = reminderTime > currentTime;
          return (
            <div key={reminder.id} className="flex items-center gap-2">
              <ReminderCard
                reminder={reminder}
                reminderTime={reminderTime}
                timelineInfo={timelineInfo}
                isUpcoming={isUpcoming}
                onDelete={onDelete}
                onReschedule={handleReschedule}
                isDeleting={isDeleting}
              />
            </div>
          );
        })}
      </div>
      {showForm && (
        <Dialog
          isOpen={showForm}
          onClose={() => {
            setShowForm(false);
            setRescheduleData(null);
          }}
        >
          <DialogHeader>
            <DialogTitle>
              {rescheduleData ? "Reschedule Reminder" : "Create New Reminder"}
            </DialogTitle>
          </DialogHeader>
          <DialogBody>
            <ReminderForm
              initialData={
                rescheduleData
                  ? convertReminderToFormData(rescheduleData)
                  : undefined
              }
              onSuccess={() => {
                onCreateSuccess();
                setShowForm(false);
                setRescheduleData(null);
              }}
            />
          </DialogBody>
        </Dialog>
      )}
    </div>
  );
}
