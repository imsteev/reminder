import { useQuery } from "@tanstack/react-query";
import React, { useState } from "react";
import { getReminders } from "../../api/reminders";
import ReminderList from "./ReminderList";
import { DEFAULT_USER_ID } from "../../constants";

export default function RemindersContainer() {
  const [includePast] = useState(true);

  const {
    data: reminders,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ["reminders", includePast],
    queryFn: () =>
      getReminders({ user_id: DEFAULT_USER_ID, include_past: includePast }),
  });

  return (
    <div className="container mx-auto px-4 py-8 space-y-6">
      <h3 className="text-xl font-bold">Reminders</h3>
      <ReminderList
        reminders={reminders || []}
        isLoading={isLoading}
        error={error}
        refetch={refetch}
      />
    </div>
  );
}
