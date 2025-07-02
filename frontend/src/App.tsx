import React, { useCallback, useEffect } from "react";
import { Outlet } from "react-router-dom";
import NavBar from "./components/NavBar";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "sonner";
import { CurrentTimeProvider } from "./contexts/CurrentTimeContext";
import ClerkSignedInComponent from "./components/ClerkSignedInComponent";
import { initAxios } from "./api";

const queryClient = new QueryClient();
initAxios();

function App() {
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
