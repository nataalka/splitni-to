import { useParams } from "react-router-dom"
import { useQuery } from "@tanstack/react-query"
import api from "@/lib/api"
import type { Group, User } from "@/types"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { ArrowUpDown, Plus, Wallet, PieChart } from "lucide-react"
import { AddMemberDialog } from "@/dialogs/AddMemberDialog.tsx"
import { UserListCard } from "@/components/UserListCard"

export default function GroupDetailPage() {
  const { id } = useParams<{ id: string }>()

  const { data: group, isLoading: groupLoading, error } = useQuery<Group>({
    queryKey: ["group", id],
    queryFn: () => api.get(`/groups/${id}`).then(res => res.data),
  })

  const { data: members, isLoading: membersLoading } = useQuery<User[]>({
    queryKey: ["group", id, "members"],
    queryFn: () => api.get(`/groups/${id}/members`).then(res => res.data),
    enabled: !!id,
  })

  if (groupLoading || membersLoading) {
    return <div className="p-20 text-center text-zinc-400 animate-pulse font-medium">Loading group details...</div>
  }
  if (error) return <div className="p-20 text-center text-red-500 font-medium">Group not found.</div>

  return (
    <div className="max-w-2xl mx-auto space-y-8">
      {/* Header */}
      <section className="space-y-6">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-black text-zinc-900 tracking-tight">{group?.name}</h1>
            <p className="text-zinc-500 text-sm mt-1">{group?.description}</p>
          </div>
          <Button variant="pinkPrimary">
            <Plus className="h-4 w-4"/> Add Expense
          </Button>
        </div>

        {/* Stats Card */}
        <div className="grid grid-cols-2 gap-4">
          <Card className="rounded-3xl border-zinc-100 shadow-sm bg-white p-4 flex flex-col justify-between h-32">
            <div className="h-8 w-8 rounded-full bg-pink-50 flex items-center justify-center text-pink-600">
              <Wallet className="h-4 w-4"/>
            </div>
            <div>
              <p className="text-xs font-bold uppercase tracking-wider text-zinc-400">Total Spent</p>
              <p className="text-2xl font-black text-zinc-900">0.00 €</p>
            </div>
          </Card>

          <Card className="rounded-3xl border-zinc-100 shadow-sm bg-white p-4 flex flex-col justify-between h-32">
            <div className="h-8 w-8 rounded-full bg-blue-50 flex items-center justify-center text-blue-600">
              <PieChart className="h-4 w-4"/>
            </div>
            <div>
              <p className="text-xs font-bold uppercase tracking-wider text-zinc-400">Your Share</p>
              <p className="text-2xl font-black text-zinc-900 text-blue-600">0.00 €</p>
            </div>
          </Card>
        </div>
      </section>

      {/* Expenses */}
      <section className="space-y-3">
        <div className="flex items-center justify-between px-1">
          <h2 className="text-xs font-bold uppercase tracking-widest text-zinc-400">Recent Expenses</h2>
          <ArrowUpDown className="h-3 w-3 text-zinc-400"/>
        </div>
        <Card className="rounded-3xl border-zinc-100 shadow-sm bg-white overflow-hidden min-h-[100px] flex items-center justify-center">
          <CardContent className="p-0">
            <p className="text-sm text-zinc-400 italic py-8">No expenses yet. Start by adding one!</p>
          </CardContent>
        </Card>
      </section>

      {/* Members */}
      <section className="space-y-3">
        <div className="flex items-center justify-between px-1">
          <h2 className="text-xs font-bold uppercase tracking-widest text-zinc-400">Group Members</h2>
          <AddMemberDialog groupId={id || ""}/>
        </div>
        <UserListCard
          users={members}
          emptyMessage="This group has no members yet."
        />
      </section>
    </div>
  )
}