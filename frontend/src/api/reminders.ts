import axios from "axios";
import { API_BASE_URL } from "../constants";
import type {
  Reminder,
  CreateReminderRequest,
  UpdateReminderRequest,
  ContactMethod,
  CreateContactMethodRequest,
  UpdateContactMethodRequest,
  DeleteResponse,
  ErrorResponse,
  GetRemindersQuery,
  GetContactMethodsQuery,
} from "../types/protocol";

const API_URL = `${API_BASE_URL}/api`;

export const getReminders = async (
  query: GetRemindersQuery
): Promise<Reminder[] | null> => {
  const response = await axios.get(
    `${API_URL}/reminders?user_id=${query.user_id}&include_past=${query.include_past}`
  );
  return response.data;
};

export const createReminder = async (
  reminder: CreateReminderRequest
): Promise<Reminder> => {
  const response = await axios.post(`${API_URL}/reminders`, reminder);
  return response.data;
};

export const updateReminder = async (
  id: number,
  reminder: UpdateReminderRequest
): Promise<Reminder> => {
  const response = await axios.put(`${API_URL}/reminders/${id}`, reminder);
  return response.data;
};

export const deleteReminder = async (id: number): Promise<DeleteResponse> => {
  const response = await axios.delete(`${API_URL}/reminders/${id}`);
  return response.data;
};

export const getContactMethods = async (
  query: GetContactMethodsQuery
): Promise<ContactMethod[] | null> => {
  const response = await axios.get(
    `${API_URL}/contact-methods?user_id=${query.user_id}`
  );
  return response.data;
};

export const createContactMethod = async (
  contactMethod: CreateContactMethodRequest
): Promise<ContactMethod> => {
  const response = await axios.post(
    `${API_URL}/contact-methods`,
    contactMethod
  );
  return response.data;
};

export const updateContactMethod = async (
  id: number,
  contactMethod: UpdateContactMethodRequest
): Promise<ContactMethod> => {
  const response = await axios.put(
    `${API_URL}/contact-methods/${id}`,
    contactMethod
  );
  return response.data;
};

export const deleteContactMethod = async (
  id: number
): Promise<DeleteResponse> => {
  const response = await axios.delete(`${API_URL}/contact-methods/${id}`);
  return response.data;
};

// Re-export types for convenience
export type {
  Reminder,
  CreateReminderRequest,
  UpdateReminderRequest,
  ContactMethod,
  CreateContactMethodRequest,
  UpdateContactMethodRequest,
  DeleteResponse,
  ErrorResponse,
  GetRemindersQuery,
  GetContactMethodsQuery,
};
