import React from "react";
import { createBrowserRouter } from "react-router-dom";
import App from "./App";
import RemindersContainer from "./containers/reminders-container/RemindersContainer";
import ContactMethodsManager from "./components/ContactMethodsManager";

// Refactor this to be declarative
export const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    children: [
      {
        index: true,
        element: <RemindersContainer />,
      },
      {
        path: "/settings",
        element: <ContactMethodsManager />,
      },
    ],
  },
]);
