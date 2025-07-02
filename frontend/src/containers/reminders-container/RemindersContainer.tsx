import { useQuery } from "@tanstack/react-query";
import React, { useState } from "react";
import { getReminders } from "../../api/reminders";
import ReminderList from "./ReminderList";
import { DEFAULT_USER_ID } from "../../constants";
import NowMarker from "./timeline/NowMarker";

export default function RemindersContainer() {
  const [includePast, setIncludePast] = useState(false);

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
      <div className="flex items-center justify-between">
        <div className="flex items-center justify-between gap-4">
          <NowMarker />
        </div>
        <div className="flex items-center space-x-3">
          <label
            htmlFor="include-past"
            className="text-sm font-medium text-gray-700"
          >
            Show past reminders
          </label>
          <input
            id="include-past"
            type="checkbox"
            checked={includePast}
            onChange={(e) => setIncludePast(e.target.checked)}
            className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
          />
        </div>
      </div>
      <ReminderList
        reminders={reminders || []}
        isLoading={isLoading}
        error={error}
        refetch={refetch}
      />
    </div>
  );
}
