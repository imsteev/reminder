import React from "react";
import { useForm } from "react-hook-form";
import { useMutation } from "@tanstack/react-query";
import { Button, Field, Input, Textarea } from "./ui";
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
  UI_TEXT,
} from "../constants";

interface ReminderFormData {
  name?: string;
  message?: string;
  contactValue?: string;
  reminderType?: "one-time" | "repeating";
  deliveryType?: "sms" | "email";
  intervalDays?: number;
  intervalHours?: number;
  intervalMinutes?: number;
  startTime?: string;
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
    watch,
    formState: {},
  } = useForm<ReminderFormData>({
    defaultValues: {
      name: "",
      reminderType: "repeating",
      deliveryType: "sms",
      intervalDays: 0,
      intervalHours: 1,
      intervalMinutes: 0,
    },
  });

  const reminderType = watch("reminderType");
  const deliveryType = watch("deliveryType");

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
    const isOneTime = data.reminderType === "one-time";

    let periodMinutes = 0;
    if (!isOneTime) {
      const days = data.intervalDays || 0;
      const hours = data.intervalHours || 0;
      const minutes = data.intervalMinutes || 0;
      periodMinutes = days * 24 * 60 + hours * 60 + minutes;
    }

    createMutation.mutate({
      user_id: DEFAULT_USER_ID,
      message: data.message || "",
      start_time: data.startTime
        ? new Date(data.startTime).toISOString()
        : new Date().toISOString(),
      type: data.reminderType || "repeating",
      period_minutes: periodMinutes,
      delivery_type: data.deliveryType || "sms",
    });
  };

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className="space-y-4 bg-white p-6 rounded-lg shadow"
    >
      <Field.Root>
        <Field.Control
          render={
            <Input
              {...register("name")}
              className="w-full"
              placeholder="Name"
            />
          }
        />
      </Field.Root>

      <Field.Root>
        <Field.Control
          render={
            <Textarea
              {...register("message")}
              placeholder={UI_TEXT.MESSAGE_PLACEHOLDER}
              rows={3}
            />
          }
        />
      </Field.Root>

      <Field.Root>
        <Field.Label>Delivery Method</Field.Label>
        <div className="flex gap-6">
          <label className="flex items-center gap-2">
            <input
              type="radio"
              value="sms"
              {...register("deliveryType")}
              className="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
            />
            <span className="text-sm">SMS</span>
          </label>
          <label className="flex items-center gap-2">
            <input
              type="radio"
              value="email"
              {...register("deliveryType")}
              className="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
            />
            <span className="text-sm">Email</span>
          </label>
        </div>
      </Field.Root>

      <Field.Root>
        <Field.Label>
          {deliveryType === "sms" ? "Phone Number" : "Email Address"}
        </Field.Label>
        <Field.Control
          render={
            <Input
              type={deliveryType === "sms" ? "tel" : "email"}
              {...register("contactValue")}
              placeholder={
                deliveryType === "sms"
                  ? UI_TEXT.PHONE_PLACEHOLDER
                  : "user@example.com"
              }
            />
          }
        />
      </Field.Root>

      <Field.Root>
        <Field.Label>Reminder Type</Field.Label>
        <div className="flex gap-6">
          <label className="flex items-center gap-2">
            <input
              type="radio"
              value="one-time"
              {...register("reminderType")}
              className="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
            />
            <span className="text-sm">One-time</span>
          </label>
          <label className="flex items-center gap-2">
            <input
              type="radio"
              value="repeating"
              {...register("reminderType")}
              className="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
            />
            <span className="text-sm">Repeating</span>
          </label>
        </div>
      </Field.Root>

      {reminderType === "repeating" && (
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Interval
            </label>
            <div className="grid grid-cols-3 gap-4">
              <Field.Root>
                <Field.Control
                  render={
                    <Input
                      type="number"
                      {...register("intervalMinutes", {
                        min: 0,
                        max: 59,
                        valueAsNumber: true,
                      })}
                      placeholder="0"
                      min={0}
                      max={59}
                    />
                  }
                />
                <Field.Description>Minutes</Field.Description>
              </Field.Root>

              <Field.Root>
                <Field.Control
                  render={
                    <Input
                      type="number"
                      {...register("intervalHours", {
                        min: 0,
                        max: 23,
                        valueAsNumber: true,
                      })}
                      placeholder="0"
                      min={0}
                      max={23}
                    />
                  }
                />
                <Field.Description>Hours</Field.Description>
              </Field.Root>

              <Field.Root>
                <Field.Control
                  render={
                    <Input
                      type="number"
                      {...register("intervalDays", {
                        min: 0,
                        valueAsNumber: true,
                      })}
                      placeholder="0"
                      min={0}
                    />
                  }
                />
                <Field.Description>Days</Field.Description>
              </Field.Root>
            </div>
          </div>
        </div>
      )}

      <div>
        <label
          htmlFor="startTime"
          className="block text-sm font-medium text-gray-700 mb-2"
        >
          {UI_TEXT.START_TIME_LABEL}
        </label>

        {/* Custom Time Input */}
        <Field.Root>
          <Field.Description>{UI_TEXT.CUSTOM_TIME_LABEL}</Field.Description>
          <Field.Control
            render={<Input type="datetime-local" {...register("startTime")} />}
          />
        </Field.Root>

        {/* Quick Time Buttons */}
        <div className="mt-3">
          <div className="flex flex-wrap gap-2">
            <Button
              type="button"
              variant="blue"
              size="sm"
              onClick={setStartTimeToNow}
              className="rounded-full"
            >
              Now
            </Button>
            <Button
              type="button"
              variant="green"
              size="sm"
              onClick={() => setStartTimeIn(TIME_PRESETS.FIFTEEN_MINUTES)}
              className="rounded-full"
            >
              In 15 min
            </Button>
            <Button
              type="button"
              variant="green"
              size="sm"
              onClick={() => setStartTimeIn(TIME_PRESETS.ONE_HOUR)}
              className="rounded-full"
            >
              In 1 hour
            </Button>
            <Button
              type="button"
              variant="purple"
              size="sm"
              onClick={() =>
                setStartTimeTomorrow(
                  TIME_PRESETS.TOMORROW_9AM.hour,
                  TIME_PRESETS.TOMORROW_9AM.minute
                )
              }
              className="rounded-full"
            >
              Tomorrow 9am
            </Button>
          </div>
        </div>
      </div>

      <Button
        type="submit"
        disabled={createMutation.isPending}
        className="w-full"
      >
        {createMutation.isPending ? "Creating..." : "Create Reminder"}
      </Button>
    </form>
  );
};

export default ReminderForm;
