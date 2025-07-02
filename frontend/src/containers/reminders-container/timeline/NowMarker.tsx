import React, { useContext } from "react";
import { format } from "date-fns";
import { CurrentTimeContext } from "../../../contexts/CurrentTimeContext";

export default function NowMarker() {
  const currentTime = useContext(CurrentTimeContext);
  if (!currentTime) return null;
  return <p className="text-xs opacity-90">{format(currentTime, "h:mm a")}</p>;
}
