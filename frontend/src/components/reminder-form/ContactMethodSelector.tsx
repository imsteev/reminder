import React from "react";
import { UseFormRegister, UseFormWatch } from "react-hook-form";
import { Link } from "react-router-dom";
import { Button, Field } from "../ui";
import { formatPhoneNumber } from "../ui/PhoneInput";
import type { ContactMethod } from "../../types/protocol";

interface ReminderFormData {
  contactMethodID?: number;
  newContactMethodType?: "phone" | "email";
}

interface Props {
  contactMethods: ContactMethod[];
  register: UseFormRegister<ReminderFormData>;
  watch: UseFormWatch<ReminderFormData>;
  onCreatePhoneMethod: () => void;
  onCreateEmailMethod: () => void;
}

export default function ContactMethodSelector({
  contactMethods,
  register,
  watch,
  onCreatePhoneMethod,
  onCreateEmailMethod,
}: Props) {
  const newContactMethodType = watch("newContactMethodType");

  return (
    <div className="space-y-3">
      <label className="flex justify-between items-center text-sm font-medium text-gray-700 rounded-lg">
        <span>Contact method</span>
        <Field.Root>
          {!!contactMethods?.length && (
            <div className="space-y-3">
              <div>
                <select
                  {...register("contactMethodID", {
                    valueAsNumber: true,
                  })}
                  value={watch("contactMethodID")}
                  className="bg-white text-sm w-full px-2 py-1 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="0">Select a contact method</option>
                  {contactMethods?.map((method) => (
                    <option key={method.id} value={method.id}>
                      {method.description}{" "}
                      {method.type === "email" && `| ${method.value}`}
                      {method.type === "phone" &&
                        `| ${formatPhoneNumber(method.value)
                          .replaceAll("(", "")
                          .replaceAll(")", "")}`}
                    </option>
                  ))}
                </select>
              </div>
            </div>
          )}

          {contactMethods?.length === 0 && (
            <div className="space-y-3 mt-3">
              <p className="text-sm text-gray-600">
                No contact methods found. Create your first one:
              </p>

              <div className="flex gap-3">
                <div className="flex gap-2">
                  <input
                    type="radio"
                    id="email"
                    value="email"
                    {...register("newContactMethodType")}
                    className="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
                  />
                  <label htmlFor="email" className="text-sm text-gray-700">
                    Email
                  </label>
                </div>
                <div className="flex gap-2">
                  <input
                    type="radio"
                    id="phone"
                    value="phone"
                    {...register("newContactMethodType")}
                    className="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
                  />
                  <label htmlFor="phone" className="text-sm text-gray-700">
                    Phone
                  </label>
                </div>
              </div>
              <div className="flex gap-2">
                <Button
                  type="button"
                  size="sm"
                  variant="secondary"
                  onClick={
                    newContactMethodType === "phone"
                      ? onCreatePhoneMethod
                      : onCreateEmailMethod
                  }
                  className="w-full"
                >
                  Create {newContactMethodType} method
                </Button>
              </div>
            </div>
          )}
        </Field.Root>
      </label>
      
      {contactMethods?.length > 0 && (
        <p className="text-xs text-gray-500">
          <Link to="/contact-methods" className="text-blue-600 hover:underline">
            Manage contact methods
          </Link>
        </p>
      )}
    </div>
  );
}