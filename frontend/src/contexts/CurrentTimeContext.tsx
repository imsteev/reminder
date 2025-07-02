import React, { createContext, useEffect, useState } from "react";

export const CurrentTimeContext = createContext<Date | null>(null);

export const CurrentTimeProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const [currentTime, setCurrentTime] = useState(new Date());
  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date());
    }, 30000); // Update every 30 seconds
    return () => clearInterval(timer);
  }, []);
  return (
    <CurrentTimeContext.Provider value={currentTime}>
      {children}
    </CurrentTimeContext.Provider>
  );
};
