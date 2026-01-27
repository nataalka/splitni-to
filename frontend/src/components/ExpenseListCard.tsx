import * as React from "react"
import { Card } from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import { ExpenseListItem } from "./ExpenseListItem"
import type { Expense } from "@/types"

export function ExpenseListCard({ expenses }: { expenses: Expense[] | undefined }) {
  return (
    <Card className="rounded-3xl border-zinc-100 shadow-sm overflow-hidden bg-white">
      {expenses && expenses.length > 0 ? (
        <div className="flex flex-col">
          {expenses.map((expense, index) => (
            <React.Fragment key={expense.id}>
              <ExpenseListItem expense={expense} />
              {index < expenses.length - 1 && (
                <Separator className="mx-4 bg-zinc-50 w-auto" />
              )}
            </React.Fragment>
          ))}
        </div>
      ) : (
        <div className="py-12 text-center">
          <p className="text-sm text-zinc-400 italic">No expenses yet.</p>
        </div>
      )}
    </Card>
  )
}