import React from "react";
import { format } from "date-fns";

interface TimelineMarkerProps {
  currentTime: Date;
}

export default function TimelineMarker({ currentTime }: TimelineMarkerProps) {
  return (
    <div className="relative flex items-center mb-6">
      <div className="flex items-center">
        <div className="w-4 h-4 bg-red-500 rounded-full border-2 border-white shadow-md z-10"></div>
        <div className="ml-4 bg-red-50 border border-red-200 rounded-lg px-4 py-2">
          <div className="flex items-center gap-2">
            <span className="text-red-600 font-semibold text-sm">NOW</span>
            <span className="text-red-500 text-xs">
              {format(currentTime, "h:mm a")}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
