import React, { useCallback } from "react";
import RemindersContainer from "./containers/reminders-container/RemindersContainer";
import NavBar from "./components/NavBar";
import { useQueryClient } from "@tanstack/react-query";

// jotai ? react-query ?
function App() {
  const queryClient = useQueryClient();
  const refetchReminders = useCallback(() => {
    queryClient.invalidateQueries({ queryKey: ["reminders"] });
  }, [queryClient]);
  return (
    <div className="min-h-screen bg-gray-50">
      <NavBar refetchReminders={refetchReminders} />
      <RemindersContainer />
    </div>
  );
}

export default App;
