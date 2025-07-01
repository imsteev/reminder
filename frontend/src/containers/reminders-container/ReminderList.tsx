import React, { useState, useEffect } from "react";
import { useMutation } from "@tanstack/react-query";
import { Button } from "../../components/ui";
import { deleteReminder, type Reminder } from "../../api/reminders";
import Timeline from "./timeline/Timeline";
import { getNextOccurrence } from "./timeline/utils";

interface Props {
  reminders?: Reminder[];
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
      refetch();
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
    return <div className="text-center text-gray-500 py-24">No reminders</div>;
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
