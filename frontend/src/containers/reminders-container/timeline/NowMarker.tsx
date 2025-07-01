import React from "react";
import { format } from "date-fns";
import TimelineDot from "./TimelineDot";

interface Props {
  currentTime: Date;
}

export default function NowMarker({ currentTime }: Props) {
  return (
    <div className="flex items-center">
      <div className=" bg-red-50 border border-red-200 rounded-lg px-4 py-2">
        <div className="flex items-center gap-2">
          <span className="text-red-600 font-semibold text-sm">NOW</span>
          <span className="text-red-500 text-xs">
            {format(currentTime, "h:mm a")}
          </span>
        </div>
      </div>
    </div>
  );
}
