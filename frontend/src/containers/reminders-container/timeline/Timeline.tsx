import React, { useState } from "react";
import { type Reminder } from "../../../api/reminders";
import ReminderCard from "../ReminderCard";
import { getNextOccurrence } from "./utils";
import {
  Dialog,
  DialogHeader,
  DialogTitle,
  DialogBody,
} from "../../../components/ui";
import { ReminderForm } from "../../../components/reminder-form";
import {
  getCurrentDateTimeString,
  toLocalDateTimeString,
} from "../../../utils/datetime";

interface Props {
  reminders: Reminder[];
  currentTime: Date;
  onDelete: (id: number) => void;
  onSuccess: () => void;
  isDeleting: boolean;
}

export default function Timeline({
  reminders,
  currentTime,
  onDelete,
  onSuccess,
  isDeleting,
}: Props) {
  const [showForm, setShowForm] = useState(false);
  const [editData, setEditData] = useState<Reminder | null>(null);

  const handleEdit = (reminder: Reminder) => {
    setEditData(reminder);
    setShowForm(true);
  };

  const convertReminderToFormData = (reminder: Reminder) => {
    const totalMinutes = reminder.period_minutes;
    const days = Math.floor(totalMinutes / (24 * 60));
    const hours = Math.floor((totalMinutes % (24 * 60)) / 60);
    const minutes = totalMinutes % 60;

    return {
      body: reminder.body || "",
      isRepeating: reminder.is_repeating,
      contactMethodID: reminder.contact_method_id,
      intervalDays: days,
      intervalHours: hours,
      intervalMinutes: minutes,
      startTime: toLocalDateTimeString(new Date(reminder.start_time)),
    };
  };
  return (
    <div className="relative">
      <div className="space-y-6">
        {reminders.map((reminder) => {
          const reminderTime = getNextOccurrence(reminder, currentTime);
          const isUpcoming = reminderTime > currentTime;
          return (
            <div key={reminder.id} className="flex items-center gap-2">
              <ReminderCard
                reminder={reminder}
                reminderTime={reminderTime}
                isUpcoming={isUpcoming}
                onDelete={onDelete}
                onEdit={handleEdit}
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
            setEditData(null);
          }}
        >
          <DialogHeader>
            <DialogTitle>
              {editData ? "Edit Reminder" : "Create New Reminder"}
            </DialogTitle>
          </DialogHeader>
          <DialogBody>
            <ReminderForm
              reminderID={editData?.id}
              initialData={
                editData ? convertReminderToFormData(editData) : undefined
              }
              onSuccess={() => {
                onSuccess();
                setShowForm(false);
                setEditData(null);
              }}
            />
          </DialogBody>
        </Dialog>
      )}
    </div>
  );
}
