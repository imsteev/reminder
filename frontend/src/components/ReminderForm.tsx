import React from "react";
import { useForm } from "react-hook-form";
import { useMutation } from "@tanstack/react-query";
import { createReminder } from "../api/reminders";

interface ReminderFormData {
  message: string;
  phoneNumber: string;
  frequency: number;
  intervalHours: number;
  startTime: string;
}

interface ReminderFormProps {
  onSuccess?: () => void;
}

const ReminderForm: React.FC<ReminderFormProps> = ({ onSuccess }) => {
  const {
    register,
    handleSubmit,
    reset,
    setValue,
    formState: { errors },
  } = useForm<ReminderFormData>();

  const createMutation = useMutation({
    mutationFn: createReminder,
    onSuccess: () => {
      reset();
      onSuccess?.();
    },
  });

  const setStartTimeToNow = () => {
    const now = new Date();
    const localDateTime = new Date(
      now.getTime() - now.getTimezoneOffset() * 60000
    )
      .toISOString()
      .slice(0, 16);
    setValue("startTime", localDateTime);
  };

  const onSubmit = (data: ReminderFormData) => {
    createMutation.mutate({
      user_id: "default-user",
      message: data.message,
      phone_number: data.phoneNumber,
      frequency: data.frequency,
      interval_hours: data.intervalHours,
      start_time: new Date(data.startTime).toISOString(),
    });
  };

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className="space-y-4 bg-white p-6 rounded-lg shadow"
    >
      <div>
        <label
          htmlFor="message"
          className="block text-sm font-medium text-gray-700 mb-1"
        >
          Message
        </label>
        <textarea
          id="message"
          {...register("message", { required: "Message is required" })}
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="Take your medication"
          rows={3}
        />
        {errors.message && (
          <p className="text-red-500 text-sm mt-1">{errors.message.message}</p>
        )}
      </div>

      <div>
        <label
          htmlFor="phoneNumber"
          className="block text-sm font-medium text-gray-700 mb-1"
        >
          Phone Number
        </label>
        <input
          type="tel"
          id="phoneNumber"
          {...register("phoneNumber", { required: "Phone number is required" })}
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="+1234567890"
        />
        {errors.phoneNumber && (
          <p className="text-red-500 text-sm mt-1">
            {errors.phoneNumber.message}
          </p>
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label
            htmlFor="frequency"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Times per day
          </label>
          <input
            type="number"
            id="frequency"
            {...register("frequency", {
              required: "Frequency is required",
              min: 1,
              valueAsNumber: true,
            })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="3"
            min="1"
          />
          {errors.frequency && (
            <p className="text-red-500 text-sm mt-1">
              {errors.frequency.message}
            </p>
          )}
        </div>

        <div>
          <label
            htmlFor="intervalHours"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Interval (hours)
          </label>
          <input
            type="number"
            id="intervalHours"
            {...register("intervalHours", {
              required: "Interval is required",
              min: 1,
              valueAsNumber: true,
            })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="8"
            min="1"
          />
          {errors.intervalHours && (
            <p className="text-red-500 text-sm mt-1">
              {errors.intervalHours.message}
            </p>
          )}
        </div>
      </div>

      <div>
        <label
          htmlFor="startTime"
          className="block text-sm font-medium text-gray-700 mb-1"
        >
          Start Time
        </label>
        <div className="flex gap-2">
          <input
            type="datetime-local"
            id="startTime"
            {...register("startTime", { required: "Start time is required" })}
            className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            type="button"
            onClick={setStartTimeToNow}
            className="px-3 py-2 bg-gray-100 text-gray-700 border border-gray-300 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            Now
          </button>
        </div>
        {errors.startTime && (
          <p className="text-red-500 text-sm mt-1">
            {errors.startTime.message}
          </p>
        )}
      </div>

      <button
        type="submit"
        disabled={createMutation.isPending}
        className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
      >
        {createMutation.isPending ? "Creating..." : "Create Reminder"}
      </button>
    </form>
  );
};

export default ReminderForm;
