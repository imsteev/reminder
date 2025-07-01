import React, { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Button, Card, CardContent, Input, Field } from "./ui";
import { 
  getContactMethods, 
  createContactMethod, 
  updateContactMethod, 
  deleteContactMethod,
  ContactMethod,
  CreateContactMethodRequest,
  UpdateContactMethodRequest
} from "../api/reminders";
import { DEFAULT_USER_ID } from "../constants";

interface ContactMethodsManagerProps {
  userId?: number;
}

export default function ContactMethodsManager({ userId = DEFAULT_USER_ID }: ContactMethodsManagerProps) {
  const [editingId, setEditingId] = useState<number | null>(null);
  const [showAddForm, setShowAddForm] = useState(false);
  const [newMethod, setNewMethod] = useState<CreateContactMethodRequest>({
    user_id: userId,
    type: "email",
    value: "",
    description: "",
  });

  const queryClient = useQueryClient();

  const { data: contactMethods = [], isLoading } = useQuery({
    queryKey: ["contactMethods", userId],
    queryFn: () => getContactMethods(userId),
  });

  const createMutation = useMutation({
    mutationFn: createContactMethod,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["contactMethods", userId] });
      setShowAddForm(false);
      setNewMethod({ user_id: userId, type: "email", value: "", description: "" });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateContactMethodRequest }) =>
      updateContactMethod(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["contactMethods", userId] });
      setEditingId(null);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteContactMethod,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["contactMethods", userId] });
    },
  });

  const handleCreate = () => {
    createMutation.mutate(newMethod);
  };

  const handleUpdate = (id: number, method: ContactMethod) => {
    updateMutation.mutate({
      id,
      data: {
        type: method.type,
        value: method.value,
        description: method.description,
      },
    });
  };

  const handleDelete = (id: number) => {
    if (window.confirm("Are you sure you want to delete this contact method?")) {
      deleteMutation.mutate(id);
    }
  };

  if (isLoading) {
    return <div className="text-center py-4">Loading contact methods...</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold text-gray-900">Contact Methods</h2>
        <Button
          onClick={() => setShowAddForm(true)}
          className="bg-blue-600 text-white hover:bg-blue-700"
          size="sm"
        >
          + Add Contact Method
        </Button>
      </div>

      {showAddForm && (
        <Card>
          <CardContent className="space-y-4">
            <h3 className="font-medium">Add New Contact Method</h3>
            
            <div className="flex gap-3">
              <label className="flex items-center gap-2">
                <input
                  type="radio"
                  checked={newMethod.type === "email"}
                  onChange={() => setNewMethod({ ...newMethod, type: "email" })}
                  className="w-4 h-4 text-blue-600"
                />
                <span className="text-sm">Email</span>
              </label>
              <label className="flex items-center gap-2">
                <input
                  type="radio"
                  checked={newMethod.type === "phone"}
                  onChange={() => setNewMethod({ ...newMethod, type: "phone" })}
                  className="w-4 h-4 text-blue-600"
                />
                <span className="text-sm">Phone</span>
              </label>
            </div>

            <Input
              value={newMethod.value}
              onChange={(e) => setNewMethod({ ...newMethod, value: e.target.value })}
              placeholder={newMethod.type === "email" ? "user@example.com" : "+1234567890"}
              type={newMethod.type === "email" ? "email" : "tel"}
            />

            <Input
              value={newMethod.description}
              onChange={(e) => setNewMethod({ ...newMethod, description: e.target.value })}
              placeholder="Description (e.g., Personal email, Work phone)"
            />

            <div className="flex gap-2">
              <Button
                onClick={handleCreate}
                disabled={createMutation.isPending || !newMethod.value}
                className="bg-green-600 text-white hover:bg-green-700"
                size="sm"
              >
                {createMutation.isPending ? "Creating..." : "Create"}
              </Button>
              <Button
                onClick={() => setShowAddForm(false)}
                variant="ghost"
                size="sm"
              >
                Cancel
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      <div className="space-y-3">
        {contactMethods.map((method) => (
          <ContactMethodCard
            key={method.id}
            method={method}
            isEditing={editingId === method.id}
            onEdit={() => setEditingId(method.id)}
            onSave={(updatedMethod) => handleUpdate(method.id, updatedMethod)}
            onCancel={() => setEditingId(null)}
            onDelete={() => handleDelete(method.id)}
            isUpdating={updateMutation.isPending}
            isDeleting={deleteMutation.isPending}
          />
        ))}
      </div>

      {contactMethods.length === 0 && (
        <div className="text-center py-8 text-gray-500">
          <p>No contact methods found.</p>
          <p className="text-sm">Add your first contact method to start receiving reminders.</p>
        </div>
      )}
    </div>
  );
}

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

function ContactMethodCard({
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
                onChange={() => setEditedMethod({ ...editedMethod, type: "email" })}
                className="w-4 h-4 text-blue-600"
              />
              <span className="text-sm">Email</span>
            </label>
            <label className="flex items-center gap-2">
              <input
                type="radio"
                checked={editedMethod.type === "phone"}
                onChange={() => setEditedMethod({ ...editedMethod, type: "phone" })}
                className="w-4 h-4 text-blue-600"
              />
              <span className="text-sm">Phone</span>
            </label>
          </div>

          <Input
            value={editedMethod.value}
            onChange={(e) => setEditedMethod({ ...editedMethod, value: e.target.value })}
            type={editedMethod.type === "email" ? "email" : "tel"}
          />

          <Input
            value={editedMethod.description}
            onChange={(e) => setEditedMethod({ ...editedMethod, description: e.target.value })}
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
              <span className="font-medium">{method.value}</span>
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