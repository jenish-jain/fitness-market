import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import Link from 'next/link'

export default function Home() {
  return (
    <main className="container mx-auto p-6">
      <div className="flex flex-col items-center justify-center min-h-screen space-y-8">
        <div className="text-center space-y-4">
          <h1 className="text-4xl font-bold tracking-tight">Fitness Market</h1>
          <p className="text-xl text-muted-foreground max-w-md">
            Your one-stop marketplace for fitness equipment and services
          </p>
        </div>
        
        <Card className="w-full max-w-md">
          <CardHeader>
            <CardTitle>Welcome</CardTitle>
            <CardDescription>
              Get started by exploring our fitness marketplace
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <Button className="w-full">Browse Equipment</Button>
            <Button variant="outline" className="w-full">Find Services</Button>
            <div className="pt-4 border-t">
              <Link href="/login">
                <Button variant="secondary" className="w-full">Login</Button>
              </Link>
            </div>
          </CardContent>
        </Card>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 w-full max-w-4xl">
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Equipment</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                Find quality fitness equipment for your home gym
              </p>
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Services</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                Connect with certified trainers and fitness professionals
              </p>
            </CardContent>
          </Card>
          
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Community</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">
                Join a community of fitness enthusiasts
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    </main>
  )
}