import {
  SignedIn,
  SignedOut,
  SignInButton,
  UserButton,
} from "@clerk/clerk-react";
import React from "react";
import { Button } from "./ui";

interface Props {
  children: React.ReactNode;
}

export default function ClerkSignedInComponent({ children }: Props) {
  return (
    <>
      <SignedIn>{children}</SignedIn>
      <SignedOut>
        <div className="min-h-screen bg-gray-50 flex items-center justify-center">
          <div className="bg-white p-8 rounded-lg shadow-lg border border-gray-200 text-center max-w-md">
            <div className="mb-6 min-w-[300px]">
              <div className="flex items-center justify-center w-16 h-16 bg-blue-100 rounded-full mx-auto mb-4">
                <span className="text-3xl">⏰</span>
              </div>
              <h1 className="text-2xl font-bold text-gray-900 mb-2">Uchi</h1>
              <p className="text-gray-600">Manage your reminders</p>
            </div>
            <SignInButton>
              <Button variant="primary" className="w-full">
                Sign in
              </Button>
            </SignInButton>
          </div>
        </div>
      </SignedOut>
    </>
  );
}

export { UserButton };
