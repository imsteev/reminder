import { useState } from "react";
import { type ContactMethod } from "../api/reminders";
import React from "react";
import { Card, CardContent, PhoneInput, Input, Button } from "./ui";
import { formatPhoneNumber } from "./ui/PhoneInput";

interface ContactMethodCardProps {
  method: ContactMethod;
  isEditing: boolean;
  onEdit: () => void;
  onSave: (method: ContactMethod) => void;
  onCancel: () => void;
  onDelete: () => void;
  isUpdating: boolean;
  isDeleting: boolean;
}

export default function ContactMethodCard({
  method,
  isEditing,
  onEdit,
  onSave,
  onCancel,
  onDelete,
  isUpdating,
  isDeleting,
}: ContactMethodCardProps) {
  const [editedMethod, setEditedMethod] = useState<ContactMethod>(method);

  React.useEffect(() => {
    setEditedMethod(method);
  }, [method]);

  if (isEditing) {
    return (
      <Card>
        <CardContent className="space-y-3">
          <div className="flex gap-3">
            <label className="flex items-center gap-2">
              <input
                type="radio"
                checked={editedMethod.type === "email"}
                onChange={() =>
                  setEditedMethod({ ...editedMethod, type: "email" })
                }
                className="w-4 h-4 text-blue-600"
              />
              <span className="text-sm">Email</span>
            </label>
            <label className="flex items-center gap-2">
              <input
                type="radio"
                checked={editedMethod.type === "phone"}
                onChange={() =>
                  setEditedMethod({ ...editedMethod, type: "phone" })
                }
                className="w-4 h-4 text-blue-600"
              />
              <span className="text-sm">Phone</span>
            </label>
          </div>

          {editedMethod.type === "phone" && (
            <PhoneInput
              value={editedMethod.value}
              onChange={(value) => setEditedMethod({ ...editedMethod, value })}
            />
          )}

          {editedMethod.type === "email" && (
            <Input
              value={editedMethod.value}
              onChange={(e) =>
                setEditedMethod({ ...editedMethod, value: e.target.value })
              }
              type={editedMethod.type === "email" ? "email" : "tel"}
            />
          )}

          <Input
            value={editedMethod.description}
            onChange={(e) =>
              setEditedMethod({ ...editedMethod, description: e.target.value })
            }
            placeholder="Description"
          />

          <div className="flex gap-2">
            <Button
              onClick={() => onSave(editedMethod)}
              disabled={isUpdating || !editedMethod.value}
              className="bg-green-600 text-white hover:bg-green-700"
              size="sm"
            >
              {isUpdating ? "Saving..." : "Save"}
            </Button>
            <Button onClick={onCancel} variant="ghost" size="sm">
              Cancel
            </Button>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardContent>
        <div className="flex items-center justify-between">
          <div>
            <div className="flex items-center gap-2">
              <span className="px-2 py-1 bg-gray-100 text-gray-700 rounded text-xs font-medium">
                {method.type.toUpperCase()}
              </span>
              <span className="font-medium">
                {method.type === "phone"
                  ? formatPhoneNumber(method.value)
                  : method.value}
              </span>
            </div>
            {method.description && (
              <p className="text-sm text-gray-600 mt-1">{method.description}</p>
            )}
          </div>
          <div className="flex gap-2">
            <Button onClick={onEdit} variant="ghost" size="sm">
              Edit
            </Button>
            <Button
              onClick={onDelete}
              variant="ghost"
              size="sm"
              disabled={isDeleting}
              className="text-red-600 hover:text-red-800"
            >
              {isDeleting ? "..." : "Delete"}
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
