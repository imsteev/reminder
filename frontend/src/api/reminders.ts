import axios from 'axios'

const API_BASE_URL = 'http://localhost:8080/api'

export interface Reminder {
  id: number
  user_id: string
  message: string
  phone_number: string
  frequency: number
  interval_hours: number
  start_time: string
  end_time?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateReminderRequest {
  user_id: string
  message: string
  phone_number: string
  frequency: number
  interval_hours: number
  start_time: string
}

export const getReminders = async (userId: string): Promise<Reminder[]> => {
  const response = await axios.get(`${API_BASE_URL}/reminders?user_id=${userId}`)
  return response.data
}

export const createReminder = async (reminder: CreateReminderRequest): Promise<Reminder> => {
  const response = await axios.post(`${API_BASE_URL}/reminders`, reminder)
  return response.data
}

export const updateReminder = async (id: number, reminder: Partial<Reminder>): Promise<void> => {
  await axios.put(`${API_BASE_URL}/reminders/${id}`, reminder)
}

export const deleteReminder = async (id: number): Promise<void> => {
  await axios.delete(`${API_BASE_URL}/reminders/${id}`)
}