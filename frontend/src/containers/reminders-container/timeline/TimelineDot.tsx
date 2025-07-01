import React from "react";
import { cn } from "../../../utils/cn";

interface Props {
  isUpcoming?: boolean;
  type?: "repeating" | "one-time";
  className?: string;
}

export default function TimelineDot({ isUpcoming, type, className }: Props) {
  const typeColor = isUpcoming
    ? type === "repeating"
      ? "bg-blue-500"
      : "bg-green-500"
    : "";
  return (
    <div
      className={cn(
        "w-4 h-4 rounded-full border-2 border-white shadow-md z-10 flex-shrink-0",
        typeColor,
        className
      )}
    />
  );
}
