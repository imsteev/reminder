import React, { useCallback, useEffect, useState } from "react";
import { Outlet } from "react-router-dom";
import NavBar from "./components/NavBar";
import { useQueryClient } from "@tanstack/react-query";
import { Toaster } from "sonner";
import { CurrentTimeProvider } from "./contexts/CurrentTimeContext";

function App() {
  const queryClient = useQueryClient();
  const refetchReminders = useCallback(() => {
    queryClient.invalidateQueries({ queryKey: ["reminders"] });
  }, [queryClient]);

  return (
    <CurrentTimeProvider>
      <div className="min-h-screen bg-gray-50">
        <NavBar refetchReminders={refetchReminders} />
        <Outlet />
        <Toaster richColors position="bottom-right" />
      </div>
    </CurrentTimeProvider>
  );
}

export default App;
