import React from "react";
import { format } from "date-fns";

interface Props {
  currentTime: Date;
}

export default function NowMarker({ currentTime }: Props) {
  return (
    <div className="inline-flex items-center">
      <div className="bg-purple-400 text-white px-3 py-1 rounded-full shadow-sm">
        <div className="flex items-center gap-2">
          <div className="w-2 h-2 bg-white rounded-full animate-pulse"></div>
          <span className="font-medium text-sm">NOW</span>
          <span className="text-xs opacity-90">
            {format(currentTime, "h:mm a")}
          </span>
        </div>
      </div>
    </div>
  );
}
