import { useQuery } from "@tanstack/react-query";
import React, { useState } from "react";
import { getReminders } from "../../api/reminders";
import ReminderForm from "./ReminderForm";
import ReminderList from "./ReminderList";
import { Button } from "../../components/ui";
import { DEFAULT_USER_ID } from "../../constants";

export default function RemindersContainer() {
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
      <div className="container mx-auto px-4 py-8 space-y-6">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">Reminders</h1>
        <ReminderList
          reminders={reminders}
          isLoading={isLoading}
          error={error}
          refetch={refetch}
        />
      </div>
    </div>
  );
}
