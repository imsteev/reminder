import React, { useState } from "react";
import { Card, CardContent, PhoneInput, Input, Button } from "./ui";
import { toast } from "sonner";
import {
  type CreateContactMethodRequest,
  createContactMethod,
} from "../api/reminders";
import { useMutation } from "@tanstack/react-query";
import { DEFAULT_USER_ID } from "../constants";

interface Props {
  onSuccess: () => void;
  onCancel: () => void;
}
export default function ContactMethodForm({ onSuccess, onCancel }: Props) {
  const [formData, setFormData] = useState<CreateContactMethodRequest>({
    user_id: DEFAULT_USER_ID,
    type: "email",
    value: "",
    description: "",
  });

  const createMutation = useMutation({
    mutationFn: createContactMethod,
    onSuccess: () => {
      toast.success("Contact method created");
      setFormData({
        user_id: DEFAULT_USER_ID,
        type: "email",
        value: "",
        description: "",
      });
      onSuccess();
    },
    onError: (error) => {
      toast.error("Failed to create contact method", {
        description: error.message,
      });
    },
  });

  return (
    <Card>
      <CardContent className="space-y-4">
        <h3 className="font-medium">New Contact Method</h3>

        <div className="flex gap-3">
          <label className="flex items-center gap-2">
            <input
              type="radio"
              checked={formData.type === "email"}
              onChange={() => setFormData({ ...formData, type: "email" })}
              className="w-4 h-4 text-blue-600"
            />
            <span className="text-sm">Email</span>
          </label>
          <label className="flex items-center gap-2">
            <input
              type="radio"
              checked={formData.type === "phone"}
              onChange={() => setFormData({ ...formData, type: "phone" })}
              className="w-4 h-4 text-blue-600"
            />
            <span className="text-sm">Phone</span>
          </label>
        </div>

        {formData.type === "phone" && (
          <PhoneInput
            value={formData.value}
            onChange={(value) => setFormData({ ...formData, value })}
          />
        )}

        {formData.type === "email" && (
          <>
            <label className="text-sm text-gray-700">Email Address</label>
            <Input
              value={formData.value}
              onChange={(e) =>
                setFormData({ ...formData, value: e.target.value })
              }
              placeholder="user@example.com"
            />
          </>
        )}
        <>
          <label className="text-sm text-gray-700">Description</label>
          <Input
            value={formData.description}
            onChange={(e) =>
              setFormData({ ...formData, description: e.target.value })
            }
            placeholder="(e.g. Personal email, Work phone)"
          />
        </>

        <div className="flex gap-2">
          <Button
            onClick={() => createMutation.mutate(formData)}
            disabled={createMutation.isPending || !formData.value}
            className="bg-green-600 text-white hover:bg-green-700"
            size="sm"
          >
            {createMutation.isPending ? "Creating..." : "Create"}
          </Button>
          <Button onClick={() => onCancel()} variant="ghost" size="sm">
            Cancel
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
