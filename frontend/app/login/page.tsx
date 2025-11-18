import { LoginForm } from '@/components/auth/login-form'
import { AuthGuard } from '@/components/auth/auth-guard'
import Link from 'next/link'

export default function LoginPage() {
  return (
    <AuthGuard requireAuth={false}>
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="max-w-md w-full space-y-8 p-8 bg-white rounded-lg shadow-md">
          <div className="text-center">
            <h2 className="text-3xl font-bold text-gray-900">Sign In</h2>
            <p className="mt-2 text-sm text-gray-600">
              Welcome back to Fitness Market
            </p>
          </div>
          
          <LoginForm />
          
          <div className="text-center">
            <p className="text-sm text-gray-600">
              Don't have an account?{' '}
              <Link href="/register" className="font-medium text-blue-600 hover:text-blue-500">
                Sign up
              </Link>
            </p>
          </div>
        </div>
      </div>
    </AuthGuard>
  )
}