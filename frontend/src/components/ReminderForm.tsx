import React from "react";
import { useForm } from "react-hook-form";
import { useMutation } from "@tanstack/react-query";
import { Button, Field, Input, Select, Textarea } from "./ui";
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
  message?: string;
  phoneNumber?: string;
  cadenceType?: "frequency" | "interval";
  frequency?: number;
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
      cadenceType: "frequency",
    },
  });

  const cadenceType = watch("cadenceType");

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
    let frequency = 1;
    let intervalHours = 1;

    if (data.cadenceType === "frequency") {
      frequency = data.frequency || 1;
      intervalHours = 24 / frequency;
    } else if (data.cadenceType === "interval") {
      if (data.intervalMinutes) {
        intervalHours = data.intervalMinutes / 60;
        frequency = Math.ceil(24 / intervalHours);
      } else {
        intervalHours = data.intervalHours || 1;
        frequency = Math.ceil(24 / intervalHours);
      }
    }

    createMutation.mutate({
      user_id: DEFAULT_USER_ID,
      message: data.message || "",
      phone_number: data.phoneNumber || "",
      frequency,
      interval_hours: intervalHours,
      start_time: data.startTime
        ? new Date(data.startTime).toISOString()
        : new Date().toISOString(),
    });
  };

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className="space-y-4 bg-white p-6 rounded-lg shadow"
    >
      <Field.Root>
        <Field.Label>Message</Field.Label>
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
        <Field.Label>Phone Number</Field.Label>
        <Field.Control
          render={
            <Input
              type="tel"
              {...register("phoneNumber")}
              placeholder={UI_TEXT.PHONE_PLACEHOLDER}
            />
          }
        />
      </Field.Root>

      <div className="space-y-4">
        <Field.Root>
          <Field.Label>Every</Field.Label>
          <Field.Control
            render={
              <Select {...register("cadenceType")}>
                <option value="day">day</option>
                <option value="hour">hour</option>
                <option value="minute">minute</option>
              </Select>
            }
          />
        </Field.Root>

        {cadenceType === "frequency" && (
          <Field.Root>
            <Field.Label>Times per day</Field.Label>
            <Field.Control
              render={
                <Input
                  type="number"
                  {...register("frequency", {
                    min: 1,
                    valueAsNumber: true,
                  })}
                  placeholder={FORM_DEFAULTS.FREQUENCY_PLACEHOLDER}
                  min={FORM_DEFAULTS.MIN_FREQUENCY}
                />
              }
            />
          </Field.Root>
        )}

        {cadenceType === "interval" && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Field.Root>
              <Field.Label>Every X hours</Field.Label>
              <Field.Control
                render={
                  <Input
                    type="number"
                    {...register("intervalHours", {
                      min: 1,
                      valueAsNumber: true,
                    })}
                    placeholder={FORM_DEFAULTS.INTERVAL_PLACEHOLDER}
                    min={FORM_DEFAULTS.MIN_INTERVAL}
                  />
                }
              />
            </Field.Root>

            <Field.Root>
              <Field.Label>OR every X minutes</Field.Label>
              <Field.Control
                render={
                  <Input
                    type="number"
                    {...register("intervalMinutes", {
                      min: 1,
                      valueAsNumber: true,
                    })}
                    placeholder="30"
                    min={1}
                  />
                }
              />
            </Field.Root>
          </div>
        )}
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
          <p className="text-xs text-gray-500 mb-2">
            {UI_TEXT.QUICK_SELECT_LABEL}
          </p>
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
            <Button
              type="button"
              variant="purple"
              size="sm"
              onClick={() =>
                setStartTimeTomorrow(
                  TIME_PRESETS.TOMORROW_1PM.hour,
                  TIME_PRESETS.TOMORROW_1PM.minute
                )
              }
              className="rounded-full"
            >
              Tomorrow 1pm
            </Button>
          </div>
        </div>

        {/* Custom Time Input */}
        <Field.Root>
          <Field.Description>{UI_TEXT.CUSTOM_TIME_LABEL}</Field.Description>
          <Field.Control
            render={<Input type="datetime-local" {...register("startTime")} />}
          />
        </Field.Root>
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
