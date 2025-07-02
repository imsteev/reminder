import { SignedIn, SignedOut, SignInButton, UserButton } from "@clerk/clerk-react";
import React from "react";

interface Props {
  children: React.ReactNode;
}

export default function ClerkSignedInComponent({ children }: Props) {
  return (
    <>
      <SignedIn>
        {children}
      </SignedIn>
      <SignedOut>
        <div className="min-h-screen bg-gray-50 flex items-center justify-center">
          <div className="bg-white p-8 rounded-lg shadow-lg text-center max-w-md">
            <div className="mb-6">
              <div className="flex items-center justify-center w-16 h-16 bg-blue-100 rounded-full mx-auto mb-4">
                <span className="text-3xl">⏰</span>
              </div>
              <h1 className="text-2xl font-bold text-gray-900 mb-2">Welcome to Uchi</h1>
              <p className="text-gray-600">
                Sign in to manage your reminders and stay organized
              </p>
            </div>
            <SignInButton mode="modal">
              <button className="w-full bg-blue-600 text-white px-6 py-3 rounded-lg font-medium hover:bg-blue-700 transition-colors">
                Sign In
              </button>
            </SignInButton>
          </div>
        </div>
      </SignedOut>
    </>
  );
}

export { UserButton };