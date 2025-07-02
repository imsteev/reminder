import React from "react";
import { Field } from "../ui";

interface Props {
  value: boolean;
  onChange: (value: boolean) => void;
}

export default function ReminderTypeSelector({ value, onChange }: Props) {
  const handleToggle = () => {
    onChange(!value);
  };

  return (
    <Field.Root>
      <Field.Label className="flex items-center justify-between mb-2 cursor-pointer">
        When should we start sending reminders?
        <div className="text-xs text-gray-500 flex items-center cursor-pointer">
          Repeating
          <input
            type="checkbox"
            checked={value}
            onChange={handleToggle}
            className="ml-2 w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
          />
        </div>
      </Field.Label>
    </Field.Root>
  );
}
