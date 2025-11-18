import { RegisterForm } from '@/components/auth/register-form'
import { AuthGuard } from '@/components/auth/auth-guard'
import Link from 'next/link'

export default function RegisterPage() {
  return (
    <AuthGuard requireAuth={false}>
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="max-w-md w-full space-y-8 p-8 bg-white rounded-lg shadow-md">
          <div className="text-center">
            <h2 className="text-3xl font-bold text-gray-900">Create Account</h2>
            <p className="mt-2 text-sm text-gray-600">
              Join Fitness Market today
            </p>
          </div>
          
          <RegisterForm />
          
          <div className="text-center">
            <p className="text-sm text-gray-600">
              Already have an account?{' '}
              <Link href="/login" className="font-medium text-blue-600 hover:text-blue-500">
                Sign in
              </Link>
            </p>
          </div>
        </div>
      </div>
    </AuthGuard>
  )
}