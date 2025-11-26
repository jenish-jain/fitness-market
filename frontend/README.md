# Fitness Market Frontend

## Setup

1. Install dependencies:
```bash
npm install
```

2. Copy environment variables:
```bash
cp .env.local.example .env.local
```

3. Update `.env.local` with your Supabase credentials:
- `NEXT_PUBLIC_SUPABASE_URL`: Your Supabase project URL
- `NEXT_PUBLIC_SUPABASE_ANON_KEY`: Your Supabase anonymous key
- `NEXT_PUBLIC_API_URL`: Your backend API URL

4. Run the development server:
```bash
npm run dev
```

## Authentication

The app uses Supabase for authentication with:
- Email/password registration and login
- JWT token validation
- Protected routes with AuthGuard
- Automatic token refresh

## API Integration

The frontend automatically includes Supabase JWT tokens in API requests to the Go backend.