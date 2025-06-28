import React from 'react'
import { useMutation } from '@tanstack/react-query'
import { deleteReminder, Reminder } from '../api/reminders'
import { format } from 'date-fns'

interface ReminderListProps {
  reminders?: Reminder[]
  isLoading: boolean
  error: any
  onDelete: () => void
}

const ReminderList: React.FC<ReminderListProps> = ({ reminders, isLoading, error, onDelete }) => {
  const deleteMutation = useMutation({
    mutationFn: deleteReminder,
    onSuccess: () => {
      onDelete()
    },
  })

  if (isLoading) return <div className="text-center py-4">Loading reminders...</div>
  if (error) return <div className="text-red-500 text-center py-4">Error loading reminders</div>

  return (
    <div className="space-y-4">
      {reminders?.length === 0 ? (
        <p className="text-gray-500 text-center py-8">No reminders yet. Create your first one!</p>
      ) : (
        reminders?.map((reminder) => (
          <div key={reminder.id} className="bg-white p-4 rounded-lg shadow border">
            <div className="flex justify-between items-start">
              <div className="flex-1">
                <p className="text-gray-900 font-medium mb-2">{reminder.message}</p>
                <div className="text-sm text-gray-600 space-y-1">
                  <p>📱 {reminder.phone_number}</p>
                  <p>🔄 {reminder.frequency}x daily, every {reminder.interval_hours} hours</p>
                  <p>⏰ Starts: {format(new Date(reminder.start_time), 'MMM d, yyyy h:mm a')}</p>
                  <p className="flex items-center">
                    <span className={`inline-block w-2 h-2 rounded-full mr-2 ${reminder.is_active ? 'bg-green-500' : 'bg-gray-400'}`}></span>
                    {reminder.is_active ? 'Active' : 'Inactive'}
                  </p>
                </div>
              </div>
              <button
                onClick={() => deleteMutation.mutate(reminder.id)}
                disabled={deleteMutation.isPending}
                className="ml-4 text-red-600 hover:text-red-800 disabled:opacity-50"
              >
                🗑️
              </button>
            </div>
          </div>
        ))
      )}
    </div>
  )
}

export default ReminderList