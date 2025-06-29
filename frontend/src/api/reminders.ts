import axios from "axios";
import { API_BASE_URL } from "../constants";

const API_URL = `${API_BASE_URL}/api`;

export interface Reminder {
  id: number;
  user_id: string;
  message: string;
  phone_number: string;
  frequency: number;
  interval_hours: number;
  start_time: string;
  end_time?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateReminderRequest {
  user_id: string;
  message: string;
  phone_number: string;
  frequency: number;
  interval_hours: number;
  start_time: string;
}

export const getReminders = async (userId: string): Promise<Reminder[]> => {
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
