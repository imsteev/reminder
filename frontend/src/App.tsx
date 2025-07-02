import React, { useCallback, useEffect, useState } from "react";
import { Outlet } from "react-router-dom";
import NavBar from "./components/NavBar";
import { QueryClientProvider, useQueryClient } from "@tanstack/react-query";
import { Toaster } from "sonner";
import { CurrentTimeProvider } from "./contexts/CurrentTimeContext";
import ClerkSignedInComponent from "./components/ClerkSignedInComponent";

function App() {
  const queryClient = useQueryClient();
  const refetchReminders = useCallback(() => {
    queryClient.invalidateQueries({ queryKey: ["reminders"] });
  }, [queryClient]);

  return (
    <QueryClientProvider client={queryClient}>
      <CurrentTimeProvider>
        <ClerkSignedInComponent>
          <div className="min-h-screen bg-gray-50">
            <NavBar refetchReminders={refetchReminders} />
            <Outlet />
            <Toaster richColors position="bottom-right" />
          </div>
        </ClerkSignedInComponent>
      </CurrentTimeProvider>
    </QueryClientProvider>
  );
}

export default App;
