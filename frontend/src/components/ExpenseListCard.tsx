import * as React from "react"
import { Card } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { ExpenseListItem } from "./ExpenseListItem"
import type { Expense, User } from "@/types"
import { useAuth } from "@/hooks/useAuth.ts";

interface ExpenseListCardProps {
  expenses: Expense[],
  members: User[],
}

export function ExpenseListCard({expenses, members }: ExpenseListCardProps) {
  const { userId } = useAuth();

  const membersMap = React.useMemo(() => {
    return (members || []).reduce((acc, user) => {
      acc[user.id] = `${user.name} ${user.surname}`;
      return acc;
    }, {} as Record<string, string>);
  }, [members]);

  return (
    <Card className="rounded-3xl border-zinc-100 shadow-sm overflow-hidden bg-white">
      {expenses && expenses.length > 0 ? (
        <div className="flex flex-col">
          {expenses.map((expense, index) => {
            const isMe = expense.payer_id === userId;
            const payerDisplayName = isMe ? "You" : (membersMap[expense.payer_id] || "Unknown");
            return (
              <React.Fragment key={expense.id}>
                <ExpenseListItem expense={expense} payerName={payerDisplayName} />
                {index < expenses.length - 1 && (
                  <Separator className="mx-4 bg-zinc-50 w-auto"/>
                )}
              </React.Fragment>
            );
          })}
        </div>
      ) : (
        <div className="py-12 text-center">
          <p className="text-sm text-zinc-400 italic">No expenses yet.</p>
        </div>
      )}
    </Card>
  )
}