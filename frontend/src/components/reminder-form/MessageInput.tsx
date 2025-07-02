import React from "react";
import { UseFormRegister } from "react-hook-form";
import { Field, Textarea } from "../ui";

interface ReminderFormData {
  body?: string;
}

interface Props {
  register: UseFormRegister<ReminderFormData>;
}

export default function MessageInput({ register }: Props) {
  return (
    <Field.Root>
      <Field.Control
        render={
          <Textarea
            {...register("body")}
            className="w-full"
            placeholder="Message"
            rows={3}
          />
        }
      />
    </Field.Root>
  );
}