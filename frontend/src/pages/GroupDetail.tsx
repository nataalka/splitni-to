import { useParams, Link } from "react-router-dom"
import { useQuery } from "@tanstack/react-query"
import api from "@/lib/api"
import type { Group } from "@/types"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ArrowUpDown, ChevronLeft, Plus, Users, Wallet } from "lucide-react"

export default function GroupDetailPage() {
  const {id} = useParams<{ id: string }>()

  const {data: group, isLoading, error} = useQuery<Group>({
    queryKey: ["group", id],
    queryFn: async () => {
      const response = await api.get(`/groups/${id}`)
      return response.data
    },
  })

  if (isLoading) return <div className="p-10 text-center text-muted-foreground">Loading group details...</div>
  if (error) return <div className="p-10 text-center text-destructive">Group not found.</div>

  return (
    <div className="max-w-5xl mx-auto p-4 md:p-6 space-y-6">
      {/* Back button & Title */}
      <div className="flex flex-col gap-4">
        <Link to="/groups">
          <Button variant="ghost" size="sm" className="-ml-2 text-muted-foreground">
            <ChevronLeft className="mr-1 h-4 w-4"/> Back to Groups
          </Button>
        </Link>
        <Card className="flex-row justify-between items-start p-4">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">{group?.name}</h1>
            <p className="text-muted-foreground">{group?.description}</p>
          </div>
          <Button className="bg-pink-600 hover:bg-pink-700">
            <Plus className="mr-2 h-4 w-4"/> Add Expense
          </Button>
        </Card>
      </div>

      {/* Stats / Overview */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle className="font-bold">Total Spent</CardTitle>
          <Wallet className="h-4 w-4"/>
        </CardHeader>
        <CardContent>
          <div className="font-bold">0.00 €</div>
        </CardContent>
      </Card>

      {/* Expenses */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>Recent Expenses</CardTitle>
          <ArrowUpDown className="h-4 w-4"/>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">No expenses yet. Start by adding one!</p>
        </CardContent>
      </Card>

      {/* Members */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>Members</CardTitle>
          <Users className="h-4 w-4"/>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">Members list coming soon...</p>
        </CardContent>
      </Card>
    </div>
  )
}