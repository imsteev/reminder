import React, { useState } from "react";
import { Button, Dialog, DialogBody, DialogHeader, DialogTitle } from "./ui";
import ReminderForm from "../containers/reminders-container/ReminderForm";

interface Props {
  refetchReminders: () => void;
}

export default function NavBar({ refetchReminders }: Props) {
  const [showForm, setShowForm] = useState(false);
  return (
    <nav className="bg-white border-b border-gray-200 shadow-sm">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          {/* Logo and Brand */}
          <div className="flex items-center space-x-3">
            <div className="flex items-center justify-center w-10 h-10 bg-blue-100 rounded-lg">
              <span className="text-xl">⏰</span>
            </div>
            <div>
              <h1 className="text-xl font-bold text-gray-900">RemindMe</h1>
              <p className="text-xs text-gray-500">Stay organized</p>
            </div>
          </div>

          {/* Navigation Links */}
          <div className="hidden md:flex items-center space-x-8"></div>

          {/* Action Buttons */}
          <div className="flex items-center space-x-4">
            {/* New Reminder Button */}
            <Button
              onClick={() => setShowForm(true)}
              className="bg-blue-600 text-white hover:bg-blue-700 font-medium"
              size="sm"
            >
              <span className="mr-2">+</span>
              New Reminder
            </Button>

            {/* User Profile */}
            <div className="flex items-center space-x-2">
              <div className="w-8 h-8 bg-gray-200 rounded-full flex items-center justify-center text-gray-600 font-medium text-sm">
                U
              </div>
            </div>
          </div>
        </div>
      </div>
      {showForm && (
        <Dialog
          isOpen={showForm}
          onClose={() => {
            setShowForm(false);
          }}
        >
          <DialogHeader>
            <DialogTitle>Create New Reminder</DialogTitle>
          </DialogHeader>
          <DialogBody>
            <ReminderForm
              onSuccess={() => {
                setShowForm(false);
                refetchReminders();
              }}
            />
          </DialogBody>
        </Dialog>
      )}
    </nav>
  );
}
