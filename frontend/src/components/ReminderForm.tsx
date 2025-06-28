import React from "react";
import { useForm } from "react-hook-form";
import { useMutation } from "@tanstack/react-query";
import { createReminder } from "../api/reminders";
import {
  getCurrentDateTimeString,
  getDateTimeStringInMinutes,
  getTomorrowDateTimeString,
} from "../utils/datetime";
import {
  DEFAULT_USER_ID,
  DEFAULT_START_DELAY_MINUTES,
  TIME_PRESETS,
  FORM_DEFAULTS,
  UI_TEXT,
} from "../constants";

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
    setValue("startTime", getCurrentDateTimeString());
  };

  const setStartTimeIn = (minutes: number) => {
    setValue("startTime", getDateTimeStringInMinutes(minutes));
  };

  const setStartTimeTomorrow = (hour: number, minute = 0) => {
    setValue("startTime", getTomorrowDateTimeString(hour, minute));
  };

  // Set default start time
  React.useEffect(() => {
    setStartTimeIn(DEFAULT_START_DELAY_MINUTES);
  }, []);

  const onSubmit = (data: ReminderFormData) => {
    createMutation.mutate({
      user_id: DEFAULT_USER_ID,
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
          placeholder={UI_TEXT.MESSAGE_PLACEHOLDER}
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
          placeholder={UI_TEXT.PHONE_PLACEHOLDER}
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
            placeholder={FORM_DEFAULTS.FREQUENCY_PLACEHOLDER}
            min={FORM_DEFAULTS.MIN_FREQUENCY}
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
            placeholder={FORM_DEFAULTS.INTERVAL_PLACEHOLDER}
            min={FORM_DEFAULTS.MIN_INTERVAL}
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
          className="block text-sm font-medium text-gray-700 mb-2"
        >
          {UI_TEXT.START_TIME_LABEL}
        </label>
        
        {/* Quick Time Buttons */}
        <div className="mb-3">
          <p className="text-xs text-gray-500 mb-2">{UI_TEXT.QUICK_SELECT_LABEL}</p>
          <div className="flex flex-wrap gap-2">
            <button
              type="button"
              onClick={setStartTimeToNow}
              className="px-3 py-1.5 text-sm bg-blue-100 text-blue-700 border border-blue-200 rounded-full hover:bg-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              Now
            </button>
            <button
              type="button"
              onClick={() => setStartTimeIn(TIME_PRESETS.FIFTEEN_MINUTES)}
              className="px-3 py-1.5 text-sm bg-green-100 text-green-700 border border-green-200 rounded-full hover:bg-green-200 focus:outline-none focus:ring-2 focus:ring-green-500"
            >
              In 15 min
            </button>
            <button
              type="button"
              onClick={() => setStartTimeIn(TIME_PRESETS.ONE_HOUR)}
              className="px-3 py-1.5 text-sm bg-green-100 text-green-700 border border-green-200 rounded-full hover:bg-green-200 focus:outline-none focus:ring-2 focus:ring-green-500"
            >
              In 1 hour
            </button>
            <button
              type="button"
              onClick={() => setStartTimeTomorrow(TIME_PRESETS.TOMORROW_9AM.hour, TIME_PRESETS.TOMORROW_9AM.minute)}
              className="px-3 py-1.5 text-sm bg-purple-100 text-purple-700 border border-purple-200 rounded-full hover:bg-purple-200 focus:outline-none focus:ring-2 focus:ring-purple-500"
            >
              Tomorrow 9am
            </button>
            <button
              type="button"
              onClick={() => setStartTimeTomorrow(TIME_PRESETS.TOMORROW_1PM.hour, TIME_PRESETS.TOMORROW_1PM.minute)}
              className="px-3 py-1.5 text-sm bg-purple-100 text-purple-700 border border-purple-200 rounded-full hover:bg-purple-200 focus:outline-none focus:ring-2 focus:ring-purple-500"
            >
              Tomorrow 1pm
            </button>
          </div>
        </div>

        {/* Custom Time Input */}
        <div>
          <p className="text-xs text-gray-500 mb-2">{UI_TEXT.CUSTOM_TIME_LABEL}</p>
          <input
            type="datetime-local"
            id="startTime"
            {...register("startTime", { required: "Start time is required" })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
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
