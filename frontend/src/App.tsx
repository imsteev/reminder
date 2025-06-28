import React from 'react'
import { useQuery } from '@tanstack/react-query'
import { getReminders } from './api/reminders'
import ReminderList from './components/ReminderList'
import ReminderForm from './components/ReminderForm'
import { DEFAULT_USER_ID } from './constants'

function App() {
  const { data: reminders, isLoading, error, refetch } = useQuery({
    queryKey: ['reminders'],
    queryFn: () => getReminders(DEFAULT_USER_ID),
  })

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-8">Reminder Manager</h1>
        
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <div>
            <h2 className="text-xl font-semibold text-gray-800 mb-4">Create New Reminder</h2>
            <ReminderForm onSuccess={refetch} />
          </div>
          
          <div>
            <h2 className="text-xl font-semibold text-gray-800 mb-4">Your Reminders</h2>
            <ReminderList 
              reminders={reminders} 
              isLoading={isLoading} 
              error={error} 
              onDelete={refetch} 
            />
          </div>
        </div>
      </div>
    </div>
  )
}

export default App