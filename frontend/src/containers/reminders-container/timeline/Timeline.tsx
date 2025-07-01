import React, { useState } from "react";
import { Reminder } from "../../../api/reminders";
import NowMarker from "./NowMarker";
import ReminderCard from "../ReminderCard";
import { getNextOccurrence, getTimelinePosition } from "./utils";
import TimelineDot from "./TimelineDot";
import {
  Button,
  Dialog,
  DialogHeader,
  DialogTitle,
  DialogBody,
} from "../../../components/ui";
import ReminderForm from "../ReminderForm";

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
  return (
    <div className="relative">
      <div className="space-y-6">
        <div className="flex items-center justify-between gap-4">
          <NowMarker currentTime={currentTime} />
          <div className="space-y-4 flex justify-center">
            <Button variant="blue" onClick={() => setShowForm(true)}>
              <span className="mr-2">+</span>
              New Reminder
            </Button>
          </div>
        </div>
        {reminders.map((reminder) => {
          const reminderTime = getNextOccurrence(reminder, currentTime);
          const timelineInfo = getTimelinePosition(reminderTime, currentTime);
          const isUpcoming = reminderTime > currentTime;
          return (
            <div key={reminder.id} className="flex items-center gap-2">
              <TimelineDot isUpcoming={isUpcoming} type={reminder.type} />
              <ReminderCard
                reminder={reminder}
                reminderTime={reminderTime}
                timelineInfo={timelineInfo}
                isUpcoming={isUpcoming}
                onDelete={onDelete}
                isDeleting={isDeleting}
              />
            </div>
          );
        })}
      </div>
      <Dialog isOpen={showForm} onClose={() => setShowForm(false)}>
        <DialogHeader>
          <DialogTitle>Create New Reminder</DialogTitle>
        </DialogHeader>
        <DialogBody>
          <ReminderForm
            onSuccess={() => {
              onCreateSuccess();
              setShowForm(false);
            }}
          />
        </DialogBody>
      </Dialog>
    </div>
  );
}
