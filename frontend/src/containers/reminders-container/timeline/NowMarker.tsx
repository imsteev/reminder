import React, { useContext } from "react";
import { format } from "date-fns";
import { CurrentTimeContext } from "../../../contexts/CurrentTimeContext";

export default function NowMarker() {
  const currentTime = useContext(CurrentTimeContext);
  if (!currentTime) return null;
  return (
    <div className="inline-flex items-center">
      <div className="bg-orange-400 text-white px-3 py-1 rounded-md shadow-sm">
        <div className="flex items-center gap-2">
          <div className="w-2 h-2 bg-white rounded-full animate-pulse"></div>
          <span className="text-xs opacity-90">
            {format(currentTime, "h:mm a")}
          </span>
        </div>
      </div>
    </div>
  );
}
