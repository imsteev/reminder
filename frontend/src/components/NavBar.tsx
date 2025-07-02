import React, { useState, useEffect, useContext } from "react";
import { Link, useLocation } from "react-router-dom";
import { Button, Dialog, DialogBody, DialogHeader, DialogTitle } from "./ui";
import { ReminderForm } from "./reminder-form";
import { CurrentTimeContext } from "../contexts/CurrentTimeContext";
import NowMarker from "../containers/reminders-container/timeline/NowMarker";
import { UserButton } from "./ClerkSignedInComponent";

interface Props {
  refetchReminders: () => void;
}

export default function NavBar({ refetchReminders }: Props) {
  const [showForm, setShowForm] = useState(false);
  const location = useLocation();

  // Keyboard shortcut for 'R' to open reminder form
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      // Only trigger if not already typing in an input/textarea/contenteditable
      const target = event.target as HTMLElement;
      const isTyping =
        target.tagName === "INPUT" ||
        target.tagName === "TEXTAREA" ||
        target.contentEditable === "true";

      if (
        !isTyping &&
        (event.key === "r" || event.key === "R") &&
        !event.metaKey &&
        !event.ctrlKey &&
        !event.altKey &&
        !event.shiftKey
      ) {
        event.preventDefault();
        setShowForm(true);
      }
    };

    document.addEventListener("keydown", handleKeyDown);

    return () => {
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, []);

  return (
    <nav className="bg-white border-b border-gray-200 shadow-sm">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          {/* Logo and Brand */}
          <div className="flex items-center space-x-3">
            <div
              onClick={() => {
                window.location.href = "/";
              }}
              className="flex items-center justify-center w-10 h-10 bg-blue-100 rounded-lg cursor-pointer"
            >
              <span className="text-xl">⏰</span>
            </div>
            <div>
              <h1 className="text-xl font-bold text-gray-900">Uchi</h1>
              <NowMarker />
            </div>
          </div>

          {/* Navigation Links */}
          <div className="hidden md:flex items-center space-x-8">
            <Link
              to="/"
              className={`font-medium transition-colors ${
                location.pathname === "/"
                  ? "text-blue-600"
                  : "text-gray-600 hover:text-blue-600"
              }`}
            >
              Reminders
            </Link>
            <Link
              to="/settings"
              className={`font-medium transition-colors ${
                location.pathname === "/settings"
                  ? "text-blue-600"
                  : "text-gray-600 hover:text-blue-600"
              }`}
            >
              Settings
            </Link>
          </div>

          {/* Action Buttons */}
          <div className="flex items-center space-x-4">
            {/* New Reminder Button */}
            <Button
              onClick={() => setShowForm(true)}
              className="bg-blue-600 text-white hover:bg-blue-700 font-medium"
              size="sm"
            >
              <span className="mr-2">+</span>
              Reminder <span className="ml-1 text-xs">(R)</span>
            </Button>

            {/* User Profile */}
            <UserButton
              appearance={{
                elements: {
                  avatarBox: "w-8 h-8",
                },
              }}
            />
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
          <DialogHeader className="flex">
            <DialogTitle className="p-0 flex items-center justify-between">
              Create New Reminder
            </DialogTitle>
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
