import React, { useState, useEffect } from "react";
import { useMutation } from "@tanstack/react-query";
import {
  Button,
  Dialog,
  DialogBody,
  DialogHeader,
  DialogTitle,
} from "../../components/ui";
import { deleteReminder, type Reminder } from "../../api/reminders";
import Timeline from "./timeline/Timeline";
import { getNextOccurrence } from "./timeline/utils";
import { toast } from "sonner";
import { ReminderForm } from "../../components/reminder-form";

interface Props {
  reminders: Reminder[];
  isLoading: boolean;
  error: any;
  refetch: () => void;
}

const ReminderList: React.FC<Props> = ({
  reminders,
  isLoading,
  error,
  refetch,
}) => {
  const [showForm, setShowForm] = useState(false);
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
      toast.success("Reminder deleted");
      refetch();
    },
    onError: (error) => {
      toast.error("Failed to delete reminder", {
        description: error.message,
      });
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

  const sortedReminders = reminders?.sort(
    (a, b) =>
      getNextOccurrence(a, currentTime).getTime() -
      getNextOccurrence(b, currentTime).getTime()
  );

  if (!isLoading && !sortedReminders?.length) {
    return (
      <div className="text-center py-24 flex flex-col items-center justify-center">
        <div className="text-xl font-medium mb-4 text-gray-500">
          No active reminders
        </div>
        <Button
          onClick={() => setShowForm(true)}
          className="bg-blue-600 text-white hover:bg-blue-700 font-medium"
          size="sm"
        >
          <span className="mr-2">+</span>
          Reminder
        </Button>

        {showForm && (
          <Dialog
            isOpen={showForm}
            onClose={() => {
              setShowForm(false);
            }}
          >
            <DialogHeader className="flex">
              <DialogTitle className="p-0 flex items-center justify-between">
                Create Reminder
              </DialogTitle>
            </DialogHeader>
            <DialogBody>
              <ReminderForm
                onSuccess={() => {
                  setShowForm(false);
                  refetch();
                }}
              />
            </DialogBody>
          </Dialog>
        )}
      </div>
    );
  }

  return (
    <div className="w-full">
      <Timeline
        reminders={sortedReminders || []}
        currentTime={currentTime}
        onDelete={(id) => deleteMutation.mutate(id)}
        isDeleting={deleteMutation.isPending}
        onCreateSuccess={refetch}
      />
    </div>
  );
};

export default ReminderList;
