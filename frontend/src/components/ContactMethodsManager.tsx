import React, { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Button, Card, CardContent, Input, Field, PhoneInput } from "./ui";
import {
  getContactMethods,
  createContactMethod,
  updateContactMethod,
  deleteContactMethod,
  type ContactMethod,
  type CreateContactMethodRequest,
  type UpdateContactMethodRequest,
} from "../api/reminders";
import { DEFAULT_USER_ID } from "../constants";
import ContactMethodCard from "./ContactMethodCard";
import ContactMethodForm from "./ContactMethodForm";

interface ContactMethodsManagerProps {
  userId?: number;
}

export default function ContactMethodsManager({
  userId = DEFAULT_USER_ID,
}: ContactMethodsManagerProps) {
  const [editingId, setEditingId] = useState<number | null>(null);
  const [showAddForm, setShowAddForm] = useState(false);

  const {
    data: contactMethods = [],
    isLoading,
    refetch,
  } = useQuery({
    queryKey: ["contactMethods", userId],
    queryFn: () => getContactMethods({ user_id: userId }),
  });

  const updateMutation = useMutation({
    mutationFn: ({
      id,
      data,
    }: {
      id: number;
      data: UpdateContactMethodRequest;
    }) => updateContactMethod(id, data),
    onSuccess: () => {
      refetch();
      setEditingId(null);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteContactMethod,
    onSuccess: () => {
      refetch();
    },
  });

  const handleDelete = (id: number) => {
    if (
      window.confirm("Are you sure you want to delete this contact method?")
    ) {
      deleteMutation.mutate(id);
    }
  };

  if (isLoading) {
    return <div className="text-center py-4">Loading contact methods...</div>;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold text-gray-900">
            Contact Methods
          </h2>
          <Button
            variant="outline"
            onClick={() => setShowAddForm(true)}
            size="sm"
          >
            + Contact Method
          </Button>
        </div>

        {showAddForm && (
          <ContactMethodForm
            onCancel={() => setShowAddForm(false)}
            onSuccess={() => refetch()}
          />
        )}

        <div className="space-y-3">
          {contactMethods.map((method) => (
            <ContactMethodCard
              key={method.id}
              method={method}
              isEditing={editingId === method.id}
              onEdit={() => setEditingId(method.id)}
              onSave={(updatedMethod) =>
                updateMutation.mutate({
                  id: method.id,
                  data: updatedMethod,
                })
              }
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
            <p className="text-sm">
              Add your first contact method to start receiving reminders.
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
