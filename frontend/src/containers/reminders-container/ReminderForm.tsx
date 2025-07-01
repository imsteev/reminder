import React from "react";
import { useForm } from "react-hook-form";
import { useMutation, useQuery } from "@tanstack/react-query";
import { Button, Field, Input, Textarea } from "../../components/ui";
import {
  createReminder,
  getContactMethods,
  createContactMethod,
} from "../../api/reminders";
import {
  getCurrentDateTimeString,
  getDateTimeStringInMinutes,
  getTomorrowDateTimeString,
} from "../../utils/datetime";
import {
  DEFAULT_USER_ID,
  DEFAULT_START_DELAY_MINUTES,
  TIME_PRESETS,
  UI_TEXT,
} from "../../constants";

interface ReminderFormData {
  name?: string;
  body?: string;
  contactMethodID?: number;
  reminderType?: "one-time" | "repeating";
  intervalDays?: number;
  intervalHours?: number;
  intervalMinutes?: number;
  startTime?: string;
  newContactMethodType?: "phone" | "email";
  newContactMethodValue?: string;
  newContactMethodDescription?: string;
  createNewContactMethod?: boolean;
}

interface ReminderFormProps {
  onSuccess?: () => void;
  initialData?: Partial<ReminderFormData>;
}

const ReminderForm: React.FC<ReminderFormProps> = ({
  onSuccess,
  initialData,
}) => {
  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
    formState: {},
  } = useForm<ReminderFormData>({
    defaultValues: {
      name: initialData?.name || "",
      body: initialData?.body || "",
      reminderType: initialData?.reminderType || "one-time",
      contactMethodID: initialData?.contactMethodID || undefined,
      intervalDays: initialData?.intervalDays || 0,
      intervalHours: initialData?.intervalHours || 1,
      intervalMinutes: initialData?.intervalMinutes || 0,
      startTime: initialData?.startTime || "",
      createNewContactMethod: false,
      newContactMethodType: "email",
    },
  });

  const reminderType = watch("reminderType");
  const createNewContactMethod = watch("createNewContactMethod");
  const newContactMethodType = watch("newContactMethodType");

  const { data: contactMethods = [] } = useQuery({
    queryKey: ["contactMethods", DEFAULT_USER_ID],
    queryFn: () => getContactMethods({ user_id: DEFAULT_USER_ID }),
  });

  const createContactMethodMutation = useMutation({
    mutationFn: createContactMethod,
  });

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
    if (!initialData?.startTime) {
      setStartTimeIn(DEFAULT_START_DELAY_MINUTES);
    }
  }, [initialData?.startTime]);

  const onSubmit = async (data: ReminderFormData) => {
    const isOneTime = data.reminderType === "one-time";

    let periodMinutes = 0;
    if (!isOneTime) {
      const days = data.intervalDays || 0;
      const hours = data.intervalHours || 0;
      const minutes = data.intervalMinutes || 0;
      periodMinutes = days * 24 * 60 + hours * 60 + minutes;
    }

    let contactMethodId = data.contactMethodID;

    if (data.createNewContactMethod) {
      try {
        const newContactMethod = await createContactMethodMutation.mutateAsync({
          user_id: DEFAULT_USER_ID,
          type: data.newContactMethodType || "email",
          value: data.newContactMethodValue || "",
          description: data.newContactMethodDescription || "",
        });
        contactMethodId = newContactMethod.id;
      } catch (error) {
        console.error("Failed to create contact method:", error);
        return;
      }
    }

    if (!contactMethodId) {
      console.error("No contact method selected or created");
      return;
    }

    createMutation.mutate({
      user_id: DEFAULT_USER_ID,
      body: data.body || "",
      start_time: data.startTime
        ? new Date(data.startTime).toISOString()
        : new Date().toISOString(),
      is_repeating: data.reminderType === "repeating",
      period_minutes: periodMinutes,
      contact_method_id: contactMethodId,
    });
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 bg-white">
      <Field.Root>
        <Field.Control
          render={
            <Input
              {...register("name")}
              className="w-full"
              placeholder="Message"
              autoFocus
            />
          }
        />
      </Field.Root>
      <Field.Root>
        <Field.Label>Reminder Type</Field.Label>
        <div className="flex gap-3">
          <label className="flex items-center gap-2 px-3 py-2 rounded-lg border border-gray-200 hover:bg-gray-50 cursor-pointer transition-colors">
            <input
              type="radio"
              value="one-time"
              {...register("reminderType")}
              className="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
            />
            <span className="text-sm">One-time</span>
          </label>
          <label className="flex items-center gap-2 px-3 py-2 rounded-lg border border-gray-200 hover:bg-gray-50 cursor-pointer transition-colors">
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
                      onFocus={(e) => e.target.select()}
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
                      onFocus={(e) => e.target.select()}
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
                      onFocus={(e) => e.target.select()}
                    />
                  }
                />
                <Field.Description>Days</Field.Description>
              </Field.Root>
            </div>
          </div>
        </div>
      )}

      {/* <Field.Root>
        <Field.Control
          render={
            <Textarea
              {...register("body")}
              placeholder={UI_TEXT.MESSAGE_PLACEHOLDER}
              rows={3}
            />
          }
        />
      </Field.Root> */}

      <Field.Root>
        <Field.Label>Contact Method</Field.Label>
        {contactMethods.length > 0 && (
          <div className="space-y-3">
            <div>
              {!createNewContactMethod && (
                <div className="mt-2">
                  <select
                    {...register("contactMethodID", { valueAsNumber: true })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="">Select a contact method</option>
                    {contactMethods.map((method) => (
                      <option key={method.id} value={method.id}>
                        {method.description ||
                          `${method.type}: ${method.value}`}
                      </option>
                    ))}
                  </select>
                </div>
              )}
            </div>
          </div>
        )}

        {(contactMethods.length === 0 || createNewContactMethod) && (
          <div className="space-y-3 mt-3">
            {contactMethods.length === 0 && (
              <p className="text-sm text-gray-600">
                No contact methods found. Create your first one:
              </p>
            )}

            <div className="flex gap-3">
              <label className="flex items-center gap-2">
                <input
                  type="radio"
                  value="email"
                  {...register("newContactMethodType")}
                  className="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
                />
                <span className="text-sm">Email</span>
              </label>
              <label className="flex items-center gap-2">
                <input
                  type="radio"
                  value="phone"
                  {...register("newContactMethodType")}
                  className="w-4 h-4 text-blue-600 border-gray-300 focus:ring-blue-500"
                />
                <span className="text-sm">Phone</span>
              </label>
            </div>

            <Input
              {...register("newContactMethodValue")}
              placeholder={
                newContactMethodType === "email"
                  ? "user@example.com"
                  : "+1234567890"
              }
              type={newContactMethodType === "email" ? "email" : "tel"}
            />

            <Input
              {...register("newContactMethodDescription")}
              placeholder="Description (e.g., Personal email, Work phone)"
            />
          </div>
        )}
      </Field.Root>

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
