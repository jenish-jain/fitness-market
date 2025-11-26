import { AuthProvider } from '@/components/providers/AuthProvider'
import './globals.css'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Fitness Market',
  description: 'Your fitness marketplace',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>
        <AuthProvider>
          {children}
        </AuthProvider>
      </body>
    </html>
  )
}