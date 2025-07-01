import { useQuery } from "@tanstack/react-query";
import React, { useState } from "react";
import { getReminders } from "../../api/reminders";
import ReminderForm from "./ReminderForm";
import ReminderList from "./ReminderList";
import { Button } from "../../components/ui";
import { DEFAULT_USER_ID } from "../../constants";

export default function RemindersContainer() {
  const [showForm, setShowForm] = useState(false);

  const {
    data: reminders,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ["reminders"],
    queryFn: () => getReminders(DEFAULT_USER_ID),
  });

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold text-gray-900">Reminders</h1>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <div className="space-y-4">
            <div>
              {!showForm ? (
                <div className="flex flex-col items-center text-center py-4">
                  <Button
                    variant="blue"
                    className="w-full"
                    onClick={() => setShowForm(true)}
                  >
                    <span className="mr-2">+</span>
                    Create New Reminder
                  </Button>
                </div>
              ) : (
                <div className="space-y-4 mt-4">
                  <div className="flex items-center justify-between w-full">
                    <Button
                      variant="outline"
                      className="w-full"
                      onClick={() => setShowForm(false)}
                    >
                      Cancel
                    </Button>
                  </div>
                  <ReminderForm
                    onSuccess={() => {
                      refetch();
                      setShowForm(false);
                    }}
                  />
                </div>
              )}
            </div>
            <ReminderList
              reminders={reminders}
              isLoading={isLoading}
              error={error}
              onDeleteSuccess={refetch}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
