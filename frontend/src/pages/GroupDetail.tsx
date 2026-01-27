import { useParams, Link } from "react-router-dom"
import { useQuery } from "@tanstack/react-query"
import api from "@/lib/api"
import type { Group } from "@/types"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ArrowUpDown, ChevronLeft, Plus, Wallet } from "lucide-react"
import { AddMemberDialog } from "@/dialogs/AddMemberDialog.tsx";

export default function GroupDetailPage() {
  const {id} = useParams<{ id: string }>()

  const {data: group, isLoading: groupLoading, error} = useQuery<Group>({
    queryKey: ["group", id],
    queryFn: async () => {
      const response = await api.get(`/groups/${id}`)
      return response.data
    },
  })

  const { data: members, isLoading: membersLoading } = useQuery<User[]>({
    queryKey: ["group", id, "members"],
    queryFn: () => api.get(`/groups/${id}/members`).then(res => res.data),
    enabled: !!id,
  })

  if (groupLoading || membersLoading) {
    return <div className="p-10 text-center text-muted-foreground">Loading group details...</div>
  }
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
          <AddMemberDialog groupId={id}/>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {members && members.length > 0 ? (
              members.map((member) => (
                <div key={member.id} className="flex items-center gap-3 p-2 rounded-lg hover:bg-zinc-50 transition-colors">
                  <div className="h-8 w-8 rounded-full bg-pink-100 text-pink-700 flex items-center justify-center font-bold text-xs">
                    {member.name[0]}{member.surname[0]}
                  </div>
                  <div>
                    <p className="text-sm font-medium leading-none">{member.name} {member.surname}</p>
                    <p className="text-xs text-muted-foreground">{member.email}</p>
                  </div>
                </div>
              ))
            ) : (
              <p className="text-sm text-muted-foreground text-center py-4">No members yet.</p>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}