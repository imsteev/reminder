import React, { useCallback } from "react";
import { Outlet } from "react-router-dom";
import NavBar from "./components/NavBar";
import { useQueryClient } from "@tanstack/react-query";
import { Toaster } from "sonner";

function App() {
  const queryClient = useQueryClient();
  const refetchReminders = useCallback(() => {
    queryClient.invalidateQueries({ queryKey: ["reminders"] });
  }, [queryClient]);

  return (
    <div className="min-h-screen bg-gray-50">
      <NavBar refetchReminders={refetchReminders} />
      <Outlet />
      <Toaster richColors position="bottom-right" />
    </div>
  );
}

export default App;
