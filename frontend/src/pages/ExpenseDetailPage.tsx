import { useParams } from "react-router-dom"
import { useQuery } from "@tanstack/react-query"
import api from "@/lib/api"
import { format } from "date-fns"
import { sk } from "date-fns/locale"
import { ArrowUpRight, Calendar, Receipt, Wallet } from "lucide-react"
import { Card } from "@/components/ui/card"
import type { ExpenseDetailed, User } from "@/types"
import { UserListCard } from "@/components/UserListCard.tsx";
import { StatsCard } from "@/components/StatsCard.tsx";
import { UserListItem } from "@/components/UserListItem.tsx";
import { ExpenseActions } from "@/actions/ExpenseActions.tsx";

export default function ExpenseDetailPage() {
  const {group_id, expense_id} = useParams<{ group_id: string; expense_id: string }>()

  const {data: expense, isLoading, error} = useQuery<ExpenseDetailed>({
    queryKey: ["expense", group_id, expense_id],
    queryFn: () => api.get(`/groups/${group_id}/expenses/${expense_id}`).then(res => res.data),
    enabled: !!expense_id,
  })

  const { data: groupMembers } = useQuery<User[]>({
    queryKey: ["group", group_id, "members"],
    queryFn: () => api.get(`/groups/${group_id}/members`).then(res => res.data),
    enabled: !!group_id,
  })

  if (isLoading) return <div className="p-20 text-center animate-pulse text-zinc-400 font-medium">Loading
    detail...</div>
  if (error || !expense) return <div className="p-20 text-center text-red-500 font-black">Expense not found.</div>

  const totalAmount = Number(expense.amount);
  const payerSplit = expense.splits.find(s => s.user.id === expense.payer.id);
  const payerShare = Number(payerSplit?.amount || 0);
  const getsBack = totalAmount - payerShare;

  return (
    <div className="max-w-2xl mx-auto space-y-4">
      {/* Header */}
      <header className="relative space-y-2">
        <div className="flex justify-between items-start gap-4">
          <div className="space-y-1">
            <p className="text-xs font-black text-zinc-600 uppercase tracking-widest">
              Expense Detail
            </p>
            <h1 className="text-3xl font-black text-zinc-900 tracking-tighter leading-[0.9]">
              {expense.description}
            </h1>
          </div>
          <ExpenseActions
            groupId={group_id!}
            expense={expense!}
            members={groupMembers || []}
          />
        </div>
        <div className="flex flex-wrap items-center gap-2 mb-2">
          {/* Group Badge */}
          <span
            className="bg-pink-600 text-white text-[9px] font-black uppercase tracking-[0.2em] px-2 py-1 rounded-full shadow-sm">
            {expense.group.name}
          </span>

          {/* Date Badge */}
          <div
            className="flex items-center gap-1.5 bg-zinc-100 text-zinc-500 text-[9px] font-black uppercase tracking-[0.1em] px-2.5 py-1 rounded-full border border-zinc-200">
            <Calendar className="h-3 w-3"/>
            {format(new Date(expense.created_at), "d. MMM yyyy", {locale: sk})}
          </div>
        </div>
      </header>

      {/* Payer Card */}
      <section className="space-y-3">
        <h2 className="text-xs font-black uppercase tracking-widest text-zinc-400 px-2">
          Paid by
        </h2>
        <Card className="p-2 shadow-sm">
          <UserListItem
            user={expense.payer}
            subtext={
              <span className="flex items-center gap-1.5 text-emerald-600 font-semibold">
                <span className="h-1 w-1 rounded-full bg-emerald-600"/>
                Covered the whole bill
              </span>
            }
          />
        </Card>
        <div className="grid grid-cols-2 gap-4">
          <StatsCard
            label="Gets Back"
            value={getsBack}
            icon={ArrowUpRight}
            variant="positive"
          />
          <StatsCard
            label="Total Paid"
            value={totalAmount}
            icon={Wallet}
            variant="default"
          />
        </div>
      </section>

      {/* Expense splits */}
      <section className="space-y-4">
        <div className="flex items-center justify-between px-2">
          <h2 className="text-xs font-black uppercase tracking-[0.2em] text-zinc-400">
            Expense Splits
          </h2>
          <Receipt className="h-4 w-4 text-zinc-300"/>
        </div>

        <UserListCard
          users={expense.splits.map(s => s.user)}
          renderActions={(user) => {
            const split = expense.splits.find(s => s.user.id === user.id);
            const amount = Number(split?.amount || 0);

            const isPayer = user.id === expense.payer.id;

            return (
              <div className="text-right min-w-[80px]">
                <p
                  className="text-[9px] uppercase tracking-widest font-bold text-pink-600/70">
                  {isPayer ? "Paid" : "Owes"}
                </p>
                <p className="text-lg font-bold leading-none text-pink-600">
                  {amount.toFixed(2)} €
                </p>
              </div>
            );
          }}
          renderSubtext={(user) => (
            <span className="text-[10px] text-zinc-400 font-medium">
              {user.id === expense.payer.id ? "Payer" : "Participant"}
            </span>
          )}
        />
      </section>
    </div>
  )
}