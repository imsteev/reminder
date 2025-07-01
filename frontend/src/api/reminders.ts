import axios from "axios";
import { API_BASE_URL } from "../constants";

const API_URL = `${API_BASE_URL}/api`;

export interface Reminder {
  id: number;
  user_id: number;
  name: string;
  body: string;
  start_time: string;
  is_repeating: boolean;
  period_minutes: number;
  delivery_type: "sms" | "email";
  phone_number?: string;
  email?: string;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
  river_job_id?: number;
  contact_method_id: number;
}

export interface CreateReminderRequest {
  user_id: number;
  body: string;
  start_time: string;
  is_repeating: boolean;
  period_minutes: number;
  contact_method_id: number;
  phone_number?: string;
  email?: string;
}

export interface ContactMethod {
  id: number;
  user_id: number;
  type: "phone" | "email";
  value: string;
  description: string;
}

export interface CreateContactMethodRequest {
  user_id: number;
  type: "phone" | "email";
  value: string;
  description: string;
}

export interface UpdateContactMethodRequest {
  type: "phone" | "email";
  value: string;
  description: string;
}

export const getReminders = async (
  userId: number,
  includePast = false
): Promise<Reminder[]> => {
  const response = await axios.get(
    `${API_URL}/reminders?user_id=${userId}&include_past=${includePast}`
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
  reminder: Partial<Reminder>
): Promise<void> => {
  await axios.put(`${API_URL}/reminders/${id}`, reminder);
};

export const deleteReminder = async (id: number): Promise<void> => {
  await axios.delete(`${API_URL}/reminders/${id}`);
};

export const getContactMethods = async (userId: number): Promise<ContactMethod[]> => {
  const response = await axios.get(`${API_URL}/contact-methods?user_id=${userId}`);
  return response.data;
};

export const createContactMethod = async (
  contactMethod: CreateContactMethodRequest
): Promise<ContactMethod> => {
  const response = await axios.post(`${API_URL}/contact-methods`, contactMethod);
  return response.data;
};

export const updateContactMethod = async (
  id: number,
  contactMethod: UpdateContactMethodRequest
): Promise<void> => {
  await axios.put(`${API_URL}/contact-methods/${id}`, contactMethod);
};

export const deleteContactMethod = async (id: number): Promise<void> => {
  await axios.delete(`${API_URL}/contact-methods/${id}`);
};
