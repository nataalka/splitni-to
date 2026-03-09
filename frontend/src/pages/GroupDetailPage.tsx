import { useParams } from "react-router-dom"
import { useQuery } from "@tanstack/react-query"
import api from "@/lib/api"
import type { Group, User, Expense, MemberBalance } from "@/types"
import { Wallet, PieChart } from "lucide-react"
import { AddMemberDialog } from "@/dialogs/AddMemberDialog.tsx"
import { UserListCard } from "@/components/UserListCard"
import { ExpenseListCard } from "@/components/ExpenseListCard.tsx";
import { ExpenseFormDialog } from "@/dialogs/ExpenseFormDialog.tsx";
import { useAuth } from "@/hooks/useAuth.ts";
import { StatsCard } from "@/components/StatsCard.tsx";
import { GroupActions } from "@/actions/GroupActions.tsx";
import { EmptyBox } from "@/components/EmptyBox.tsx";

export default function GroupDetailPage() {
  const { id } = useParams<{ id: string }>()
  const { userId } = useAuth()

  {/* General Info */}
  const { data: group, isLoading: groupLoading, error } = useQuery<Group>({
    queryKey: ["group", id],
    queryFn: () => api.get(`/groups/${id}`).then(res => res.data),
  })

  {/* Members */}
  const { data: members, isLoading: membersLoading } = useQuery<User[]>({
    queryKey: ["group", id, "members"],
    queryFn: () => api.get(`/groups/${id}/members`).then(res => res.data),
    enabled: !!id,
  })

  {/* Expenses */}
  const { data: expenses } = useQuery<Expense[]>({
    queryKey: ["group", id, "expenses"],
    queryFn: () => api.get(`/groups/${id}/expenses`).then(res => res.data),
    enabled: !!id,
  })

  {/* Group Total */}
  const { data: totalData } = useQuery<{ total_spent: string }>({
    queryKey: ["group", id, "total-spent"],
    queryFn: () => api.get(`/groups/${id}/expenses/total`).then(res => res.data),
    enabled: !!id,
  })

  {/* User's balance */}
  const { data: usersBalanceData } = useQuery<MemberBalance>({
    queryKey: ["group", id, "balance", userId],
    queryFn: () => api.get(`/groups/${id}/balances/${userId}`).then(res => res.data),
    enabled: !!id && !!userId,
  })

  const myBalance = Number(usersBalanceData?.balance || 0);
  const isPositive = myBalance > 0;
  const isNegative = myBalance < 0;

  if (groupLoading || membersLoading) {
    return <div className="p-20 text-center text-zinc-400 animate-pulse font-medium">Loading group details...</div>
  }
  if (error) return <div className="p-20 text-center text-red-500 font-medium">Group not found.</div>

  return (
    <div className="max-w-2xl mx-auto space-y-4">
      {/* Header */}
      <section className="space-y-6">
        <div className="flex-col items-start justify-between gap-4">
          <p className="text-xs font-black text-zinc-600 uppercase tracking-widest">
            Group Detail
          </p>
          <div className="flex justify-between">
            <h1 className="text-3xl font-black text-zinc-900 tracking-tight">{group?.name}</h1>
            {group && <GroupActions group={group} />}
          </div>
          <p className="text-zinc-500 text-sm">{group?.description}</p>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-2 gap-4">
          <StatsCard
            label={isPositive ? "You are owed" : isNegative ? "You owe" : "Your balance"}
            value={Math.abs(myBalance)}
            icon={PieChart}
            variant={isPositive ? "positive" : isNegative ? "negative" : "default"}
          />
          <StatsCard
            label="Group Total"
            value={Number(totalData?.total_spent || 0)}
            icon={Wallet}
          />
        </div>

      </section>

      {/* Expenses */}
      <section className="space-y-3">
        <div className="flex items-center justify-between px-1">
          <h2 className="text-xs font-black uppercase tracking-widest text-zinc-400 px-2">
            Recent Expenses
          </h2>
          {members && <ExpenseFormDialog groupId={id!} members={members} />}
        </div>
        {expenses && members ? (
          <ExpenseListCard expenses={expenses} members={members}/>
        ) : (
          <EmptyBox description={"No expenses yet."}/>
        )}
      </section>

      {/* Members */}
      <section className="space-y-3">
        <div className="flex items-center justify-between px-1">
          <h2 className="text-xs font-black uppercase tracking-widest text-zinc-400 px-2">
            Group Members
          </h2>
          <AddMemberDialog groupId={id || ""}/>
        </div>
        {members ? (
          <UserListCard
            users={members}
            emptyMessage="This group has no members yet."
          />
        ) : (
          <EmptyBox description={"This group has no members yet."}/>
        )}
      </section>
    </div>
  )
}