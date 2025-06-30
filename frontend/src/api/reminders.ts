import axios from "axios";
import { API_BASE_URL } from "../constants";

const API_URL = `${API_BASE_URL}/api`;

export interface Reminder {
  id: number;
  user_id: number;
  message: string;
  start_time: string;
  type: "one-time" | "repeating";
  period_minutes: number;
  delivery_type: "sms" | "email";
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface CreateReminderRequest {
  user_id: number;
  message: string;
  start_time: string;
  type: "one-time" | "repeating";
  period_minutes: number;
  delivery_type: "sms" | "email";
}

export interface ContactMethod {
  user_id: number;
  type: "phone" | "email";
  value: string;
  description: string;
}

export const getReminders = async (userId: number): Promise<Reminder[]> => {
  const response = await axios.get(`${API_URL}/reminders?user_id=${userId}`);
  return response.data;
};

export const createReminder = async (
  reminder: CreateReminderRequest,
): Promise<Reminder> => {
  const response = await axios.post(`${API_URL}/reminders`, reminder);
  return response.data;
};

export const updateReminder = async (
  id: number,
  reminder: Partial<Reminder>,
): Promise<void> => {
  await axios.put(`${API_URL}/reminders/${id}`, reminder);
};

export const deleteReminder = async (id: number): Promise<void> => {
  await axios.delete(`${API_URL}/reminders/${id}`);
};
