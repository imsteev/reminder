import React from "react";
import { Field } from "../ui";

interface Props {
  value: "one-time" | "repeating";
  onChange: (value: "one-time" | "repeating") => void;
}

export default function ReminderTypeSelector({ value, onChange }: Props) {
  const handleToggle = () => {
    onChange(value === "repeating" ? "one-time" : "repeating");
  };

  return (
    <Field.Root>
      <Field.Label className="flex items-center justify-between mb-2 cursor-pointer">
        When should we start sending reminders?
        <div className="text-xs text-gray-500 flex items-center cursor-pointer">
          Repeating
          <input
            type="checkbox"
            checked={value === "repeating"}
            onChange={handleToggle}
            className="ml-2 w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
          />
        </div>
      </Field.Label>
    </Field.Root>
  );
}